package main

import (
	"fmt"
	"reflect"
	"testing"
)

type EmptyInt int

func (em *EmptyInt) DoNothing(v interface{}) {
	if v == nil {
		return
	}
	*em += 1
	fmt.Printf("DoNothing %p  %d 成功\n", &v, *em)
}

func (em EmptyInt) ShowMySelf(v interface{}) {
	fmt.Printf("show %d 成功\n", em)
}

func DoNothing(v interface{}) {
	if v == nil {
		return
	}
	fmt.Printf("doSomething %p 成功\n", &v)
}

/////////////接口型函数的demo////////////
type Func interface {
	DoFun(v interface{})
}

func FuncDriver(fun Func, v interface{}) {
	fun.DoFun(v)
	fmt.Println("func interface method called")
}

//以下是对以上接口函数调用的优化: 把任意的func(interface{})函数转化为Func类型传入FuncDriver

//在其他语言里面，有些函数可以直接作为参数传递，有些是以函数指针进行传递，但是都没有办法像go这样可以给函数类型“增加”新方法
type FuncHandle func(interface{})

func (fun FuncHandle) DoFun(v interface{}) {
	fun(v)
}

func FuncCaller(v interface{}, fun func(interface{})) {
	//h := FuncHandle(fun) //类型装换
	var h FuncHandle = fun // var em EmptyInt = val //error
	FuncDriver(h, v)
}

/////////////接口型函数的demo////////////

func TestFunCaller(t *testing.T) {
	var em EmptyInt = 1
	var val int = 1
	em++

	em.DoNothing(val)
	FuncCaller(val, em.DoNothing)
	var emn EmptyInt = 10
	FuncCaller(val, emn.DoNothing)
	//fmt.Println(em == val) //em 的值是int，但是类型不是int

	FuncCaller(val, DoNothing)

	klass := reflect.TypeOf(&em)
	for i := 0; i < klass.NumMethod(); i++ {
		m := klass.Method(i)
		fmt.Printf("%s: %v\n", m.Name, m.Type)
	}

}
