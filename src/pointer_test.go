package main

import (
	"fmt"
	"testing"
	"unsafe"
)

//任何类型的指针都可以被转化为 unsafe.Pointer
//unsafe.Pointer 可以被转化为任何类型的指针
//uintptr 可以被转化为 unsafe.Pointer
//unsafe.Pointer 可以被转化为 uintptr

func TestPoniter(t *testing.T) {

	i := 10
	p := &i
	fp := (*float32)(unsafe.Pointer(p))
	fmt.Println("float32 addr value", fp, *fp)
	*fp = *fp * 10
	fmt.Println("float32 addr value", fp, *fp)
	fmt.Println("int addr value", &i, i)

	test := Float64bits(125.34)
	fmt.Println("test", test)

	arr := []int{1, 2, 3, 4}
	ptr := unsafe.Pointer(&arr[1])
	qtr := uintptr(unsafe.Pointer(&arr[0])) + unsafe.Sizeof(arr[0])
	fmt.Println("p q addr value", ptr, qtr)

}

//test pointer uintptr
func Float64bits(f float64) uint64 {
	fmt.Printf("float64 address %p , %T\n", &f, f)
	//转化的目标类型（uint64) 的 size 一定不能比原类型 (float64) 还大（二者size都是8个字节）
	//前后两种类型有等价的 memory layout
	return *(*uint64)(unsafe.Pointer(&f))
}
