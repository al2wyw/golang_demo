package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"testing"
	"time"
)

func TestDefer(t *testing.T) {
	testDefer(t)
	ret := testDeferV2(t)
	assert.Equal(t, 2, ret)
	fmt.Println("defer return", ret)

	testTypeAssert(1)
	testTypeAssert("test")
}

// TestDeferExit 优雅退出
func TestDeferExit(t *testing.T) {
	defer func() {
		fmt.Println("exit")
	}()

	sign := make(chan os.Signal)
	signal.Notify(sign, syscall.SIGTERM) //如果没有signal， kill会让程序直接退出，defer不会运行

	timer := time.NewTimer(3 * time.Second)
	for {
		select {
		case s := <-sign:
			fmt.Println("signal", s)
			return
		case <-timer.C:
			fmt.Println("alive")
		}
		timer.Reset(3 * time.Second)
	}

}

func TestRoutineDefer(t *testing.T) {
	var deferCount = 0
	wg := new(sync.WaitGroup)
	max := 10
	wg.Add(max)
	for loop := 0; loop < max; loop++ {
		go func(i int) {
			defer func() {
				wg.Done()
				fmt.Printf("defer is called\n")
				if err := recover(); err != nil {
					fmt.Printf("panic error %+v\n%s", err, debug.Stack())
				}
			}()
			for {
				time.Sleep(100 * time.Millisecond)
				fmt.Println("wake up to work", i)
				deferCount++
				if deferCount%10 == 0 {
					panic(fmt.Errorf("this is 10, %d", i)) //如果没有recover，一旦panic整个程序都结束，recover只是防止整个程序结束
				}
			}
		}(loop)
	}

	wg.Wait()
}

func TestFatalErrorCannotRecover(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic error", err)
		}
	}()
	fatalMap()

	fmt.Println("done") //never executed
}

// fatalMap fatal error: concurrent map read and map write
func fatalMap() {
	loop := 100000
	m := map[string]int{}
	go func() {
		for i := 0; i < loop; i++ {
			m["煎鱼1"] = 1
		}
	}()
	for i := 0; i < loop; i++ {
		_ = m["煎鱼2"]
	}
}

// 参考goroutine + 闭包的笔记
func testDefer(t *testing.T) {
	var i = 1

	defer fmt.Println("result: ", func() int {
		fmt.Println("call: ", i)
		assert.Equal(t, 2, i*2)
		return i * 2
	}()) //2
	i++
	fmt.Println("main result: ", i) //2
	assert.Equal(t, 2, i)
}

func testDeferV2(t *testing.T) int {
	var i = 1

	defer func() int {
		i *= 2
		fmt.Println("result: ", i)
		assert.Equal(t, 4, i)
		return i
	}() //4
	i++
	fmt.Println("main result: ", i) //2
	assert.Equal(t, 2, i)
	return i //2
}

func testTypeAssert(object interface{}) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic error", err)
		}
	}()

	t, ok := object.(string)
	fmt.Println("this is type assert", ok)

	if t == "test" {
		fmt.Println("this is a test")
	}
}
