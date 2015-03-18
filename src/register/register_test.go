package register

import (
	"actor"
	"testing"
)

func TestStart(t *testing.T) {
	r := Start()
	a := actor.NewActor()
	err := r.Set(a)
	if err != nil {
		t.Error("Error on Register set ", err)
	}
	ra, err := r.Get(actor.Pid(2))
	if err != nil {
		t.Error("Error on register get ", err, ra.GetPid())
	}

	if ra.GetPid() != a.GetPid() {
		t.Error("Register Pids don't match ", ra.GetPid(), a.GetPid())
	}

	v, err := r.Get(actor.Pid(99))
	if err == nil {
		t.Error("Unexpected error nil ", err, v)
	}

	a.Terminate()
	r.Terminate()
}
