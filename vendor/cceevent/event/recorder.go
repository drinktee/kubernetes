package event

import (
	"cceevent/bigpiperpc"
	utilruntime "cceevent/runtime"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/satori/go.uuid"

	"github.com/golang/protobuf/proto"
)

// Recorder knows how to record events on behalf of an EventSource.
type Recorder interface {
	// Event constructs an event from the given information and puts it in the queue for sending.
	Event(involvedCluster, resourceType, resourceID, resourceName, source string, desc string) error
}

// SimpleRecorder define a simple recorder
type SimpleRecorder struct {
	FiberClient bigpiperpc.Client
	Cfg         *bigpiperpc.BigpipeConfig
}

// NewSimpleRecorder create a SimpleRecorder
func NewSimpleRecorder(host string, port int, timeout int, cfg *bigpiperpc.BigpipeConfig) *SimpleRecorder {
	fc := bigpiperpc.NewFiberProxyClient(cfg, host, port, timeout)
	return &SimpleRecorder{
		FiberClient: fc,
		Cfg:         cfg,
	}
}

// Event sent an CceEvent
func (recorder *SimpleRecorder) Event(involvedCluster, resourceType, resourceID, resourceName, source, desc string) error {
	if involvedCluster == "" || source == "" || desc == "" {
		return fmt.Errorf("Event involvedCluster,source,desc shoud not nil")
	}
	event := CceEvent{
		EventTimestamp:  fmt.Sprintf("%d", time.Now().Unix()),
		InvolvedCluster: involvedCluster,
		ResourceType:    resourceType,
		ResourceName:    resourceName,
		ResourceID:      resourceID,
		EventSource:     source,
		Desc:            desc,
	}
	entry, err := generateBpFromEvent(event, recorder.Cfg)
	if err != nil {
		return err
	}
	res, err := recorder.FiberClient.SendBp(entry)
	if err != nil {
		return err
	}
	glog.V(4).Infof("Send CceEvent success: %+v  SendBp response: %v", event, res)
	return err
}

func generateBpFromEvent(event CceEvent, cfg *bigpiperpc.BigpipeConfig) (*bigpiperpc.BpEntry, error) {
	mess, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	now := time.Now().Unix() + 5
	before := fmt.Sprintf("%d%s", now, cfg.AclToken)
	token := md5.Sum([]byte(before))
	entry := &bigpiperpc.BpEntry{
		Pipe:      proto.String(cfg.Pipe),
		AclName:   proto.String(cfg.AclName),
		AclToken:  proto.String(fmt.Sprintf("%x", token)),
		PipeletId: proto.Uint32(cfg.PipeletId),
		Message:   []byte(CceEventBegin + string(mess) + CceEventEnd),
		Expired:   proto.Int64(now),
		MsgGuid:   proto.String(fmt.Sprintf("%s%s%s", cfg.Pipe, cfg.AclName, uuid.NewV4().String())),
	}
	glog.V(4).Infof("Generate bpentry: %s", entry.String())
	return entry, nil
}

const incomingQueueLength = 10

// QueueRecorder define a recorder with channel
type QueueRecorder struct {
	lock sync.Mutex

	distributing sync.WaitGroup
	FiberClient  bigpiperpc.Client
	Cfg          *bigpiperpc.BigpipeConfig
	incoming     chan CceEvent
}

// NewQueueRecorder create a QueueRecorder
func NewQueueRecorder(host string, port int, timeout int, cfg *bigpiperpc.BigpipeConfig, fc bigpiperpc.Client) *QueueRecorder {
	if fc == nil {
		fc = bigpiperpc.NewFiberProxyClient(cfg, host, port, timeout)
	}
	q := &QueueRecorder{
		incoming:    make(chan CceEvent, incomingQueueLength),
		FiberClient: fc,
		Cfg:         cfg,
	}
	q.distributing.Add(1)
	go q.loop()
	return q
}

// loop receives from m.incoming and distributes to all watchers.
func (q *QueueRecorder) loop() {
	for {
		event, ok := <-q.incoming
		if !ok {
			break
		}
		q.distribute(event)
	}
	q.distributing.Done()
}

// distribute sends event to all watchers. Blocking.
func (q *QueueRecorder) distribute(event CceEvent) {
	q.lock.Lock()
	defer q.lock.Unlock()
	entry, err := generateBpFromEvent(event, q.Cfg)
	if err != nil {
		glog.Errorf("QueueRecorder generateBpFromEvent error: %v", err)
		return
	}
	res, err := q.FiberClient.SendBp(entry)
	if err != nil {
		glog.Errorf("QueueRecorder use FiberClient SendBp error: %v", err)
		return
	}
	glog.V(4).Infof("QueueRecorder Send CceEvent success: %+v  SendBp response: %v", event, res)
}

// Shutdown close incoming
func (q *QueueRecorder) Shutdown() {
	close(q.incoming)
	q.distributing.Wait()
}

// Event sent an CceEvent
func (q *QueueRecorder) Event(involvedCluster, resourceType, resourceID, resourceName, source, desc string) error {
	if involvedCluster == "" || source == "" || desc == "" {
		return fmt.Errorf("Event involvedCluster,source,desc shoud not nil")
	}
	event := CceEvent{
		EventTimestamp:  fmt.Sprintf("%d", time.Now().Unix()),
		InvolvedCluster: involvedCluster,
		ResourceType:    resourceType,
		ResourceName:    resourceName,
		ResourceID:      resourceID,
		EventSource:     source,
		Desc:            desc,
	}
	// NOTE: events should be a non-blocking operation
	go func() {
		defer utilruntime.HandleCrash()
		q.incoming <- event
	}()
	return nil
}
