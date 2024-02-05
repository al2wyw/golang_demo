package main

import (
	"fmt"
	"reflect"
	"testing"
)

type myStr struct {
	Id int
}
type myStrType myStr
type myType int

func TestDataType(t *testing.T) {
	var bigVal int64 = 2
	fmt.Println(bigVal)
	//bigVal = val

	var data *byte
	var in interface{}

	fmt.Println(data, data == nil) // <nil> true
	fmt.Println(in, in == nil)     // <nil> true

	fmt.Println("type", reflect.TypeOf(in))
	fmt.Println("value", reflect.ValueOf(in))

	in = data
	fmt.Println(in, in == nil) // <nil> false
	fmt.Println("type", reflect.TypeOf(in))
	fmt.Println("value", reflect.ValueOf(in))

	var test interface{} = nil
	fmt.Println(test, test == nil) // <nil> true

	fmt.Printf("%p\n", in) // 0x0

	var myInt myType = 1
	fmt.Println("type", reflect.TypeOf(myInt))
	fmt.Println("kind", reflect.TypeOf(myInt).Kind())
	fmt.Println("value", reflect.ValueOf(myInt))

	myStruct := myStrType{1}
	fmt.Println("type", reflect.TypeOf(myStruct))
	fmt.Println("kind", reflect.TypeOf(myStruct).Kind())
	fmt.Println("value", reflect.ValueOf(myStruct))
}
