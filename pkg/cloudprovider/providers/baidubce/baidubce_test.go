package baidubce

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	baidubce "github.com/drinktee/bce-sdk-go/bce"
	"github.com/drinktee/bce-sdk-go/clientset"
)

var cloudconfig = `{
    "AccessKeyID": "8e2fdc833cf44b4895afd0bce14f43cf",
    "SecretAccessKey": "7ae4ae1828694bbc814bb06fa87a43fa",
    "region": "bj",
    "masterId": "i-advasdv"
}`

func TestNewBCECloud(t *testing.T) {
	bc, err := NewBCECloud(bytes.NewBufferString(cloudconfig))
	if err != nil {
		t.Error(err)
	}
	if bc.ProviderName() != "baidubce" {
		t.Error("ProviderName error")
	}
}
func TestNewCloud(t *testing.T) {
	bc, err := newBceCloud()
	if err != nil {
		t.Error(err)
	}
	if bc.AccessKeyID != "8e2fdc833cf44b4895afd0bce14f43cf" {
		t.Error("accesskey error")
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
	bc.clientSet, err = clientset.NewFromConfig(bceConfig)
	if err != nil {
		return nil, err
	}
	bc.clientSet.Blb().SetDebug(true)
	bc.clientSet.Eip().SetDebug(true)
	return &bc, nil
}
