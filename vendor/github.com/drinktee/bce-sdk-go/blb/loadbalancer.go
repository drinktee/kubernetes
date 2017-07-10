package blb

import (
	"bytes"
	"encoding/json"

	"fmt"

	"github.com/drinktee/bce-sdk-go/bce"
	"github.com/drinktee/bce-sdk-go/util"
)

type LoadBalancer struct {
	BlbId    string `json:"blbId"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Address  string `json:"address"`
	Status   string `json:"status"`
	PublicIp string `json:"publicIp"`
}

type DescribeLoadBalancersArgs struct {
	LoadBalancerId   string
	LoadBalancerName string
	BCCId            string
	Address          string
}

type DescribeLoadBalancersResponse struct {
	Marker      string         `json:"marker"`
	IsTruncated bool           `json:"isTruncated"`
	NextMarker  string         `json:"nextMarker"`
	MaxKeys     int            `json:"maxKeys"`
	BLBList     []LoadBalancer `json:"blbList"`
}

// CreateLoadBalancerArgs create blb args
type CreateLoadBalancerArgs struct {
	Desc  string `json:"desc"`
	Name  string `json:"name"`
	VpcID string `json:"vpcId,omitempty"`
}

type CreateLoadBalancerResponse struct {
	LoadBalancerId string `json:"blbId"`
	Address        string `json:"address"`
	Desc           string `json:"desc,omitempty"`
	Name           string `json:"name"`
}

type UpdateLoadBalancerArgs struct {
	LoadBalancerId string `json:"blbId"`
	Desc           string `json:"desc,omitempty"`
	Name           string `json:"name,omitempty"`
}

type DeleteLoadBalancerArgs struct {
	LoadBalancerId string `json:"blbId"`
}

// DescribeLoadBalancers Describe loadbalancers
// TODO: args need to validate
func (c *Client) DescribeLoadBalancers(args *DescribeLoadBalancersArgs) ([]LoadBalancer, error) {
	var params map[string]string
	if args != nil {
		params = map[string]string{
			"blbId":   args.LoadBalancerId,
			"name":    args.LoadBalancerName,
			"bccId":   args.BCCId,
			"address": args.Address,
		}
	}
	req, err := bce.NewRequest("GET", c.GetURL("v1/blb", params), nil)

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
	var blbsResp *DescribeLoadBalancersResponse
	err = json.Unmarshal(bodyContent, &blbsResp)

	if err != nil {
		return nil, err
	}
	return blbsResp.BLBList, nil
}

// CreateLoadBalancer Create a  loadbalancer
// TODO: args need to validate
func (c *Client) CreateLoadBalancer(args *CreateLoadBalancerArgs) (*CreateLoadBalancerResponse, error) {
	var params map[string]string
	if args != nil {
		params = map[string]string{
			"clientToken": c.GenerateClientToken(),
		}
	}
	postContent, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}
	req, err := bce.NewRequest("POST", c.GetURL("v1/blb", params), bytes.NewBuffer(postContent))
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
	var blbsResp *CreateLoadBalancerResponse
	err = json.Unmarshal(bodyContent, &blbsResp)

	if err != nil {
		return nil, err
	}
	return blbsResp, nil
}

// UpdateLoadBalancer update a loadbalancer
// TODO: args need to validate
func (c *Client) UpdateLoadBalancer(args *UpdateLoadBalancerArgs) error {
	params := map[string]string{
		"clientToken": c.GenerateClientToken(),
	}
	if args == nil {
		return fmt.Errorf("UpdateLoadBalancer need args")
	}
	postContent, err := util.ToJson(args, "desc", "name")
	// postContent, err := json.Marshal(args)
	if err != nil {
		return err
	}
	req, err := bce.NewRequest("PUT", c.GetURL("v1/blb"+"/"+args.LoadBalancerId, params), bytes.NewBuffer(postContent))
	if err != nil {
		return err
	}
	_, err = c.SendRequest(req, nil)
	if err != nil {
		return err
	}
	return nil
}

// DeleteLoadBalancer delete a loadbalancer
func (c *Client) DeleteLoadBalancer(args *DeleteLoadBalancerArgs) error {
	params := map[string]string{
		"clientToken": c.GenerateClientToken(),
	}
	if args == nil {
		return fmt.Errorf("DeleteLoadBalancer need args")
	}
	req, err := bce.NewRequest("DELETE", c.GetURL("v1/blb"+"/"+args.LoadBalancerId, params), nil)
	if err != nil {
		return err
	}
	_, err = c.SendRequest(req, nil)
	if err != nil {
		return err
	}
	return nil
}
