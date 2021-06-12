package EventBus

import (
	"fmt"
	"reflect"
	"sync"
)

type BusSubscriber interface {
	Subscribe(topic string, fn interface{}) error
	Unsubscribe(topic string, handler interface{}) error
}

type BusPublisher interface {
	Publish(topic string, args ...interface{})
}

type BusController interface {
	HasCallback(topic string) bool
	WaitAsync()
}

type Bus interface {
	BusController
	BusSubscriber
	BusPublisher
}

type EventBus struct {
	handlers map[string][]*eventHandler
	lock     sync.Mutex
	wg       sync.WaitGroup
}

type eventHandler struct {
	callBack      reflect.Value
	flagOnce      bool
	async         bool
	transactional bool
	sync.Mutex
}

func New() Bus {
	b := &EventBus{
		make(map[string][]*eventHandler),
		sync.Mutex{},
		sync.WaitGroup{},
	}
	return Bus(b)
}

func (bus *EventBus) Subscribe(topic string, fn interface{}) error {
	return bus.doSubscribe(topic, fn, &eventHandler{
		reflect.ValueOf(fn), false, false, false, sync.Mutex{},
	})
}

func (bus *EventBus) doSubscribe(topic string, fn interface{}, handler *eventHandler) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	if !(reflect.TypeOf(fn).Kind() == reflect.Func) {
		return fmt.Errorf("$s is not of type reflect.Func", reflect.TypeOf(fn).Kind())
	}
	bus.handlers[topic] = append(bus.handlers[topic], handler)
	return nil
}

func (bus *EventBus) HasCallBack(topic string) bool {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	_, ok := bus.handlers[topic]
	if ok {
		return len(bus.handlers[topic]) > 0
	}
	return false
}
