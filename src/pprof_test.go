package main

import "testing"
import "net/http"
import "time"
import _ "net/http/pprof"

var leak [][1024]byte

func TestProf(t *testing.T) {
	go func() {
		time.Sleep(time.Second)
		v := [1024]byte{}
		leak = append(leak, v)
	}()

	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		return
	}
}
