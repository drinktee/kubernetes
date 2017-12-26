package bigpiperpc

import (
	"fmt"
	"time"

	"github.com/baidu-golang/pbrpc"
	"github.com/golang/protobuf/proto"
)

// Client define fiber rpc client
type Client interface {
	SendBp(entry *BpEntry) (*BpResponse, error)
}

var (
	ServiceName string = "ProxyService"
	MethodName  string = "SendBp"
)

// BigpipeConfig define configs of bigpipe
type BigpipeConfig struct {
	Pipe      string `json:"Pipe"`
	AclName   string `json:"AclName"`
	AclPasswd string `json:"AclPasswd"`
	AclToken  string `json:"AclToken"`
	PipeletId uint32 `json:"PipeletId"`
}

type FiberProxyClient struct {
	Config           *BigpipeConfig
	Url              pbrpc.URL
	Timeout          int
	SendBpInvocation *pbrpc.RpcInvocation
}

// NewFiberProxyClient create a client
func NewFiberProxyClient(config *BigpipeConfig, host string, port int, timeout int) *FiberProxyClient {
	rpcInvocation := pbrpc.NewRpcInvocation(&ServiceName, &MethodName)
	rpcInvocation.CompressType = proto.Int32(pbrpc.COMPRESS_GZIP)
	url := pbrpc.URL{}
	url.SetHost(&host).SetPort(&port)
	return &FiberProxyClient{
		Config:           config,
		Timeout:          timeout,
		SendBpInvocation: rpcInvocation,
		Url:              url,
	}
}

// SendBp send a BpEntry
func (fc *FiberProxyClient) SendBp(entry *BpEntry) (*BpResponse, error) {
	timeout := time.Second * time.Duration(fc.Timeout)
	connection, err := pbrpc.NewTCPConnection(fc.Url, &timeout)
	if err != nil {
		return nil, fmt.Errorf("NewTCPConnection error:%v", err)
	}
	defer connection.Close()
	rpcClient, err := pbrpc.NewRpcCient(connection)
	if err != nil {
		return nil, fmt.Errorf("NewRpcCient error:%v", err)
	}
	fc.SendBpInvocation.SetParameterIn(entry)
	fc.SendBpInvocation.LogId = proto.Int64(1)
	parameterOut := BpResponse{}
	_, err = rpcClient.SendRpcRequest(fc.SendBpInvocation, &parameterOut)
	if err != nil {
		return nil, fmt.Errorf("SendRpcRequest error:%v", err)
	}
	return &parameterOut, nil
}
