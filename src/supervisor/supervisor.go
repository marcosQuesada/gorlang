package supervisor

//supervisor
import (
	"actor"
	"fmt"
)

type supervisor struct {
	childs map[actor.Pid]chan (bool)
}

func Start() *supervisor {
	return &supervisor{
		childs: make(map[actor.Pid]chan (bool)),
	}
}

func (s *supervisor) Watch(a *actor.Actor) chan bool {
	fmt.Println("Monitoring exit signal from ", a.GetPid())
	s.childs[a.GetPid()] = a.ExitChan

	done := make(chan bool)
	go func(done chan bool) {
		select {
		case <-a.ExitChan:
			fmt.Println("Received exit signal from ", a.GetPid())
			delete(s.childs, a.GetPid())
			done <- true
		}
		defer fmt.Println("closing loop")
	}(done)

	return done
}

func (s *supervisor) IsAlive(p actor.Pid) bool {
	_, ok := s.childs[p]
	return ok
}
