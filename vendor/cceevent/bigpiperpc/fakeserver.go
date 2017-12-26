package bigpiperpc

import (
	"errors"
	"fmt"
	"testing"

	"github.com/baidu-golang/pbrpc"
	"github.com/golang/protobuf/proto"
)

var (
	FakeBigConfig = BigpipeConfig{
		Pipe:     "testpipe",
		AclName:  "testacl",
		AclToken: "testtoken",
	}
)

func DoRpcServerStart(t *testing.T, port int) *pbrpc.TcpServer {

	rpcServer := createRpcServer(port)

	err := rpcServer.Start()
	if err != nil && t != nil {
		t.Error(err)
	}
	return rpcServer
}

func createRpcServer(port int) *pbrpc.TcpServer {
	serverMeta := pbrpc.ServerMeta{}
	serverMeta.Host = proto.String("localhost")
	serverMeta.Port = Int(port)
	rpcServer := pbrpc.NewTpcServer(&serverMeta)

	ss := NewSimpleService("ProxyService", "SendBp")
	rpcServer.Register(ss)
	return rpcServer
}

func Int(v int) *int {
	p := new(int)
	*p = int(v)
	return p
}

type SimpleService struct {
	serviceName string
	methodName  string
}

func (ss *SimpleService) GetServiceName() string {
	return ss.serviceName
}

func (ss *SimpleService) GetMethodName() string {
	return ss.methodName
}

func (ss *SimpleService) NewParameter() proto.Message {
	ret := BpEntry{}
	return &ret
}

func NewSimpleService(serviceName, methodName string) *SimpleService {
	ret := SimpleService{serviceName, methodName}
	return &ret
}

func (ss *SimpleService) DoService(msg proto.Message, attachment []byte, logId *int64) (proto.Message, []byte, error) {
	dmRes := BpResponse{}
	if msg != nil {
		m, ok := msg.(*BpEntry)
		if !ok {
			errStr := "message type is not type of 'DataMessage'"
			return nil, nil, errors.New(errStr)
		}
		if *m.AclName != FakeBigConfig.AclName || *m.AclToken != FakeBigConfig.AclToken {
			return nil, nil, fmt.Errorf("AclName error")
		}
	}

	dmRes.Errmsg = proto.String("nil")
	dmRes.Status = proto.Int32(109)
	dmRes.PipeletMsgId = proto.Uint64(123)
	return &dmRes, nil, nil
}
