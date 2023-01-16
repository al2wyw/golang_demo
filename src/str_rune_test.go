package main

import (
	"fmt"
	"testing"
)

func TestStr(t *testing.T) {
	testStrTsf()
}

func testStrTsf() {
	s := "我是ABC"
	//s[0] is uint8
	fmt.Println("this is a utf8 str:", s[0], s[6])

	for _, v := range s {
		//v is int32
		fmt.Println("this is a range:", v)
		fmt.Println("this is a utf8 range:", string(v))
	}

	arr := []rune(s)
	arr[0] = rune('他') //Redundant type conversion
	fmt.Println("this is a utf8 str:", string(arr))

}
