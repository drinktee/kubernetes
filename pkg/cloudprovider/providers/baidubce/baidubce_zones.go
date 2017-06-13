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
	"os"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

// GetZone returns the Zone containing the current failure zone and locality region that the program is running in
func (bc *BCECloud) GetZone() (cloudprovider.Zone, error) {
	host, err := os.Hostname()
	zone := cloudprovider.Zone{}
	if err != nil {
		return zone, err
	}
	ins, err := bc.getVirtualMachine(types.NodeName(host))
	if err != nil {
		return zone, err
	}
	zone.FailureDomain = ins.ZoneName
	zone.Region = bc.Region
	return zone, nil
}
