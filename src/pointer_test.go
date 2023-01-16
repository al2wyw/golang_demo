package main

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestPoniter(t *testing.T) {
	test := Float64bits(124.34)
	fmt.Println(test)
}

//test pointer uintptr
func Float64bits(f float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&f))
}
