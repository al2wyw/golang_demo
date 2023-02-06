package main

import (
	"fmt"
	"testing"
)

type empty int

func (em *empty) doNothing(v interface{}) {
	if v == nil {
		return
	}
	*em += 1
	fmt.Printf("doNothing %p  %d 成功\n", &v, *em)
}

func doNothing(v interface{}) {
	if v == nil {
		return
	}
	fmt.Printf("doSomething %p 成功\n", &v)
}

/////////////接口型函数的demo////////////
type Func interface {
	doFun(v interface{})
}

func testFuncDriver(fun Func) {
	fmt.Println("func interface method called", fun)
}

//以下是对以上接口函数调用的优化

//在其他语言里面，有些函数可以直接作为参数传递，有些是以函数指针进行传递，但是都没有办法像go这样可以给函数类型“增加”新方法
type FuncHandle func(interface{})

func (fun FuncHandle) doFun(v interface{}) {
	fun(v)
}

func FuncCaller(v interface{}, fun func(interface{})) {
	h := FuncHandle(fun) //类型装换
	h.doFun(v)
	testFuncDriver(h)
}

/////////////接口型函数的demo////////////

func TestFunCaller(t *testing.T) {
	var em empty = 1
	var val int = 1
	em++

	em.doNothing(val)
	FuncCaller(val, em.doNothing)
	var emn empty = 10
	FuncCaller(val, emn.doNothing)
	//fmt.Println(em == val) //em 的值是int，但是类型不是int

	FuncCaller(val, doNothing)
}
