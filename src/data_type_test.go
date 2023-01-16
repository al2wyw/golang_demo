package main

import (
	"fmt"
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

	in = data
	fmt.Println(in, in == nil) // <nil> false
}
