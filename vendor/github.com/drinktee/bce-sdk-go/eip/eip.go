package eip

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/drinktee/bce-sdk-go/bce"
)

type Eip struct {
	Name            string          `json:"name"`
	Eip             string          `json:"eip"`
	Status          string          `json:"status"`
	EipInstanceType EipInstanceType `json:"eipInstanceType"`
	InstanceType    InstanceType    `json:"instanceType"`
	InstanceId      string          `json:"instanceId"`
	ShareGroupId    string          `json:"shareGroupId"`
	BandwidthInMbps int             `json:"bandwidthInMbps"`
	PaymentTiming   string          `json:"paymentTiming"`
	BillingMethod   string          `json:"billingMethod"`
	CreateTime      string          `json:"createTime"`
	ExpireTime      string          `json:"expireTime"`
}
type Billing struct {
	PaymentTiming string `json:"paymentTiming"`
	BillingMethod string `json:"billingMethod"`
}
type Reservation struct {
	ReservationLength   int    `json:"reservationLength"`
	ReservationTimeUnit string `json:"reservationTimeUnit"`
}
type CreateEipArgs struct {
	//  公网带宽，单位为Mbps。
	// 对于prepay以及bandwidth类型的EIP，限制为为1~200之间的整数，
	// 对于traffic类型的EIP，限制为1~1000之前的整数。
	BandwidthInMbps int      `json:"bandwidthInMbps"`
	Billing         *Billing `json:"billing"`
	Name            string   `json:"name,omitempty"`
}

type CreateEipResponse struct {
	Ip string `json:"eip"`
}

type InstanceType string

const (
	BCC InstanceType = "BCC"
	BLB InstanceType = "BLB"
)

const (
	PAYMENTTIMING_PREPAID     string = "Prepaid"
	PAYMENTTIMING_POSTPAID    string = "Postpaid"
	BILLINGMETHOD_BYTRAFFIC   string = "ByTraffic"
	BILLINGMETHOD_BYBANDWIDTH string = "ByBandwidth"
)

type EipInstanceType string

const (
	NORMAL EipInstanceType = "normal"
	SHARED EipInstanceType = "shared"
)

func (args *CreateEipArgs) validate() error {
	if args == nil {
		return fmt.Errorf("CreateEipArgs need args")
	}
	if args.BandwidthInMbps == 0 {
		return fmt.Errorf("CreateEipArgs need BandwidthInMbps")
	}
	if args.Billing == nil {
		return fmt.Errorf("CreateEipArgs need Billing")
	}
	return nil
}
func (c *Client) CreateEip(args *CreateEipArgs) (string, error) {
	err := args.validate()
	if err != nil {
		return "", err
	}
	params := map[string]string{
		"clientToken": c.GenerateClientToken(),
	}
	postContent, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	req, err := bce.NewRequest("POST", c.GetURL("v1/eip", params), bytes.NewBuffer(postContent))
	if err != nil {
		return "", err
	}
	resp, err := c.SendRequest(req, nil)
	if err != nil {
		return "", err
	}
	bodyContent, err := resp.GetBodyContent()

	if err != nil {
		return "", err
	}
	var blbsResp *CreateEipResponse
	err = json.Unmarshal(bodyContent, &blbsResp)

	if err != nil {
		return "", err
	}
	return blbsResp.Ip, nil

}

type ResizeEipArgs struct {
	BandwidthInMbps int    `json:"newBandwidthInMbps"`
	Ip              string `json:"-"`
}

func (args *ResizeEipArgs) validate() error {
	if args == nil {
		return fmt.Errorf("ResizeEipArgs need args")
	}
	if args.Ip == "" {
		return fmt.Errorf("ResizeEipArgs need ip")
	}
	if args.BandwidthInMbps == 0 {
		return fmt.Errorf("ResizeEipArgs need BandwidthInMbps")
	}
	return nil
}

func (c *Client) ResizeEip(args *ResizeEipArgs) error {
	err := args.validate()
	if err != nil {
		return err
	}
	params := map[string]string{
		"resize":      "",
		"clientToken": c.GenerateClientToken(),
	}
	postContent, err := json.Marshal(args)
	if err != nil {
		return err
	}
	// url := "http://" + Endpoint[c.GetRegion()] + "/v1/eip" + "/" + args.Ip + "?" + "resize&" + "clientToken=" + c.GenerateClientToken()
	// req, err := bce.NewRequest("PUT", url, bytes.NewBuffer(postContent))
	req, err := bce.NewRequest("PUT", c.GetURL("v1/eip"+"/"+args.Ip, params), bytes.NewBuffer(postContent))
	if err != nil {
		return err
	}
	_, err = c.SendRequest(req, nil)
	if err != nil {
		return err
	}
	return nil

}

type BindEipArgs struct {
	Ip           string       `json:"-"`
	InstanceType InstanceType `json:"instanceType"`
	InstanceId   string       `json:"instanceId"`
}

func (args *BindEipArgs) validate() error {
	if args == nil {
		return fmt.Errorf("BindEip need args")
	}
	if args.Ip == "" {
		return fmt.Errorf("BindEip need ip")
	}
	if args.InstanceType == "" {
		return fmt.Errorf("BindEip need InstanceType")
	}
	if args.InstanceId == "" {
		return fmt.Errorf("BindEip need InstanceId")
	}
	return nil
}

func (c *Client) BindEip(args *BindEipArgs) error {
	err := args.validate()
	if err != nil {
		return err
	}
	params := map[string]string{
		"bind":        "",
		"clientToken": c.GenerateClientToken(),
	}
	postContent, err := json.Marshal(args)
	if err != nil {
		return err
	}
	req, err := bce.NewRequest("PUT", c.GetURL("v1/eip"+"/"+args.Ip, params), bytes.NewBuffer(postContent))
	if err != nil {
		return err
	}
	_, err = c.SendRequest(req, nil)
	if err != nil {
		return err
	}
	return nil

}

type EipArgs struct {
	Ip string `json:"-"`
}

func (args *EipArgs) validate() error {
	if args == nil {
		return fmt.Errorf("EipArgs need args")
	}
	if args.Ip == "" {
		return fmt.Errorf("EipArgs need ip")
	}
	return nil
}

func (c *Client) UnbindEip(args *EipArgs) error {
	err := args.validate()
	if err != nil {
		return err
	}
	params := map[string]string{
		"unbind":      "",
		"clientToken": c.GenerateClientToken(),
	}

	req, err := bce.NewRequest("PUT", c.GetURL("v1/eip"+"/"+args.Ip, params), nil)
	if err != nil {
		return err
	}
	_, err = c.SendRequest(req, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteEip(args *EipArgs) error {
	err := args.validate()
	if err != nil {
		return err
	}
	params := map[string]string{
		"clientToken": c.GenerateClientToken(),
	}
	req, err := bce.NewRequest("DELETE", c.GetURL("v1/eip"+"/"+args.Ip, params), nil)
	if err != nil {
		return err
	}
	_, err = c.SendRequest(req, nil)
	if err != nil {
		return err
	}
	return nil
}

type GetEipsArgs struct {
	Ip           string       `json:"-"`
	InstanceType InstanceType `json:"instanceType"`
	InstanceId   string       `json:"instanceId"`
}

type GetEipsResponse struct {
	EipList     []Eip  `json:"eipList"`
	Marker      string `json:"marker"`
	IsTruncated bool   `json:"isTruncated"`
	NextMarker  string `json:"nextMarker"`
	MaxKeys     int    `json:"maxKeys"`
}

func (c *Client) GetEips(args *GetEipsArgs) ([]Eip, error) {
	if args == nil {
		args = &GetEipsArgs{}
	}
	params := map[string]string{
		"eip":          args.Ip,
		"instanceType": string(args.InstanceType),
		"instanceId":   args.InstanceId,
	}
	req, err := bce.NewRequest("GET", c.GetURL("v1/eip", params), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.SendRequest(req, nil)
	if err != nil {
		return nil, err
	}
	bodyContent, err := resp.GetBodyContent()

	if err != nil {
		return nil, err
	}
	var blbsResp *GetEipsResponse
	err = json.Unmarshal(bodyContent, &blbsResp)

	if err != nil {
		return nil, err
	}
	return blbsResp.EipList, nil
}

// not implemented
func (c *Client) PurchaseReservedEips() {

}
