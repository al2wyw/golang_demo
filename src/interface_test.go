package main

import (
	"fmt"
	"testing"
)

type Duck interface {
	Quack()
}

type Cat struct{}

func (c *Cat) Quack() {
	fmt.Println("meow")
}

type empty int

func (em *empty) doNothing(v interface{}) {
	if v == nil {
		return
	}
	fmt.Printf("doNothing %p 成功\n", &v)
}

type Func interface {
	doFun(v interface{})
}

type FuncHandle func(interface{})

func (fun FuncHandle) doFunc(v interface{}) {
	fun(v)
}

func FuncCaller(v interface{}, fun func(interface{})) {
	h := FuncHandle(fun)
	h.doFunc(v)
}

func TestFunCaller(t *testing.T) {
	var em empty = 1
	var val int = 1
	em++

	em.doNothing(val)
	FuncCaller(val, em.doNothing)
	//fmt.Println(em == val) //em 的值是int，但是类型不是int

	var duck Duck = &Cat{}
	duck.Quack()
}
