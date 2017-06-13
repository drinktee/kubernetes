package bcc

import (
	"encoding/json"

	"github.com/drinktee/bce-sdk-go/bce"
)

type Instance struct {
	InstanceId            string `json:"id"`
	InstanceName          string `json:"name"`
	Description           string `json:"desc"`
	Status                string `json:"status"`
	PaymentTiming         string `json:"paymentTiming"`
	CreationTime          string `json:"createTime"`
	ExpireTime            string `json:"expireTime"`
	PublicIP              string `json:"publicIp"`
	InternalIP            string `json:"internalIp"`
	CpuCount              int    `json:"cpuCount"`
	MemoryCapacityInGB    int    `json:"memoryCapacityInGB"`
	localDiskSizeInGB     int    `json:"localDiskSizeInGB"`
	ImageId               string `json:"imageId"`
	NetworkCapacityInMbps int    `json:"networkCapacityInMbps"`
	PlacementPolicy       string `json:"placementPolicy"`
	ZoneName              string `json:"zoneName"`
	SubnetId              string `json:"subnetId"`
	VpcId                 string `json:"vpcId"`
}

type ListInstancesResponse struct {
	Marker      string     `json:"marker"`
	IsTruncated bool       `json:"isTruncated"`
	NextMarker  string     `json:"nextMarker"`
	MaxKeys     int        `json:"maxKeys"`
	Instances   []Instance `json:"instances"`
}

type GetInstanceResponse struct {
	Ins Instance `json:"instance"`
}

// ListInstances gets all Instances.
func (c *Client) ListInstances(option *bce.SignOption) ([]Instance, error) {

	req, err := bce.NewRequest("GET", c.GetURL("v2/instance", nil), nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.SendRequest(req, option)

	if err != nil {
		return nil, err
	}

	bodyContent, err := resp.GetBodyContent()

	if err != nil {
		return nil, err
	}

	var insList *ListInstancesResponse
	err = json.Unmarshal(bodyContent, &insList)

	if err != nil {
		return nil, err
	}

	return insList.Instances, nil
}

func (c *Client) DescribeInstance(instanceId string, option *bce.SignOption) (*Instance, error) {

	req, err := bce.NewRequest("GET", c.GetURL("v2/instance"+"/"+instanceId, nil), nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.SendRequest(req, option)

	if err != nil {
		return nil, err
	}

	bodyContent, err := resp.GetBodyContent()

	if err != nil {
		return nil, err
	}

	var ins GetInstanceResponse
	err = json.Unmarshal(bodyContent, &ins)

	if err != nil {
		return nil, err
	}

	return &ins.Ins, nil
}
