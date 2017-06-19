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

	"github.com/drinktee/bce-sdk-go/bcc"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

// NodeAddresses returns the addresses of the specified instance.
func (bc *BCECloud) NodeAddresses(name types.NodeName) ([]v1.NodeAddress, error) {
	ip, err := bc.getIPForMachine(name)
	if err != nil {
		return nil, err
	}
	return []v1.NodeAddress{
		{Type: v1.NodeInternalIP, Address: ip},
		{Type: v1.NodeHostName, Address: string(name)},
	}, nil
}

func (bc *BCECloud) getIPForMachine(name types.NodeName) (string, error) {
	ins, err := bc.clientSet.Bcc().ListInstances(nil)
	if err != nil {
		return "", err
	}
	insName := string(name)
	for _, i := range ins {
		if i.InstanceName == insName {
			return i.InternalIP, nil
		}
	}
	return "", cloudprovider.InstanceNotFound
}

func (bc *BCECloud) getVirtualMachine(name types.NodeName) (vm bcc.Instance, err error) {
	ins, err := bc.clientSet.Bcc().ListInstances(nil)
	if err != nil {
		return vm, err
	}
	insName := string(name)
	for _, i := range ins {
		if i.InstanceName == insName {
			return i, nil
		}
	}
	return vm, cloudprovider.InstanceNotFound
}

// ExternalID returns the cloud provider ID of the specified instance (deprecated).
func (bc *BCECloud) ExternalID(name types.NodeName) (string, error) {
	return bc.InstanceID(name)
}

// InstanceID returns the cloud provider ID of the specified instance.
// Note that if the instance does not exist or is no longer running, we must return ("", cloudprovider.InstanceNotFound)
func (bc *BCECloud) InstanceID(name types.NodeName) (string, error) {
	machine, err := bc.getVirtualMachine(name)
	if err != nil {
		return "", err
	}
	return machine.InstanceId, nil
}

// InstanceType returns the type of the specified instance.
// Note that if the instance does not exist or is no longer running, we must return ("", cloudprovider.InstanceNotFound)
// (Implementer Note): This is used by kubelet. Kubelet will label the node. Real log from kubelet:
//       Adding node label from cloud provider: beta.kubernetes.io/instance-type=[value]
func (bc *BCECloud) InstanceType(name types.NodeName) (string, error) {
	_, err := bc.getVirtualMachine(name)
	if err != nil {
		return "", err
	}
	return string("NoType"), nil
}

// AddSSHKeyToAllInstances adds an SSH public key as a legal identity for all instances
// expected format for the key is standard ssh-keygen format: <protocol> <blob>
func (bc *BCECloud) AddSSHKeyToAllInstances(user string, keyData []byte) error {
	return fmt.Errorf("not supported")
}

// CurrentNodeName returns the name of the node we are currently running on
// On most clouds (e.g. GCE) this is the hostname, so we provide the hostname
func (bc *BCECloud) CurrentNodeName(hostname string) (types.NodeName, error) {
	return types.NodeName(hostname), nil
}
