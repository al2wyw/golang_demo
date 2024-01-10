package main

import "testing"
import "net/http"
import "time"
import _ "net/http/pprof"

// profile 和 trace: profile是cpu时间，trace是io wait时间
// go tool pprof http://127.0.0.1:8080/debug/pprof/profile?seconds=10
// curl http://127.0.0.1:8080/debug/pprof/trace?seconds=10 > trace.out
// go tool pprof -http=127.0.0.1:9888 profile.out
// go tool trace -http=127.0.0.1:9998 trace.out
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
