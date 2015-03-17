package actor

import (
	"fmt"
	"testing"
)

type TestActor struct {
	Actor
}

func TestStartLink(t *testing.T) {
	sndCh := make(chan (*Message), 10)
	rcvCh := make(chan (*Message), 10)
	exitCh := make(chan (bool), 1)

	a := &Actor{
		pid:           getPid(),
		alive:         true,
		sendChan:      sndCh,
		receiveChan:   rcvCh,
		terminateChan: exitCh,
	}

	state := &state{
		index: make(map[string]string),
	}

	pid := a.StartLink(state)
	if pid.id != 1 {
		t.Fail()
	}

	m := &Message{
		destination: pid,
		mode:        "synch",
	}

	a.receiveChan <- m
	result := a.sendChan
	fmt.Printf("Result is %s \n", result)
}

type state struct {
	index map[string]string
}

func (s *state) get(key string) string {
	return s.index[key]
}

func (s *state) set(key, value string) {
	s.index[key] = value
}
