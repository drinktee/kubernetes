package event

// CceEvent define cce events
type CceEvent struct {
	EventTimestamp  string `json:"eventTimestamp"`  //事件发生时间戳
	InvolvedCluster string `json:"involvedCluster"` //关联集群c-1xxxxx
	ResourceType    string `json:"resourceType"`    //资源类型BLB
	ResourceName    string `json:"resourceName"`    //资源名称LBxxxxx
	ResourceID      string `json:"resourceID"`      //lb-xxxxx
	EventSource     string `json:"eventSource"`     //事件源
	Desc            string `json:"desc"`            //事件描述创建BLB
}

//  resources
const (
	BlbCceResource   string = "BLB"
	BccCceResource   string = "BCC"
	RouteCceResource string = "Route"
	EIPCceResource   string = "EIP"

	EventSourceKubernetes  string = "kubernetes"
	EventSourceCaaSManager string = "caas-api-manager"
	EventSourceCceService  string = "cce-service"
)

const (
	// CceEventBegin is header of message
	CceEventBegin string = "CCE_EVENT_BEGIN"
	// CceEventEnd is end of message
	CceEventEnd string = "CCE_EVENT_END"
)
