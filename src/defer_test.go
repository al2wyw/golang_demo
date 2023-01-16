package main

import (
	"fmt"
	"testing"
)

func TestDefer(t *testing.T) {
	testDefer()

	testTypeAssert(1)
	testTypeAssert("test")
}

func testDefer() {
	var i = 1
	//预先设置好引用的变量
	defer fmt.Println("result: ", func() int { return i * 2 }()) //2
	i++
	fmt.Println("main result: ", i) //2
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
