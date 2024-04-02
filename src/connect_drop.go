package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

// ss -lnt
// State       Recv-Q Send-Q                                     Local Address:Port                                                    Peer Address:Port
// LISTEN      129    128                                                   :::8083                                                              :::*

// netstat -s | grep -i listen
// 954 times the listen queue of a socket overflowed
// 954 SYNs to LISTEN sockets dropped

// net.core.somaxconn = 128
// net.core.netdev_max_backlog = 1000
// net.ipv4.tcp_max_syn_backlog = 4096

var _count = 0

func MyConnect(ctx context.Context, c net.Conn) context.Context {
	fmt.Println(time.Now(), "conn remote address", c.RemoteAddr().String())
	if _count < 1 {
		time.Sleep(1 * time.Hour)
		_count += 1
	}
	return ctx
}

type MyHandler struct {
}

func (MyHandler) ServeHTTP(write http.ResponseWriter, request *http.Request) {
	fmt.Println(time.Now(), "serve http", request)
	_, err := write.Write([]byte("ok done"))
	if err != nil {
		fmt.Println(time.Now(), "serve http err", err)
	}
}

func Run() {
	server := &http.Server{Addr: "0.0.0.0:8083", Handler: &MyHandler{}, ConnContext: MyConnect}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("serve error", err)
	}
}
