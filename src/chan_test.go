package main

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestChan(t *testing.T) {
	testBufferedChan()

	testPip()
}

func TestContext(t *testing.T) {
	cxt, cancel := context.WithTimeout(context.Background(), 3*time.Second) //WithDeadline用法一致，只是第二个参数变成了固定时间点
	defer cancel()
	//调用cancel或者WithTimeout的duration到了，Done()都会返回chan

	go func() {
		select {
		case <-cxt.Done():
			fmt.Println("testContext cxt is canceled", cxt.Err())
		case <-time.After(2 * time.Second):
			fmt.Println("testContext time out")
		}
	}()

	time.Sleep(4 * time.Second)
	fmt.Println("testContext main cancel")
}

func testPip() {
	pip := make(chan int, 100)
	pip <- 100
	pip <- 120
	pip <- 130
	close(pip)
	go func() {
		count := 10
		for i := 0; i < count; i++ {
			val, ok := <-pip
			time.Sleep(100 * time.Millisecond)
			fmt.Println("wake up start to work", val, ok)
		}
	}()

	time.Sleep(1 * time.Second)
}

func testBufferedChan() {
	ch := make(chan int, 10)
	for i := 0; i < 10; i++ {
		ch <- i
	}

	go func() {
		for {
			i, ok := <-ch
			if ok {
				fmt.Println("this is value from chan", i)
			} else {
				fmt.Println("chan is false")
				return
			}
		}
	}()

	time.Sleep(100 * time.Millisecond)
	close(ch)
	time.Sleep(100 * time.Millisecond)
}

func TestSelect(t *testing.T) {
	ch := make(chan int)
	timeout := make(chan bool)
	go func() {
		for {
			select {
			case v := <-ch:
				fmt.Println("read from ch", v)
			case <-timeout:
				fmt.Println("read from timeout")
				//default 一直执行，导致 CPU 忙转
				//default:
				//	fmt.Println("this is default")
			}
		}
	}()
	_ = <-time.After(1 * time.Second)
	timeout <- true
	fmt.Println("testSelect end")
}
