/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package baidubce

import (
	"fmt"

	"time"

	"github.com/drinktee/bce-sdk-go/blb"
	"github.com/drinktee/bce-sdk-go/eip"
	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

// GetLoadBalancer returns whether the specified load balancer exists, and
// if so, what its status is.
// Implementations must treat the *v1.Service parameter as read-only and not modify it.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
// TODO: Break this up into different interfaces (LB, etc) when we have more than one type of service
func (bc *BCECloud) GetLoadBalancer(clusterName string, service *v1.Service) (status *v1.LoadBalancerStatus, exists bool, err error) {
	lb, exists, err := bc.getBCELoadBalancer(cloudprovider.GetLoadBalancerName(service))
	if err != nil || !exists {
		return nil, exists, err
	}
	ip := lb.Address
	if lb.PublicIp != "" {
		ip = lb.PublicIp
	}
	glog.V(4).Infof("GetLoadBalancer ip: %s", ip)
	return &v1.LoadBalancerStatus{
		Ingress: []v1.LoadBalancerIngress{{IP: ip}},
	}, true, nil
}

func (bc *BCECloud) getBCELoadBalancer(name string) (lb blb.LoadBalancer, exists bool, err error) {
	args := blb.DescribeLoadBalancersArgs{
		LoadBalancerName: name,
	}
	lbs, err := bc.clientSet.Blb().DescribeLoadBalancers(&args)
	if err != nil {
		glog.V(2).Infof("getBCELoadBalancer  not exists blb! %v", args)
		return blb.LoadBalancer{}, false, err
	}

	if len(lbs) < 1 {
		glog.V(2).Infof("getBCELoadBalancer  not exists blb! len(lbs) < 1  %v", args)
		return blb.LoadBalancer{}, false, nil
	}

	return lbs[0], true, nil
}

// EnsureLoadBalancer creates a new load balancer 'name', or updates the existing one. Returns the status of the balancer
// Implementations must treat the *v1.Service and *v1.Node
// parameters as read-only and not modify them.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (bc *BCECloud) EnsureLoadBalancer(clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	glog.V(2).Infof("baidubce.EnsureLoadBalancer(%v, %v, %v, %v, %v, %v, %v, %v,%v)",
		clusterName, service.Namespace, service.Name, bc.Region, service.Spec.LoadBalancerIP, service.Spec.Ports, nodes)
	// if service.Spec.SessionAffinity != v1.ServiceAffinityNone {
	// 	// Does not support SessionAffinity
	// 	return nil, fmt.Errorf("unsupported load balancer affinity: %v", service.Spec.SessionAffinity)
	// }
	if len(service.Spec.Ports) == 0 {
		return nil, fmt.Errorf("requested load balancer with no ports")
	}
	if service.Spec.LoadBalancerIP != "" {
		return nil, fmt.Errorf("LoadBalancerIP cannot be specified for BLB")
	}
	for _, port := range service.Spec.Ports {
		if port.Protocol != v1.ProtocolTCP {
			return nil, fmt.Errorf("Only TCP LoadBalancer is supported for Baidu K8S")
		}
	}
	lb, exists, err := bc.getBCELoadBalancer(cloudprovider.GetLoadBalancerName(service))
	if err != nil {
		return nil, err
	}
	lbName := cloudprovider.GetLoadBalancerName(service)
	if !exists {
		glog.V(4).Infoln("EnsureLoadBalancer create not exists blb!")
		args := blb.CreateLoadBalancerArgs{
			Name: lbName,
		}
		_, err := bc.clientSet.Blb().CreateLoadBalancer(&args)
		if err != nil {
			return nil, err
		}
		argsDesc := blb.DescribeLoadBalancersArgs{
			LoadBalancerName: lbName,
		}
		lbs, err := bc.clientSet.Blb().DescribeLoadBalancers(&argsDesc)
		if err != nil {
			return nil, err
		}
		lb = lbs[0]
	} else {
		glog.V(4).Infoln("EnsureLoadBalancer: blb already exists!")
	}
	lb, err = bc.waitForLoadBalancer(&lb)
	if err != nil {
		return nil, err
	}
	// time.Sleep(60 * time.Second)
	glog.V(2).Infoln("EnsureLoadBalancer: reconcileListeners!")
	err = bc.reconcileListeners(service, &lb)
	if err != nil {
		return nil, err
	}
	glog.V(2).Infoln("EnsureLoadBalancer: reconcileBackendServers!")
	// time.Sleep(60 * time.Second)
	lb, err = bc.waitForLoadBalancer(&lb)
	if err != nil {
		return nil, err
	}
	err = bc.reconcileBackendServers(nodes, &lb)
	if err != nil {
		return nil, err
	}

	glog.V(2).Infoln("EnsureLoadBalancer: createEIP!")
	// TODO
	pubIP, err := bc.createEIP(&lb)
	if err != nil {
		return nil, err
	}
	glog.V(4).Infof("EnsureLoadBalancer: LoadBalancerIngress= %v  pubIP is %s", lb.PublicIp, pubIP)
	return &v1.LoadBalancerStatus{
		Ingress: []v1.LoadBalancerIngress{{IP: pubIP}},
	}, nil
}

// UpdateLoadBalancer updates hosts under the specified load balancer.
func (bc *BCECloud) UpdateLoadBalancer(clusterName string, service *v1.Service, nodes []*v1.Node) error {
	_, err := bc.EnsureLoadBalancer(clusterName, service, nodes)
	return err
}

// EnsureLoadBalancerDeleted deletes the specified load balancer if it
// exists, returning nil if the load balancer specified either didn't exist or
// was successfully deleted.
// This construction is useful because many cloud providers' load balancers
// have multiple underlying components, meaning a Get could say that the LB
// doesn't exist even if some part of it is still laying around.
func (bc *BCECloud) EnsureLoadBalancerDeleted(clusterName string, service *v1.Service) error {
	lbName := cloudprovider.GetLoadBalancerName(service)
	serviceName := getServiceName(service)
	glog.V(2).Infof("delete(%s): START clusterName=%q lbName=%q", serviceName, clusterName, lbName)

	// reconcile logic is capable of fully reconcile, so we can use this to delete
	service.Spec.Ports = []v1.ServicePort{}

	lb, existsLb, err := bc.getBCELoadBalancer(lbName)
	glog.V(4).Infof("EnsureLoadBalancerDeleted getBCELoadBalancer : %s", lb.Name)
	if err != nil {
		glog.V(4).Infof("EnsureLoadBalancerDeleted get error: %s", err.Error())
		return err
	}
	if !existsLb {
		glog.V(4).Infof("BCELoadBalancer not exists: %s", lbName)
		return nil
	}
	// start delete blb and eip, delete blb first
	glog.V(4).Infof("Start delete BLB: %s", lb.Name)
	args := blb.DeleteLoadBalancerArgs{
		LoadBalancerId: lb.BlbId,
	}
	err = bc.clientSet.Blb().DeleteLoadBalancer(&args)
	if err != nil {
		return err
	}
	// delete EIP
	if lb.PublicIp != "" {
		glog.V(4).Infof("Start delete EIP: %s", lb.PublicIp)
		err = bc.deleteEIP(&lb)
		if err != nil {
			return err
		}
	}
	glog.V(2).Infof("delete(%s): FINISH", serviceName)
	return nil
}

// This returns a human-readable version of the Service used to tag some resources.
// This is only used for human-readable convenience, and not to filter.
func getServiceName(service *v1.Service) string {
	return fmt.Sprintf("%s/%s", service.Namespace, service.Name)
}

// PortListener describe listener port
type PortListener struct {
	Port     int
	Protocol string
	NodePort int32
}

func (bc *BCECloud) reconcileListeners(service *v1.Service, lb *blb.LoadBalancer) error {
	expected := make(map[int]PortListener)
	// add expected ports
	for _, v1 := range service.Spec.Ports {
		expected[int(v1.Port)] = PortListener{
			Port:     int(v1.Port),
			Protocol: string(v1.Protocol),
			NodePort: (v1.NodePort),
		}
	}
	// delete or update unexpected ports
	all, err := bc.getAllListeners(lb)
	if err != nil {
		return err
	}
	deleteList := []PortListener{}
	for _, l := range all {
		port, ok := expected[l.Port]
		if !ok {
			// delete listener port
			// add to deleteList
			deleteList = append(deleteList, l)
		} else {
			if l != port {
				// update listener port
				err := bc.updateListener(lb, port)
				if err != nil {
					return err
				}
				delete(expected, l.Port)
			}
		}
	}
	// delete listener
	if len(deleteList) > 0 {
		err = bc.deleteListener(lb, deleteList)
		if err != nil {
			return err
		}
	}

	// create expected listener
	for _, pl := range expected {
		err := bc.createListener(lb, pl)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bc *BCECloud) findPortListener(lb *blb.LoadBalancer, port int, proto string) (PortListener, error) {
	switch proto {
	case "http":
	case "tcp":
		args := blb.DescribeTCPListenerArgs{
			LoadBalancerId: lb.BlbId,
			ListenerPort:   port,
		}
		ls, err := bc.clientSet.Blb().DescribeTCPListener(&args)
		if err != nil {
			return PortListener{}, err
		}
		if len(ls) < 1 {
			return PortListener{}, fmt.Errorf("there is no tcp listener blb:%s  port:%d", lb.Name, port)
		}
		return PortListener{
			Port:     ls[0].ListenerPort,
			NodePort: int32(ls[0].BackendPort),
			Protocol: proto,
		}, nil
	case "https":
	case "udp":
	}
	return PortListener{}, fmt.Errorf("protocol not match: %s", proto)
}

func (bc *BCECloud) getAllListeners(lb *blb.LoadBalancer) ([]PortListener, error) {
	allListeners := []PortListener{}
	// add TCPlisteners
	args := blb.DescribeTCPListenerArgs{
		LoadBalancerId: lb.BlbId,
	}
	ls, err := bc.clientSet.Blb().DescribeTCPListener(&args)
	if err != nil {
		return nil, err
	}
	for _, listener := range ls {
		allListeners = append(allListeners, PortListener{
			Port:     listener.ListenerPort,
			Protocol: "TCP",
			NodePort: int32(listener.BackendPort),
		})
	}

	// add HTTPlisteners HTTPS UDP
	// TODO
	return allListeners, nil
}

func (bc *BCECloud) createListener(lb *blb.LoadBalancer, pl PortListener) error {
	switch pl.Protocol {
	case "HTTP":
	case "TCP":
		args := blb.CreateTCPListenerArgs{
			LoadBalancerId: lb.BlbId,
			ListenerPort:   pl.Port,
			BackendPort:    int(pl.NodePort),
			Scheduler:      "RoundRobin",
		}
		err := bc.clientSet.Blb().CreateTCPListener(&args)
		if err != nil {
			return err
		}
		return nil
	case "HTTPS":
	case "UDP":
	}
	return fmt.Errorf("CreateListener protocol not match: %s", pl.Protocol)
}

func (bc *BCECloud) updateListener(lb *blb.LoadBalancer, pl PortListener) error {
	switch pl.Protocol {
	case "http":
	case "tcp":
		args := blb.UpdateTCPListenerArgs{
			LoadBalancerId: lb.BlbId,
			ListenerPort:   pl.Port,
			BackendPort:    int(pl.NodePort),
			Scheduler:      "RoundRobin",
		}
		err := bc.clientSet.Blb().UpdateTCPListener(&args)
		if err != nil {
			return err
		}
		return nil
	case "https":
	case "udp":
	}
	return fmt.Errorf("updateListener protocol not match: %s", pl.Protocol)
}

func (bc *BCECloud) deleteListener(lb *blb.LoadBalancer, pl []PortListener) error {
	portList := []int{}
	for _, l := range pl {
		portList = append(portList, l.Port)
	}
	args := blb.DeleteListenersArgs{
		LoadBalancerId: lb.BlbId,
		PortList:       portList,
	}
	err := bc.clientSet.Blb().DeleteListeners(&args)
	if err != nil {
		return err
	}
	return nil
}

const DEFAULT_SERVER_WEIGHT = 100

func (bc *BCECloud) getAllBackendServer(lb *blb.LoadBalancer) ([]blb.BackendServer, error) {
	args := blb.DescribeBackendServersArgs{
		LoadBalancerId: lb.BlbId,
	}
	bs, err := bc.clientSet.Blb().DescribeBackendServers(&args)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func (bc *BCECloud) reconcileBackendServers(nodes []*v1.Node, lb *blb.LoadBalancer) error {
	expectedServer := make(map[string]string)
	for _, node := range nodes {
		expectedServer[node.Spec.ExternalID] = node.ObjectMeta.Name
	}
	allBS, err := bc.getAllBackendServer(lb)
	if err != nil {
		return err
	}
	removeList := []string{}
	// remove unexpected servers
	for _, bs := range allBS {
		_, exists := expectedServer[bs.InstanceId]
		if !exists {
			removeList = append(removeList, bs.InstanceId)
			delete(expectedServer, bs.InstanceId)
		}
	}
	if len(removeList) > 0 {
		args := blb.RemoveBackendServersArgs{
			LoadBalancerId:    lb.BlbId,
			BackendServerList: removeList,
		}
		err = bc.clientSet.Blb().RemoveBackendServers(&args)
		if err != nil {
			return err
		}

	}
	addList := []blb.BackendServer{}
	// add expected servers
	for insID, nodeName := range expectedServer {
		addList = append(addList, blb.BackendServer{
			InstanceId: insID,
			Weight:     DEFAULT_SERVER_WEIGHT,
		})
		glog.V(4).Infof("add node %s", nodeName)
	}
	if len(addList) > 0 {
		args := blb.AddBackendServersArgs{
			LoadBalancerId:    lb.BlbId,
			BackendServerList: addList,
		}
		err = bc.clientSet.Blb().AddBackendServers(&args)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bc *BCECloud) createEIP(lb *blb.LoadBalancer) (string, error) {
	bill := &eip.Billing{
		PaymentTiming: "Postpaid",
		BillingMethod: "ByTraffic",
	}
	args := &eip.CreateEipArgs{
		BandwidthInMbps: 1000,
		Billing:         bill,
		Name:            lb.Name,
	}
	glog.V(4).Infof("CreateEip:  %v", args)
	ip, err := bc.clientSet.Eip().CreateEip(args)
	if err != nil {
		return "", err
	}
	argsGet := eip.GetEipsArgs{
		Ip: ip,
	}
	eips, err := bc.clientSet.Eip().GetEips(&argsGet)
	if err != nil {
		return "", err
	}
	if len(eips) > 0 {
		eipStatus := eips[0].Status
		for index := 0; (index < 10) && (eipStatus != "available"); index++ {
			glog.V(4).Infof("Eip: %s is not available, retry:  %d", ip, index)
			time.Sleep(10 * time.Second)
			eips, err := bc.clientSet.Eip().GetEips(&argsGet)
			if err != nil {
				return "", err
			}
			eipStatus = eips[0].Status
		}
		glog.V(4).Infof("Eip status is: %s", eipStatus)
	}

	for index := 0; (index < 10) && (lb.Status != "available"); index++ {
		glog.V(4).Infof("BLB: %s is not available, retry:  %d", lb.Name, index)
		time.Sleep(10 * time.Second)
		newlb, exist, err := bc.getBCELoadBalancer(lb.Name)
		if err != nil {
			glog.V(4).Infof("getBCELoadBalancer error: %s", lb.Name)
			return "", err
		}
		if !exist {
			glog.V(4).Infof("getBCELoadBalancer not exist: %s", lb.Name)
			return "", fmt.Errorf("BLB not exists:%s", lb.Name)
		}
		lb = &newlb
		glog.V(4).Infof("BLB status is : %s", lb.Status)
	}
	argsBind := &eip.BindEipArgs{
		Ip:           ip,
		InstanceId:   lb.BlbId,
		InstanceType: eip.BLB,
	}
	glog.V(4).Infof("BindEip:  %v", argsBind)
	glog.V(4).Infof("Bind BLB: %v", lb)
	err = bc.clientSet.Eip().BindEip(argsBind)
	if err != nil {
		glog.V(4).Infof("BindEip error: %v", err)
		return "", err
	}
	lb.PublicIp = ip
	glog.V(4).Infof("createEIP: lb.PublicIp is %s", lb.PublicIp)
	return ip, nil
}

func (bc *BCECloud) deleteEIP(lb *blb.LoadBalancer) error {

	// err := bc.clientSet.Eip().UnbindEip(&args)
	// if err != nil {
	// 	return err
	// }
	argsGet := eip.GetEipsArgs{
		Ip: lb.PublicIp,
	}
	eips, err := bc.clientSet.Eip().GetEips(&argsGet)
	if err != nil {
		return err
	}
	if len(eips) > 0 {
		eipStatus := eips[0].Status
		for index := 0; (index < 10) && (eipStatus != "available"); index++ {
			glog.V(4).Infof("Eip: %s is not available, retry:  %d", lb.PublicIp, index)
			time.Sleep(10 * time.Second)
			eips, err := bc.clientSet.Eip().GetEips(&argsGet)
			if err != nil {
				return err
			}
			eipStatus = eips[0].Status
		}
	}
	args := eip.EipArgs{
		Ip: lb.PublicIp,
	}
	err = bc.clientSet.Eip().DeleteEip(&args)
	if err != nil {
		return err
	}
	return nil
}

func (bc *BCECloud) waitForLoadBalancer(lb *blb.LoadBalancer) (blb.LoadBalancer, error) {
	// var newlb blb.LoadBalancer
	for index := 0; (index < 10) && (lb.Status != "available"); index++ {
		glog.V(4).Infof("BLB: %s is not available, retry:  %d", lb.Name, index)
		time.Sleep(10 * time.Second)
		newlb, exist, err := bc.getBCELoadBalancer(lb.Name)
		if err != nil {
			glog.V(4).Infof("getBCELoadBalancer error: %s", lb.Name)
			return newlb, err
		}
		if !exist {
			glog.V(4).Infof("getBCELoadBalancer not exist: %s", lb.Name)
			return newlb, fmt.Errorf("BLB not exists:%s", lb.Name)
		}
		lb = &newlb
		glog.V(4).Infof("BLB status is : %s", lb.Status)
	}

	return *lb, nil
}
