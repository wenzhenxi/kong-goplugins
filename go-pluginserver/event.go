package main

import (
	"fmt"
	"github.com/sunmi-OS/go-pdk"
	"time"
)

// Incoming data for a new event.
// TODO: add some relevant data to reduce number of callbacks.
type StartEventData struct {
	InstanceId int    // Instance ID to start the event
	EventName  string // event name (not handler method name)
	// ....
}

type eventData struct {
	id       int              // event id
	instance *instanceData    // plugin instance
	ipc      chan interface{} // communication channel (TODO: use decoded structs)
	pdk      *pdk.PDK         // go-pdk instance
}

// HandleEvent starts the call/{callback/response}*/finish cycle.
// More than one event can be run concurrenty for a single plugin instance,
// they all receive the same object instance, so should be careful if it's
// mutated or holds references to mutable data.
//
// RPC exported method

// HandleEvent开始呼叫/ {callback / response} * /结束周期。
//一个插件实例可以同时运行多个事件，
//它们都接收相同的对象实例，因此如果
//变异或保留对可变数据的引用。
//
// RPC导出方法
func (s *PluginServer) HandleEvent(in StartEventData, out *StepData) error {
	s.lock.RLock()
	instance, ok := s.instances[in.InstanceId]
	s.lock.RUnlock()
	if !ok {
		return fmt.Errorf("No plugin instance %d", in.InstanceId)
	}

	h, ok := instance.handlers[in.EventName]
	if !ok {
		return fmt.Errorf("undefined method %s on plugin %s",
			in.EventName, instance.plugin.name)
	}

	ipc := make(chan interface{})

	event := eventData{
		instance: instance,
		ipc:      ipc,
		pdk:      pdk.Init(ipc),
	}

	s.lock.Lock()
	event.id = s.nextEventId
	s.nextEventId++
	s.events[event.id] = &event
	s.lock.Unlock()

	//log.Printf("Will launch goroutine for key %d / operation %s\n", key, op)
	go func() {
		_ = <-ipc
		h(event.pdk)

		func() {
			defer func() { recover() }()
			ipc <- "ret"
		}()

		s.lock.Lock()
		defer s.lock.Unlock()
		event.instance.lastEvent = time.Now()
		delete(s.events, event.id)
	}()

	// ipc 当前为2
	ipc <- "run" // kickstart the handler

	// 这条先与上面的go执行 ipc拿到了内容
	*out = StepData{EventId: event.id, Data: <-ipc}
	return nil
}

// A callback's response/request.
type StepData struct {
	EventId int         // event cycle to which this belongs
	Data    interface{} // carried data
}

// Step carries a callback's anser back from Kong to the plugin,
// the return value is either a new callback request or a finish signal.
//
// RPC exported method
//步骤将回调的分析服务从Kong返回到插件，
//返回值是新的回调请求或完成信号。
//
// RPC导出方法
func (s *PluginServer) Step(in StepData, out *StepData) error {
	s.lock.RLock()
	event, ok := s.events[in.EventId]
	s.lock.RUnlock()
	if !ok {
		return fmt.Errorf("No running event %d", in.EventId)
	}

	event.ipc <- in.Data
	outStr := <-event.ipc
	*out = StepData{EventId: in.EventId, Data: outStr}

	return nil
}
