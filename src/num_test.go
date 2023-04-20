package main

import (
	"fmt"
	"testing"
)

func TestOverflow(tt *testing.T) {
	testIntOver()
	testUnsignedIntOver()
}

func testIntOver() {
	var i int8 = 0x7f
	var j int8 = 0x1
	var x = i + j
	var y = x - i
	fmt.Println(x, y, x < i, y < 0, x-i < 0)
}

func testUnsignedIntOver() {
	var i uint8 = 0xff
	var j uint8 = 0x1
	var x = i + j
	var y = x - i
	fmt.Println(x, y, x < i, y < 0, x-i < 0)
}
