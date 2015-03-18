package register

// register State implementation

import (
	"actor"
	"errors"
)

// ACTOR STATE & MANIPULATORS
// Define Register State
type Register struct {
	index map[actor.Pid]*actor.Actor
}

// Get Command
func (s *Register) Get(pid actor.Pid) (*actor.Actor, error) {
	if r, ok := s.index[pid]; ok {
		return r, nil
	}
	return nil, errors.New("Not found")
}

// Set Command
func (s *Register) Set(pid actor.Pid, value *actor.Actor) {
	s.index[pid] = value
}
