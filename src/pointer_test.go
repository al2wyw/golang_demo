package main

import (
	"testing"
	"unsafe"
)

func TestPoniter(t *testing.T) {

}

//test pointer uintptr
func Float64bits(f float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&f))
}
