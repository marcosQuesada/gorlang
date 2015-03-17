package register

import (
	"testing"
)

func TestStart(t *testing.T) {
	r := Start()
	a := &Actor{
		pid: &Pid{id: int64(1231232)},
	}
	r.Set(a)

	ra, err := r.Get(a.pid)
	if err != nil {
		t.Fail()
	}

	if ra.pid != a.pid {
		t.Fail()
	}

	_, err = r.Get(&Pid{id: int64(111)})
	if err == nil {
		t.Fail()
	}
}
