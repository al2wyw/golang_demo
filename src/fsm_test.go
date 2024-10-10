package main

import (
	"context"
	"fmt"
	"github.com/looplab/fsm"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFsm(t *testing.T) {
	var afterFinishCalled bool
	machine := fsm.NewFSM(
		"start",
		fsm.Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
			{Name: "finish", Src: []string{"end"}, Dst: "finished"},
			{Name: "reset", Src: []string{"end", "finished"}, Dst: "start"},
		},
		fsm.Callbacks{
			"enter_end": func(ctx context.Context, e *fsm.Event) {
				if err := e.FSM.Event(ctx, "finish"); err != nil {
					fmt.Println(err)
				}
			},
			"after_finish": func(ctx context.Context, e *fsm.Event) {
				afterFinishCalled = true
				if e.Src != "end" {
					panic(fmt.Sprintf("source should have been 'end' but was '%s'", e.Src))
				}
				if err := e.FSM.Event(ctx, "reset"); err != nil {
					fmt.Println(err)
				}
			},
			"enter_state": func(ctx context.Context, e *fsm.Event) {
				fmt.Printf("transit from %s -> %s\n", e.Src, e.Dst)
			},
		},
	)

	println(fsm.Visualize(machine))

	if err := machine.Event(context.Background(), "run"); err != nil {
		panic(fmt.Sprintf("Error encountered when triggering the run event: %v", err))
	}

	if !afterFinishCalled {
		panic(fmt.Sprintf("After finish callback should have run, current state: '%s'", machine.Current()))
	}

	currentState := machine.Current()
	if currentState != "start" {
		panic(fmt.Sprintf("expected state to be 'start', was '%s'", currentState))
	}

	fmt.Println("Successfully ran state machine.")
}

func TestCaller(t *testing.T) {
	Where()
	Where()
}

func Where() {
	_, file, no, _ := runtime.Caller(1)
	fmt.Printf("file: %s, line: %d\n", filepath.Base(file), no)
}
