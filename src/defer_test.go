package main

import (
	"fmt"
	"runtime/debug"
	"sync"
	"testing"
	"time"
)

func TestDefer(t *testing.T) {
	testDefer()
	fmt.Println("defer return", testDeferV2())

	testTypeAssert(1)
	testTypeAssert("test")

	testRoutineDefer()
}

func testRoutineDefer() {
	var deferCount = 0
	wg := new(sync.WaitGroup)
	max := 10
	wg.Add(max)
	for loop := 0; loop < max; loop++ {
		go func(i int) {
			defer func() {
				wg.Done()
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

func testDefer() {
	var i = 1

	defer fmt.Println("result: ", func() int { fmt.Println("call: ", i); return i * 2 }()) //2
	i++
	fmt.Println("main result: ", i) //2
}

func testDeferV2() int {
	var i = 1

	defer func() int { i *= 2; fmt.Println("result: ", i); return i }() //2
	i++
	fmt.Println("main result: ", i) //2
	return i
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
