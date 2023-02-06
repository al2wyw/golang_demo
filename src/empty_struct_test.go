package main

import (
	"fmt"
	"testing"
	"unsafe"
)

type EmptyUser struct {
	name struct{}
	age  struct{}
}

type Empty struct{}

func TestEmpty(t *testing.T) {
	var emp Empty
	fmt.Println(unsafe.Sizeof(emp))

	var u EmptyUser
	fmt.Println(unsafe.Sizeof(u))

	var e = Empty{}
	fmt.Printf("empty address %p\n", &e)
	fmt.Printf("empty address %p\n", &struct{}{})
}
