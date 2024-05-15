package main

import (
	"fmt"
	"testing"
)

type canPrint interface {
	print()
}

type genericBase struct {
	name  string
	value int32
}

func (t *genericBase) print() {
	fmt.Println("genericBase print", t.name, t.value)
}

type genericExtend struct {
	genericBase
	salary float64
}

func (t *genericExtend) print() {
	fmt.Println("genericExtend print", t.name, t.value, t.salary)
}

// 不能使用 struct 作为泛型约束, 只能使用 interface 作为泛型约束
func testPrint[T genericBase](t T) {
	//fmt.Println("testPrint", t.name, t.value)
	//t.print()
	fmt.Println("testPrint", t)
}

func testGeneric[T canPrint](t T) {
	t.print()
}

type genericContainer[T canPrint] struct {
	data []T
}

func (t *genericContainer[T]) print() {
	fmt.Printf("genericContainer print: ")
	for _, i := range t.data {
		i.print()
	}
}

type Container[T canPrint] interface {
	Len() int
	Add(T)
	Remove() (T, error)
}

func (t *genericContainer[T]) Add(t2 T) {
	t.data = append(t.data, t2)
}

func (t *genericContainer[T]) Remove() (T, error) {
	length := t.Len()
	if length > 0 {
		ret := t.data[length-1]
		t.data = t.data[0 : length-1]
		return ret, nil
	}
	var a T
	return a, fmt.Errorf("size error")
	//不能return nil, fmt.Errorf("") !!!
}

func (t *genericContainer[T]) Len() int {
	return len(t.data)
}

func TestGeneric(tt *testing.T) {
	base := &genericBase{name: "test", value: 34}
	testGeneric(base)

	arr := []canPrint{base}
	container := &genericContainer[canPrint]{arr}
	testGeneric(container)

	container.Add(&genericBase{name: "peter", value: 33})
	container.Add(&genericBase{name: "anne", value: 30})
	testGeneric(container)
	if _, err := container.Remove(); err != nil {
		fmt.Println(err)
	}
	testGeneric(container)

	if _, err := container.Remove(); err != nil {
		fmt.Println(err)
	}
	if _, err := container.Remove(); err != nil {
		fmt.Println(err)
	}
	if _, err := container.Remove(); err != nil {
		fmt.Println(err)
	}

	testPrint(*base)
}
