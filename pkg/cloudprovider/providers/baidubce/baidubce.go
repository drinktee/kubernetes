/*
Copyright 2014 The Kubernetes Authors.

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
	"encoding/json"
	"io"
	"io/ioutil"
	"time"

	"cceevent/bigpiperpc"
	cceevent "cceevent/event"
	"fmt"

	baidubce "github.com/drinktee/bce-sdk-go/bce"
	"github.com/drinktee/bce-sdk-go/clientset"
	"k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/controller"
)

// ProviderName is the name of this cloud provider.
const ProviderName = "baidubce"

// CceUserAgent is prefix of http header UserAgent
const CceUserAgent = "cce-k8s:"

// BCECloud is an implementation of Interface, LoadBalancer and Instances for Baidu Compute Engine.
type BCECloud struct {
	CloudConfig
	clientSet clientset.Interface
	recorder  cceevent.Recorder
}

// CloudConfig wraps the settings for the BCE cloud provider.
type CloudConfig struct {
	ClusterID       string `json:"ClusterId"`
	ClusterName     string `json:"ClusterName"`
	AccessKeyID     string `json:"AccessKeyID"`
	SecretAccessKey string `json:"SecretAccessKey"`
	Region          string `json:"Region"`
	VpcID           string `json:"VpcId"`
	SubnetID        string `json:"SubnetId"`
	MasterID        string `json:"MasterId"`
	Endpoint        string `json:"Endpoint"`
	NodeIP          string `json:"NodeIP"`
	Debug           bool   `json:"Debug"`
	FiberHost       string `json:"FiberHost"`
	FiberPort       int    `json:"FiberPort"`
	bigpiperpc.BigpipeConfig
}

func init() {
	cloudprovider.RegisterCloudProvider(ProviderName, NewBCECloud)
}

// NewBCECloud returns a Cloud with initialized clients
func NewBCECloud(configReader io.Reader) (cloudprovider.Interface, error) {
	var bce BCECloud
	configContents, err := ioutil.ReadAll(configReader)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(configContents, &bce)
	if err != nil {
		return nil, err
	}
	err = validationBCECloud(&bce)
	if err != nil {
		return nil, err
	}
	cred := baidubce.NewCredentials(bce.AccessKeyID, bce.SecretAccessKey)
	bceConfig := baidubce.NewConfig(cred)
	bceConfig.Region = bce.Region
	// timeout need to set
	bceConfig.Timeout = 10 * time.Second
	// fix endpoint
	fixEndpoint := bce.Endpoint + "/internal-api"
	bceConfig.Endpoint = fixEndpoint
	// http request from cce's kubernetes has an useragent header
	// example: useragent: cce-k8s:c-adfdf
	bceConfig.UserAgent = CceUserAgent + bce.ClusterID
	bce.clientSet, err = clientset.NewFromConfig(bceConfig)
	if err != nil {
		return nil, err
	}
	// set debug for testing
	if bce.Debug {
		bce.clientSet.Blb().SetDebug(true)
		bce.clientSet.Eip().SetDebug(true)
		bce.clientSet.Bcc().SetDebug(true)
		bce.clientSet.Cce().SetDebug(true)
	}
	if hasBigpipe(&bce) {
		bigpipecfg := &bigpiperpc.BigpipeConfig{
			Pipe:      bce.Pipe,
			AclName:   bce.AclName,
			AclPasswd: bce.AclPasswd,
			AclToken:  bce.AclToken,
			PipeletId: bce.PipeletId,
		}
		bce.recorder = cceevent.NewQueueRecorder(bce.FiberHost, bce.FiberPort, 5, bigpipecfg, nil)
	}
	return &bce, nil
}

func validationBCECloud(bce *BCECloud) error {
	if bce.MasterID == "" {
		return fmt.Errorf("Cloud config mast have a Master ID")
	}
	if bce.ClusterID == "" {
		return fmt.Errorf("Cloud config mast have a ClusterID")
	}
	if bce.Endpoint == "" {
		return fmt.Errorf("Cloud config mast have a Endpoint")
	}
	return nil
}

func hasBigpipe(bce *BCECloud) bool {
	if bce.Pipe != "" && bce.AclName != "" && bce.AclPasswd != "" && bce.AclToken != "" && bce.PipeletId != 0 && bce.FiberHost != "" && bce.FiberPort != 0 {
		return true
	}
	return false
}

func (bc *BCECloud) CceEvent(resourceType string, resourceID string, resourceName string, source string, desc string) error {
	if bc.recorder != nil {
		return bc.recorder.Event(bc.ClusterID, resourceType, resourceID, resourceName, source, desc)
	}
	return nil
}

// LoadBalancer returns a balancer interface. Also returns true if the interface is supported, false otherwise.
func (bc *BCECloud) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return bc, true
}

// Instances returns an instances interface. Also returns true if the interface is supported, false otherwise.
func (bc *BCECloud) Instances() (cloudprovider.Instances, bool) {
	return bc, true
}

// Zones returns a zones interface. Also returns true if the interface is supported, false otherwise.
func (bc *BCECloud) Zones() (cloudprovider.Zones, bool) {
	return bc, true
}

// Clusters returns a clusters interface.  Also returns true if the interface is supported, false otherwise.
func (bc *BCECloud) Clusters() (cloudprovider.Clusters, bool) {
	return nil, false
}

// Routes returns a routes interface along with whether the interface is supported.
func (bc *BCECloud) Routes() (cloudprovider.Routes, bool) {
	return bc, true
}

// ScrubDNS provides an opportunity for cloud-provider-specific code to process DNS settings for pods.
func (bc *BCECloud) ScrubDNS(nameservers, searches []string) (nsOut, srchOut []string) {
	return nameservers, searches
}

// ProviderName returns the cloud provider ID.
func (bc *BCECloud) ProviderName() string {
	return ProviderName
}

// HasClusterID returns true if the cluster has a clusterID
func (bc *BCECloud) HasClusterID() bool {
	return true
}

// Initialize passes a Kubernetes clientBuilder interface to the cloud provider
func (bc *BCECloud) Initialize(clientBuilder controller.ControllerClientBuilder) {}
