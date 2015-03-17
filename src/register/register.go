package register

// register implementation

import (
	"errors"
	"fmt"
	//"sync"
	"time"
)

type registerMsg struct {
	pid   *Pid
	actor *Actor
	cmd   string
}

type Register struct {
	index map[int64]*Actor
	//mutex sync.Mutex //@LOCK MECHANISM will be not required on channel patter SCP
	alive bool

	sendChan chan (*Actor)
	rcvChan  chan (registerMsg)
	exitChan chan (bool)
}

//Start register instance
func Start() *Register {
	r := &Register{
		index:    make(map[int64]*Actor),
		alive:    true,
		sendChan: make(chan (*Actor), 10),
		rcvChan:  make(chan (registerMsg), 10),
		exitChan: make(chan (bool), 1),
	}
	go r.loop()

	return r
}

//Set Actor on Register
func (r *Register) Set(a *Actor) {
	msg := &registerMsg{
		pid:   a.pid,
		actor: a,
		cmd:   "set",
	}
	r.rcvChan <- *msg
}

//Get Actor from PID on register
func (r *Register) Get(p *Pid) (*Actor, error) {
	msg := &registerMsg{
		pid: p,
		cmd: "get",
	}
	r.rcvChan <- *msg
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(1 * time.Second)
		timeout <- true
	}()
	select {
	case a := <-r.sendChan:
		return a, nil
	case <-timeout:
		return nil, errors.New("timeout")
	}
}

func (r *Register) loop() {
	for r.alive {
		select {
		case msg := <-r.rcvChan:
			switch msg.cmd {
			case "set":
				err := r.set(msg.actor)
				if err != nil {
					fmt.Printf("pid %s error On Set %s \n", msg.pid, err)
				}
			case "get":
				a, err := r.get(msg.pid)
				if err != nil {
					fmt.Printf("pid %s error On get %s \n", msg.pid, err)
				} else {
					fmt.Printf("pid %s Found get \n", msg.pid)

					r.sendChan <- a
				}
			}
		case msg := <-r.exitChan:
			fmt.Printf("finish %d \n", msg)
			r.alive = false
			return
		}
	}
}

func (r *Register) set(a *Actor) error {
	//r.mutex.Lock()
	r.index[a.pid.id] = a
	//r.mutex.Unlock()
	return nil
}

func (r *Register) get(pid *Pid) (a *Actor, err error) {
	if a, ok := r.index[pid.id]; ok {
		return a, nil
	}

	return nil, errors.New("Not found")
}
