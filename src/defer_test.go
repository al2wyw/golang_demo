package main

import (
	"fmt"
	"testing"
)

func TestDefer(t *testing.T) {
	testDefer()
	fmt.Println("defer return", testDeferV2())

	testTypeAssert(1)
	testTypeAssert("test")
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
