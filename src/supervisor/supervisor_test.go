package supervisor

import (
	"actor"
	"testing"
	//"time"
)

func TestStart(t *testing.T) {
	s := Start()

	b := actor.NewActor()

	b.Spawn(&stateTest{})

	done := s.Watch(b)

	if !s.IsAlive(b.GetPid()) {
		t.Error("Actor Not supervised!")
	}

	b.Terminate()
	<-done
	if s.IsAlive(b.GetPid()) {
		t.Error("Actor Must be died, but is supervised!")
	}
}

type stateTest struct {
	index int
}
