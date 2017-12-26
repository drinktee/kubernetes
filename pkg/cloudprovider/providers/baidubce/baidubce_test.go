package baidubce

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"fmt"

	baidubce "github.com/drinktee/bce-sdk-go/bce"
	"github.com/drinktee/bce-sdk-go/clientset"
)

var cloudconfig = `{
    "AccessKeyID": "8e2fdc833cf44b4895afd0bce14f43cf",
    "SecretAccessKey": "7ae4ae1828694bbc814bb06fa87a43fa",
    "Region": "bj",
    "MasterID": "i-advasdv",
	"ClusterID":"k7s",
	"Debug":true,
	"Endpoint":"www.baidu.com",
	"FiberHost": "10.16.80.214",
    "FiberPort": 8563,
    "Pipe": "cce-test-pipe",
    "AclName": "cce-test",
    "AclPasswd": "cce-test",
	"AclToken": "$1$J02KsXCf$8S8EL4A.h1J/fmPvY.76e/",
    "PipeletId": 1
}`

func TestNewBCECloud(t *testing.T) {
	bc, err := NewBCECloud(bytes.NewBufferString(cloudconfig))
	if err != nil {
		t.Error(err)
	}
	if bc.ProviderName() != "baidubce" {
		t.Error("ProviderName error")
	}
	b := bc.(*BCECloud)
	fmt.Println(b.ClusterID)
	fmt.Println(b.clientSet.Cce().Endpoint)
	fmt.Println(b.clientSet.Blb().Endpoint)
}
func TestNewCloud(t *testing.T) {
	bc, err := newBceCloud()
	if err != nil {
		t.Error(err)
	}
	if bc.AccessKeyID != "8e2fdc833cf44b4895afd0bce14f43cf" {
		t.Error("accesskey error")
	}
	if bc.Endpoint != "" {
		fmt.Println(bc.clientSet.Cce().Endpoint)
		fmt.Println(bc.clientSet.Bcc().Endpoint)
		cceEnd := bc.clientSet.Cce().Endpoint
		fix := bc.Endpoint + "/internal-api"
		if fix != bc.clientSet.Blb().Endpoint {
			t.Errorf("fix endpoint error %s", fix)
		}
		if cceEnd != fix {
			t.Errorf("cceend error %s", cceEnd)
		}

	}

}
func newBceCloud() (*BCECloud, error) {
	var bc BCECloud
	var cfg CloudConfig
	err := json.Unmarshal(
		bytes.NewBufferString(cloudconfig).Bytes(),
		&cfg)
	if err != nil {
		return &bc, err
	}
	bc.CloudConfig = cfg
	cred := baidubce.NewCredentials(bc.AccessKeyID, bc.SecretAccessKey)
	bceConfig := baidubce.NewConfig(cred)
	bceConfig.Region = bc.Region
	// timeout need to set
	bceConfig.Timeout = 10 * time.Second
	// fix endpoint
	fixEndpoint := bc.Endpoint + "/internal-api"
	bceConfig.Endpoint = fixEndpoint
	bc.clientSet, err = clientset.NewFromConfig(bceConfig)
	if err != nil {
		return nil, err
	}
	bc.clientSet.Blb().SetDebug(true)
	bc.clientSet.Eip().SetDebug(true)
	bc.clientSet.Cce().SetDebug(true)
	return &bc, nil
}
