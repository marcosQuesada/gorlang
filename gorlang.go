package main

import (
	"actor"
	"fmt"
	"os"
	"os/signal"
	"register"
	"runtime"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	r := register.Start()
	fmt.Println("Started Register on pid", r.GetPid())
	a := actor.NewActor()
	err := r.Set(a)
	if err != nil {
		fmt.Println("Error on Register set ", err)
	}
	ra, err := r.Get(actor.Pid(2))
	if err != nil {
		fmt.Println("Error on register get ", err, ra.GetPid())
	}

	// Block until a signal is received.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-c
}
