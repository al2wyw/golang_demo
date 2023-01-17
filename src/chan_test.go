package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestChan(t *testing.T) {
	testWaitG()
	testBufferedChan()

	//testPip()

	//testSelect()

	testMutex()
}

var comtact_count int64 = 0
var lock = new(sync.Mutex)

func testMutex() {
	wg := new(sync.WaitGroup)
	loop := 100
	wg.Add(loop)

	for i := 0; i < loop; i++ {
		go func() {
			//lock.Lock()
			defer wg.Done() //; lock.Unlock()
			comtact_count++
			fmt.Println("comtact_count", comtact_count) //没有锁顺序错乱但是结果是对的，没有出现互相覆盖， wg有同步能力???
		}()
	}

	wg.Wait()
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

func testWaitG() {
	ch := make(chan string)
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		defer wg.Done()
		//ch不close, range就不会结束
		for m := range ch {
			fmt.Println("Processed:", m)
			time.Sleep(100 * time.Millisecond) // 模拟需要长时间运行的操作
		}
	}()

	ch <- "cmd.1"
	ch <- "cmd.2" // 不会被接收处理
	close(ch)
	wg.Wait()
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
				//default 一直执行
				//default:
				//	fmt.Println("this is default")
			}
		}
	}()
	_ = <-time.After(10 * time.Second)
	timeout <- true
	fmt.Println("testSelect end")
}
