package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
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

type SafeClosable interface {
	isClosed() bool
	safeWrite(val interface{}) bool
	safeClose() bool
	safeRead() (interface{}, bool)
}

type MyChannel struct {
	C      chan interface{}
	locker sync.Mutex
	closed bool
	//once sync.Once
	//closed unsafe.Pointer //atomic 不保证可见性
}

func (m *MyChannel) safeRead() (interface{}, bool) {
	ret, ok := <-m.C
	return ret, ok
}

func (m *MyChannel) safeWrite(val interface{}) bool {
	if m.closed {
		return false
	}
	m.locker.Lock()
	defer m.locker.Unlock()
	if m.closed {
		return false
	}
	m.C <- val
	return true
}

func (m *MyChannel) isClosed() bool {
	return m.closed
	//return *(*bool)(atomic.LoadPointer(&m.closed))
}

func (m *MyChannel) safeClose() bool {
	/*
		m.once.Do(func() {
			m.closed = true
			//atomic.StorePointer(&m.closed, unsafe.Pointer(&flag))
			close(m.C)
		})
	*/
	if m.closed {
		return true
	}
	m.locker.Lock()
	defer m.locker.Unlock()
	if m.closed {
		return true
	}
	close(m.C)
	m.closed = true
	return m.closed
}

type MyCloseChannel struct {
	C      chan interface{}
	close  chan struct{} //信号
	locker sync.Mutex
	once   sync.Once
	closed bool
}

func (m *MyCloseChannel) isClosed() bool {
	return m.closed
}

func (m *MyCloseChannel) safeWrite(val interface{}) bool {
	if m.closed {
		return false
	}
	select {
	case <-m.close:
		m.once.Do(func() {
			close(m.C)
		})
		return false
	case m.C <- val:
		return true
	}
}

func (m *MyCloseChannel) safeClose() bool {
	if m.closed {
		return true
	}
	close(m.close)
	m.closed = true
	return true
}

func (m *MyCloseChannel) safeRead() (interface{}, bool) {
	ret, ok := <-m.C
	return ret, ok
}

func TestSafeClose(t *testing.T) {
	SafeClose(&MyChannel{C: make(chan interface{})})
	SafeClose(&MyCloseChannel{C: make(chan interface{}), close: make(chan struct{})})
}

func SafeClose(safe SafeClosable) {
	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup

	loop := 10
	wg.Add(loop)
	// consumer
	for i := 0; i < loop; i++ {
		go func(h int) {
			defer wg.Done()
			for {
				ret, ok := safe.safeRead()
				if !ok {
					break
				}
				fmt.Println("consume", ret, h)
			}
		}(i)
	}
	wg.Add(loop)
	// producer
	for i := 0; i < loop; i++ {
		go func(h int) {
			defer wg.Done()
			for {
				value := rand.Intn(1000)
				if !safe.safeWrite(value) {
					break
				}
				if value > 980 {
					safe.safeClose()
					break
				}
			}
		}(i)
	}

	wg.Wait()

}
