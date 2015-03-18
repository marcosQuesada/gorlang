package actor

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

var (
	ErrCmdNotFound = errors.New("Command Not found")
	ErrTimeout     = errors.New("Timeout command")
)

const (
	RcvChBuffer  = 10
	SndChBuffer  = 1
	ExitChBuffer = 1
)

type VM struct {
	name string
	host string
	port int
}

type Pid int64

//Used to store Pid assignation
var lastPid Pid = 1

type Message struct {
	destination Pid
	origin      Pid
	cmd         string
	attr        interface{}
}

type ComplexResponse struct {
	Response interface{}
	Err      error
}

//Command definition to be executed by the actor
type Command func(state interface{}) (result interface{}, err error)

// Actor type Interface
type ActorInterface interface {
	Spawn(state interface{}) Pid //@TODO: handle link
	Send(*Pid, string, ...interface{}) (interface{}, error)
	Terminate() error
}

type Actor struct {
	ActorInterface

	pid         Pid
	state       interface{}
	commands    map[string]Command
	receiveChan chan (*Message)
	sendChan    chan (interface{})
	ExitChan    chan (bool)
	alive       bool
	methodIndex map[string]reflect.Value
}

func NewActor() *Actor {
	return &Actor{
		pid:         getPid(),
		alive:       true,
		receiveChan: make(chan (*Message), RcvChBuffer),
		sendChan:    make(chan (interface{}), SndChBuffer),
		ExitChan:    make(chan (bool), ExitChBuffer),
	}
}

func (a *Actor) Spawn(state interface{}) Pid {
	a.state = state
	a.loadMethods(state)
	go a.loop()

	return a.pid
}

// Send message to destPid , define cmd and attributes
func (a *Actor) Send(destPid Pid, cmd string, attr ...interface{}) (interface{}, error) {
	m := &Message{
		destination: destPid,
		cmd:         cmd,
	}
	m.SetAttrs(attr...)
	a.receiveChan <- m
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(100 * time.Millisecond)
		timeout <- true
	}()
	select {
	case <-timeout:
		return nil, ErrTimeout
	case result := <-a.sendChan:
		switch result.(type) {
		case ComplexResponse:
			cast := result.(ComplexResponse)
			return cast.Response, cast.Err
		default:
			return result, nil
		}
	}
}

//Terminate Actor
func (a *Actor) Terminate() error {
	a.alive = false
	return nil
}

func (a *Actor) GetPid() Pid {
	return a.pid
}

func (a *Actor) loop() {
	defer fmt.Println("closing loopppp")
	defer a.ExitSignal()
	defer close(a.sendChan)
	defer close(a.receiveChan)
	//defer close(a.ExitChan)

	for a.alive {
		select {
		case msg := <-a.receiveChan:
			if _, ok := a.methodIndex[msg.cmd]; !ok {
				a.sendChan <- createComplexResponse(nil, ErrCmdNotFound)
				continue
			}

			//@TODO: needs to handle list types
			in := make([]reflect.Value, 0)
			for _, v := range msg.attr.([]interface{}) {
				in = append(in, reflect.ValueOf(v))
			}

			result := a.methodIndex[msg.cmd].Call(in)
			if len(result) > 0 {
				if len(result) > 1 {
					v := cast(result[0])
					var err error
					if result[1].IsNil() {
						err = nil
					} else {
						err = cast(result[1]).(error)
					}
					a.sendChan <- createComplexResponse(v, err)
					continue
				}
				a.sendChan <- cast(result[0])

			} else {
				a.sendChan <- nil
			}
		}
	}
}

func (a *Actor) ExitSignal() {
	fmt.Println("Out of the loop")
	a.ExitChan <- true
	fmt.Println("Sended exit signal from ", a.GetPid())
}

func (a *Actor) loadMethods(state interface{}) {
	stateType := reflect.TypeOf(state)
	stateValue := reflect.ValueOf(state)

	index := make(map[string]reflect.Value, stateType.NumMethod())

	for i := 0; i < stateType.NumMethod(); i++ {
		methodType := stateType.Method(i)
		methodValue := stateValue.MethodByName(methodType.Name)
		index[methodType.Name] = methodValue
	}
	a.methodIndex = index
}

func (m *Message) SetAttrs(attrs ...interface{}) {
	m.attr = attrs
}

func getPid() (p Pid) {
	p = Pid(lastPid)
	lastPid++

	return
}

func createComplexResponse(c interface{}, err error) ComplexResponse {
	return ComplexResponse{
		Response: c,
		Err:      err,
	}
}

func cast(result reflect.Value) interface{} {
	switch result.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return result.Int()
	case reflect.String:
		return result.String()
	case reflect.Map:
		//fmt.Println("Reflect Map")
	case reflect.Slice:
		//fmt.Println("Reflect SLice")
	case reflect.Array:
		//fmt.Println("Reflect Array")
	case reflect.Struct:
		//fmt.Println("Reflect Struct")
	case reflect.Ptr:
		//fmt.Println("Reflect Ptr", result)
		return result.Interface()

	}
	return errors.New(result.String())
}
