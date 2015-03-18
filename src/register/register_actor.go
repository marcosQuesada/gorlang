package register

import (
	"actor"
)

// register Actor implementation
type RegisterActor struct {
	*actor.Actor
}

//Start register instance
func Start() *RegisterActor {
	r := &RegisterActor{actor.NewActor()}

	st := &Register{
		index: make(map[actor.Pid]*actor.Actor),
	}
	r.Spawn(st)
	return r
}

//Set Actor on Register
func (r *RegisterActor) Set(a *actor.Actor) (err error) {
	_, err = r.Send(r.GetPid(), "Set", a.GetPid(), a)
	return
}

//Get Actor from PID on register
func (r *RegisterActor) Get(p actor.Pid) (a *actor.Actor, err error) {
	res, err := r.Send(r.GetPid(), "Get", p)
	if err != nil {
		return nil, err
	}
	return res.(*actor.Actor), nil
}
