package actor

import (
	"fmt"
)

type VM struct {
	name string
	host string
	port int
}

type Pid struct {
	id int64
	//vm *VM
}

type Message struct {
	destination *Pid
	origin      *Pid
	mode        string //@TODO: think about! Synch/ASynch no-response
	cmd         string
	payload     interface{} //used to carry response content
}

//Command definition to be executed by the actor
type Command func(state interface{}) (result interface{}, err error)

// Actor type Interface
type ActorInterface interface {
	RegisterCmd(string, Command)
	StartLink(state interface{}) *Pid  //@TODO: handle link
	Send(Message) (interface{}, error) //manipulates state and get response
	Terminate() error

	//Internals
	loop()
}

type Actor struct {
	ActorInterface

	pid           *Pid
	state         interface{}
	commands      map[string]Command
	receiveChan   chan (*Message)
	sendChan      chan (*Message)
	terminateChan chan (bool)
	alive         bool
}

func GetActor(sndCh, rcvCh chan (*Message), exitCH chan (bool)) *Actor {
	return &Actor{
		pid:           getPid(),
		alive:         true,
		sendChan:      sndCh,  // make(chan (*Message), 10),
		receiveChan:   rcvCh,  // make(chan (*Message), 10),
		terminateChan: exitCH, // make(chan (bool), 1),
	}
}

// RegisterCmd register command on actor loop execution
func (a *Actor) RegisterCmd(cmd string, f Command) {
	a.commands[cmd] = f
}

func (a *Actor) StartLink(state interface{}) *Pid {
	a.state = state
	go a.loop()

	return a.pid
}

//must be public to override original init ?¿?¿
func (a *Actor) loop() {
	for a.alive {
		select {
		case msg := <-a.receiveChan:
			fmt.Printf("pid %s has received %s \n", a.pid, msg.cmd)
			if f, ok := a.commands[msg.cmd]; ok {
				if msg.mode == "synch" {
					result, err := f(a.state)
					if err != nil {
						fmt.Printf("pid %s error executing command %s \n", msg.destination, msg.cmd)
						continue
					}
					a.sendChan <- createResponse(msg, result) // @TODO: maybe here will be better to create some MessageResponse
				} else {
					_, err := f(a.state)
					if err != nil {
						fmt.Printf("pid %s error executing command %s \n", msg.destination, msg.cmd)
						continue
					}
				}
				//@TODO: SURE???
				//				if msg.mode == "synch" { // @Or typecasting?
				//				} else {
				//					err := f(a.state)
				//				}

			} else {
				fmt.Printf("pid %s has Not registered cmd %s \n", a.pid, msg.cmd)
			}
		case _ = <-a.terminateChan:
			fmt.Printf("pid %s finish \n", a.pid)
			a.alive = false
			return
		}
	}
}

func (a *Actor) Terminate() error {
	a.terminateChan <- true

	return nil
}

//@TODO: needs to handle Pid generation...solve it or shame on that!!
var lastPid = int64(1)

func getPid() (p *Pid) {
	p = &Pid{id: lastPid}
	lastPid++

	return
}

func createResponse(msg *Message, result interface{}) *Message {
	return &Message{
		destination: msg.origin,
		origin:      msg.destination,
		mode:        "response", //@TODO: think about! Synch/ASynch no-response response!
		payload:     result,
	}
}
