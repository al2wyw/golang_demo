package main

import (
	"fmt"
	"reflect"
	"testing"
)

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
}
