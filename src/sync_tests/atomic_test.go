package sync_tests

import (
	"fmt"
	"sync/atomic"
	"testing"
	"unsafe"
)

type AtomicData struct {
	name string
	age  int
}

type AtomicIf interface {
	show() string
}

func (a AtomicData) show() string {
	return a.name
}

func TestAtomic(t *testing.T) {
	var target AtomicData
	ptr := unsafe.Pointer(&target)
	val := AtomicData{name: "test", age: 100}
	atomic.StorePointer(&ptr, unsafe.Pointer(&val)) //只是改变了ptr的值

	fmt.Println(*(*AtomicData)(ptr), target) //target 没变

	var targetPtr *AtomicData
	unPtr := (*unsafe.Pointer)(unsafe.Pointer(&targetPtr)) //这个不好理解啊
	atomic.StorePointer(unPtr, unsafe.Pointer(&val))
	fmt.Println(*targetPtr)

	var targetIntf interface{}
	unIntfPtr := (*unsafe.Pointer)(unsafe.Pointer(&targetIntf))
	atomic.StorePointer(unIntfPtr, unsafe.Pointer(&val))
	fmt.Println(targetIntf) // <invalid reflect.Value>

	/*
		//panic: invalid memory address or nil pointer dereference
		var targetI AtomicIf
		unIntf := (*unsafe.Pointer)(unsafe.Pointer(&targetI))
		atomic.StorePointer(unIntf, unsafe.Pointer(&val))
		fmt.Println(targetI)
	*/
}
