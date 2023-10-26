package statemachine

import (
	"errors"
	"fmt"
	"sync"
)

type EventType string

const (
	OnErrorEvent = EventType("error")
)

var (
	ErrorEventDoesNotExist = "event type does not exist"
	ErrorCbExist           = "callback already registered"
)

type Event struct {
	Type EventType
	Data interface{}
	Err  error
}

type Action interface {
	Execute(event *Event) (interface{}, error)
	Type() EventType
}

type StateMachine struct {
	callbacks    map[EventType]Callback
	inputChannel chan Event
	errorChannel chan Event
	doneChannel  chan bool
	mutex        sync.Mutex
}

type Callback func(event Event) error

type IStateMachine interface {
	Close() error
	Send(event *Event) error
	On(eventType EventType, cb Callback) error
}

func NewGameMachine() IStateMachine {
	machine := &StateMachine{
		inputChannel: make(chan Event),
		callbacks:    map[EventType]Callback{},
		errorChannel: make(chan Event),
		doneChannel:  make(chan bool),
	}

	go func() {
		for i := 0; i < 2; i++ {
			select {
			case done := <-machine.doneChannel:
				fmt.Println("closing game machine: ", done)
				return
			case event := <-machine.errorChannel:
				err := machine.processEvent(event)
				fmt.Println(err)
			case event := <-machine.inputChannel:
				err := machine.processEvent(event)
				if err != nil {
					event.Err = err
					machine.errorChannel <- event
				}
			}
		}
	}()

	return machine
}

func (machine *StateMachine) processEvent(event Event) error {
	machine.mutex.Lock()
	defer machine.mutex.Unlock()
	if event.Err != nil {
		cb, ok := machine.callbacks[OnErrorEvent]
		if !ok {
			return errors.New(ErrorEventDoesNotExist)
		}
		err := cb(event)
		if err != nil {
			return err
		}
	}

	cb, ok := machine.callbacks[event.Type]
	if !ok {
		return errors.New(ErrorEventDoesNotExist)
	}

	cb(event)

	return nil
}

func (machine *StateMachine) Close() error {
	machine.doneChannel <- true
	return nil
}

func (machine *StateMachine) Send(event *Event) error {
	machine.mutex.Lock()
	defer machine.mutex.Unlock()

	_, ok := machine.callbacks[event.Type]
	if !ok {
		return errors.New(ErrorEventDoesNotExist)
	}

	machine.inputChannel <- *event

	return nil
}

func (machine *StateMachine) On(eventType EventType, cb Callback) error {
	_, ok := machine.callbacks[eventType]
	if ok {
		return fmt.Errorf("%s: %s", ErrorCbExist, eventType)
	}

	machine.callbacks[eventType] = cb

	return nil
}
