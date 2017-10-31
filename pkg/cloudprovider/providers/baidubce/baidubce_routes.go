package baidubce

import (
	"fmt"

	"github.com/drinktee/bce-sdk-go/cce"

	"github.com/drinktee/bce-sdk-go/bcc"
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

// ListRoutes lists all managed routes that belong to the specified clusterName
func (bc *BCECloud) ListRoutes(clusterName string) (routes []*cloudprovider.Route, err error) {
	vpcid, err := bc.getVpcID()
	if err != nil {
		return nil, err
	}
	args := bcc.ListRouteArgs{
		VpcID: vpcid,
	}
	rs, err := bc.clientSet.Bcc().ListRouteTable(&args)
	if err != nil {
		return nil, err
	}
	inss, err := bc.clientSet.Cce().ListInstances(bc.ClusterID)
	if err != nil {
		return nil, err
	}
	var kubeRoutes []*cloudprovider.Route
	nodename := make(map[string]string)
	for _, ins := range inss {
		nodename[ins.InstanceId] = ins.InternalIP
	}
	for _, r := range rs {
		// filter instance route
		if r.NexthopType != "custom" {
			continue
		}
		var insName string
		insName, ok := nodename[r.NexthopID]
		if !ok {
			glog.V(4).Infof("Instance name not exist: %s", r.NexthopID)
			insName = "unknow"
		}
		route := &cloudprovider.Route{
			Name:            r.RouteRuleID,
			DestinationCIDR: r.DestinationAddress,
			TargetNode:      types.NodeName(insName),
		}
		kubeRoutes = append(kubeRoutes, route)
	}
	return kubeRoutes, nil
}

func (bc *BCECloud) getVpcRouteTable() ([]bcc.RouteRule, error) {
	vpcid, err := bc.getVpcID()
	if err != nil {
		return nil, err
	}
	args := bcc.ListRouteArgs{
		VpcID: vpcid,
	}
	rs, err := bc.clientSet.Bcc().ListRouteTable(&args)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

// CreateRoute creates the described managed route
// route.Name will be ignored, although the cloud-provider may use nameHint
// to create a more user-meaningful name.
func (bc *BCECloud) CreateRoute(clusterName string, nameHint string, kubeRoute *cloudprovider.Route) error {
	glog.V(4).Infof("create: creating route. clusterName=%q instance=%q cidr=%q", clusterName, kubeRoute.TargetNode, kubeRoute.DestinationCIDR)
	vpcRoutes, err := bc.getVpcRouteTable()
	if err != nil {
		return err
	}
	if len(vpcRoutes) < 1 {
		return fmt.Errorf("VPC route length error: length is : %d", len(vpcRoutes))
	}
	var insID string
	inss, err := bc.clientSet.Cce().ListInstances(bc.ClusterID)
	if err != nil {
		return err
	}
	for _, ins := range inss {
		if ins.InternalIP == string(kubeRoute.TargetNode) {
			insID = ins.InstanceId
		}
		if ins.Status == cce.InstanceStatusCreateFailed || ins.Status == cce.InstanceStatusDeleted || ins.Status == cce.InstanceStatusDeleting || ins.Status == cce.InstanceStatusError {
			glog.V(4).Infof("No need to create route, instance has a wrong status: %s", ins.Status)
			return nil
		}
	}
	var needDelete []string
	for _, vr := range vpcRoutes {
		if vr.DestinationAddress == kubeRoute.DestinationCIDR && vr.SourceAddress == "0.0.0.0/0" && vr.NexthopID == insID {
			glog.V(4).Infof("Route rule already exists.")
			return nil
		}
		if vr.DestinationAddress == kubeRoute.DestinationCIDR && vr.SourceAddress == "0.0.0.0/0" {
			needDelete = append(needDelete, vr.RouteRuleID)
		}
	}
	if len(needDelete) > 0 {
		for _, delRoute := range needDelete {
			err := bc.clientSet.Bcc().DeleteRoute(delRoute)
			if err != nil {
				glog.V(4).Infof("Delete VPC route error %s", err)
				return err
			}
		}
	}

	args := bcc.CreateRouteRuleArgs{
		RouteTableID:       vpcRoutes[0].RouteTableID,
		NexthopType:        "custom",
		Description:        "generated by bce-k8s",
		DestinationAddress: kubeRoute.DestinationCIDR,
		SourceAddress:      "0.0.0.0/0",
		NexthopID:          insID,
	}
	glog.V(4).Infof("CreateRoute: DestinationAddress is %s . NexthopID is : %s .", args.DestinationAddress, args.NexthopID)
	_, err = bc.clientSet.Bcc().CreateRouteRule(&args)
	return err
}

// DeleteRoute deletes the specified managed route
// Route should be as returned by ListRoutes
func (bc *BCECloud) DeleteRoute(clusterName string, kubeRoute *cloudprovider.Route) error {
	glog.V(4).Infof("delete: deleting route. clusterName=%q instance=%q cidr=%q", clusterName, kubeRoute.TargetNode, kubeRoute.DestinationCIDR)
	vpcTable, err := bc.getVpcRouteTable()
	if err != nil {
		glog.V(4).Infof("getVpcRouteTable error %s", err.Error())
		return err
	}
	for _, vr := range vpcTable {
		if vr.DestinationAddress == kubeRoute.DestinationCIDR && vr.SourceAddress == "0.0.0.0/0" {
			glog.V(4).Infof("DeleteRoute: DestinationAddress is %s .", vr.DestinationAddress)
			err := bc.clientSet.Bcc().DeleteRoute(vr.RouteRuleID)
			if err != nil {
				glog.V(4).Infof("Delete VPC route error %s", err.Error())
				return err
			}
		}
	}

	glog.V(4).Infof("delete: route deleted. clusterName=%q instance=%q cidr=%q", clusterName, kubeRoute.TargetNode, kubeRoute.DestinationCIDR)
	return nil
}
