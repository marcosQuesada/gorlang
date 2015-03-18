package actor

import (
	"errors"
	"os"
	"testing"
)

var a *Actor
var st *state

func TestMain(m *testing.M) {
	a = NewActor()
	st = &state{
		index: make(map[string]string),
	}

	os.Exit(m.Run())
}

func TestSpawn(t *testing.T) {
	pid := a.Spawn(st)
	if pid != 1 {
		t.Fail()
	}

	resultGet, err := a.Send(pid, "Set", "key_foo", "13456789765")
	if err != nil {
		t.Error("Error in Send", err)
	}
	if resultGet != nil {
		t.Error("Result don't match, %d", resultGet)
	}
	resultGet, err = a.Send(pid, "Get", "key_foo")
	if err != nil {
		t.Error("Error in Get", err)
	}
	if resultGet != "13456789765" {
		t.Error("Result don't match, %d", resultGet)
	}
	a.Terminate()
}

func TestErrorOnNotExistentMethod(t *testing.T) {
	//inexistent command
	ab := NewActor()
	pid := ab.Spawn(st)
	_, err := ab.Send(pid, "foo", "key_foo")
	if err == nil {
		t.Error("Expected not found")
	}
	ab.Terminate()
}

func TestLoadMethods(t *testing.T) {
	ac := NewActor()
	ac.Spawn(st)
	if len(ac.methodIndex) != 2 {
		t.Fail()
	}
	ac.Terminate()

}

func TestMultiPidActors(t *testing.T) {
	b := NewActor()
	stb := &state{
		index: make(map[string]string),
	}

	b.Spawn(stb)
	if b.pid != 4 {
		t.Error("Bad Pid Number ", b.pid)
	}

	// check that is a different actor than the firt one
	resultGet, err := b.Send(b.pid, "Get", "key_foo")
	if resultGet != "" {
		t.Error("Get on unexistent key ", resultGet, b.pid)
	}

	if err == nil {
		t.Error("Expected Not found error not raised ")
	}
	b.Terminate()
}

// Internal State
type state struct {
	index map[string]string
}

func (s *state) Get(key string) (string, error) {
	if _, ok := s.index[key]; ok {
		return s.index[key], nil
	}
	return "", errors.New("Not found")

}

func (s *state) Set(key, value string) {
	s.index[key] = value
}
