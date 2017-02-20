package scheduler

/*
Scheduler struct defines an interface of the default calls for the scheduler, as well as holding information
regarding the framework, client, and events handler.

End users may wish to make their own scheduler and satisfy this interface.

A default scheduler is provided for those who wish to keep the default implementations: All calls simply create
the protobuf required for the call and send it off to the client.

End users should only create their own scheduler if they wish to change the behavior of their calls.
*/
import (
	"fmt"

	"io/ioutil"
	"log"
	"mesos-framework-sdk/client"
	mesos "mesos-framework-sdk/include/mesos"
	sched "mesos-framework-sdk/include/scheduler"
	"mesos-framework-sdk/recordio"
	"time"
)

const (
	subscribeRetry = 2
)

type Scheduler interface {
	// Scheduler must also hold framework information, a client and an event handler.
	FrameworkInfo() *mesos.FrameworkInfo
	Client() *client.Client

	// Default Calls for scheduler
	Subscribe(chan *sched.Event) error
	Teardown()
	Accept(offerIds []*mesos.OfferID, tasks []*mesos.Offer_Operation, filters *mesos.Filters)
	Decline(offerIds []*mesos.OfferID, filters *mesos.Filters)
	Revive()
	Kill(taskId *mesos.TaskID, agentid *mesos.AgentID)
	Shutdown(execId *mesos.ExecutorID, agentId *mesos.AgentID)
	Acknowledge(agentId *mesos.AgentID, taskId *mesos.TaskID, uuid []byte)
	Reconcile(tasks []*mesos.Task)
	Message(agentId *mesos.AgentID, executorId *mesos.ExecutorID, data []byte)
	SchedRequest(resources []*mesos.Request)
	Suppress()
}

// Default Scheduler can be used as a higher-level construct.
type DefaultScheduler struct {
	FramworkInfo *mesos.FrameworkInfo
	client       *client.Client
}

func NewDefaultScheduler(c *client.Client, info *mesos.FrameworkInfo) *DefaultScheduler {
	return &DefaultScheduler{
		client:       c,
		FramworkInfo: info,
	}
}

func (c *DefaultScheduler) FrameworkInfo() *mesos.FrameworkInfo {
	return c.FramworkInfo
}

// Make a subscription call to mesos.
// Channel passed is the "listener" channel for Event Controller.
func (c *DefaultScheduler) Subscribe(eventChan chan *sched.Event) error {
	call := &sched.Call{
		Type: sched.Call_SUBSCRIBE.Enum(),
		Subscribe: &sched.Call_Subscribe{
			FrameworkInfo: c.FramworkInfo,
		},
	}
	go func() {
		for {
			resp, err := c.client.Request(call)
			if err != nil {
				log.Println(err.Error())
			} else {
				log.Println(recordio.Decode(resp.Body, eventChan))
			}

			// If we disconnect we need to reset the stream ID.
			// Otherwise we'll never be able to reconnect.
			c.client.StreamID = ""
			time.Sleep(time.Duration(subscribeRetry) * time.Second)
		}
	}()
	return nil
}

// Send a teardown request to mesos master.
func (c *DefaultScheduler) Teardown() {
	teardown := &sched.Call{
		FrameworkId: c.FramworkInfo.GetId(),
		Type:        sched.Call_TEARDOWN.Enum(),
	}
	resp, err := c.client.Request(teardown)
	if err != nil {
		log.Println(err.Error())
	}

	fmt.Println(resp)
}

// Accepts offers from mesos master
func (c *DefaultScheduler) Accept(offerIds []*mesos.OfferID, tasks []*mesos.Offer_Operation, filters *mesos.Filters) {
	accept := &sched.Call{
		FrameworkId: c.FramworkInfo.GetId(),
		Type:        sched.Call_ACCEPT.Enum(),
		Accept:      &sched.Call_Accept{OfferIds: offerIds, Operations: tasks, Filters: filters},
	}

	resp, err := c.client.Request(accept)
	if err != nil {
		log.Println(err.Error())
		return
	}
	// NOTE: If we get back an improper response, this will panic.
	// TODO: recover from panics.
	k, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(k))

}

func (c *DefaultScheduler) Decline(offerIds []*mesos.OfferID, filters *mesos.Filters) {
	// Get a list of the offer ids to decline and any filters.
	decline := &sched.Call{
		FrameworkId: c.FramworkInfo.GetId(),
		Type:        sched.Call_DECLINE.Enum(),
		Decline:     &sched.Call_Decline{OfferIds: offerIds, Filters: filters},
	}

	resp, err := c.client.Request(decline)
	if err != nil {
		log.Println(err.Error())
	}
	a, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(a))
	return
}

// Sent by the scheduler to remove any/all filters that it has previously set via ACCEPT or DECLINE calls.
func (c *DefaultScheduler) Revive() {
	revive := &sched.Call{
		FrameworkId: c.FramworkInfo.GetId(),
		Type:        sched.Call_REVIVE.Enum(),
	}

	resp, err := c.client.Request(revive)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
	return
}

func (c *DefaultScheduler) Kill(taskId *mesos.TaskID, agentid *mesos.AgentID) {
	// Probably want some validation that this is a valid task and valid agentid.
	kill := &sched.Call{
		FrameworkId: c.FramworkInfo.GetId(),
		Type:        sched.Call_KILL.Enum(),
		Kill:        &sched.Call_Kill{TaskId: taskId, AgentId: agentid},
	}

	resp, err := c.client.Request(kill)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
	return
}

func (c *DefaultScheduler) Shutdown(execId *mesos.ExecutorID, agentId *mesos.AgentID) {
	shutdown := &sched.Call{
		FrameworkId: c.FramworkInfo.GetId(),
		Type:        sched.Call_SHUTDOWN.Enum(),
		Shutdown: &sched.Call_Shutdown{
			ExecutorId: execId,
			AgentId:    agentId,
		},
	}
	resp, err := c.client.Request(shutdown)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
	return
}

// UUID should be a type
// TODO import extras uuid funcs.
func (c *DefaultScheduler) Acknowledge(agentId *mesos.AgentID, taskId *mesos.TaskID, uuid []byte) {
	fmt.Println("Acknowledge event.")
	acknowledge := &sched.Call{
		FrameworkId: c.FramworkInfo.GetId(),
		Type:        sched.Call_ACKNOWLEDGE.Enum(),
		Acknowledge: &sched.Call_Acknowledge{
			AgentId: agentId,
			TaskId:  taskId,
			Uuid:    uuid,
		},
	}
	resp, err := c.client.Request(acknowledge)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
}

func (c *DefaultScheduler) Reconcile(tasks []*mesos.Task) {
	reconcile := &sched.Call{
		FrameworkId: c.FramworkInfo.GetId(),
		Type:        sched.Call_RECONCILE.Enum(),
	}
	resp, err := c.client.Request(reconcile)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
}

func (c *DefaultScheduler) Message(agentId *mesos.AgentID, executorId *mesos.ExecutorID, data []byte) {
	message := &sched.Call{
		FrameworkId: c.FramworkInfo.GetId(),
		Type:        sched.Call_MESSAGE.Enum(),
		Message: &sched.Call_Message{
			AgentId:    agentId,
			ExecutorId: executorId,
			Data:       data,
		},
	}
	resp, err := c.client.Request(message)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)

}

// Sent by the scheduler to request resources from the master/allocator.
// The built-in hierarchical allocator simply ignores this request but other allocators (modules) can interpret this in
// a customizable fashion.
func (c *DefaultScheduler) SchedRequest(resources []*mesos.Request) {
	request := &sched.Call{
		FrameworkId: c.FramworkInfo.GetId(),
		Type:        sched.Call_REQUEST.Enum(),
		Request: &sched.Call_Request{
			Requests: resources,
		},
	}
	resp, err := c.client.Request(request)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
}

func (c *DefaultScheduler) Suppress() {
	supress := &sched.Call{
		FrameworkId: c.FramworkInfo.GetId(),
		Type:        sched.Call_SUPPRESS.Enum(),
	}
	resp, err := c.client.Request(supress)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
}
