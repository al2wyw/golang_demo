package main

import (
	"fmt"
	"testing"
	"time"
)

func TestChan(t *testing.T) {
	testBufferedChan()

	//testPip()

	testSelect()
}

func testPip() {
	pip := make(chan int, 100)
	pip <- 100

	go func() {
		count := 10
		for i := 0; i < count; i++ {
			val, ok := <-pip
			time.Sleep(1000 * time.Millisecond)
			fmt.Println("wake up start to work", val, ok)
		}
	}()
	close(pip)
	time.Sleep(10 * time.Second)
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

func testSelect() {
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
