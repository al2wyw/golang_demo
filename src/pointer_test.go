package main

import (
	"fmt"
	"strconv"
	"testing"
	"unsafe"
)

//任何类型的指针都可以被转化为 unsafe.Pointer
//unsafe.Pointer 可以被转化为任何类型的指针
//uintptr 可以被转化为 unsafe.Pointer
//unsafe.Pointer 可以被转化为 uintptr

type Human struct {
	i int64
	j int8
	k int32
	n string
}

func TestPointer(t *testing.T) {

	i := 10
	p := &i
	fp := (*float32)(unsafe.Pointer(p))
	fmt.Println("float32 addr value", fp, *fp)
	*fp = *fp * 10
	fmt.Println("float32 addr value", fp, *fp)
	fmt.Println("int addr value", &i, i)

	test := Float64bits(125.34)
	fmt.Println("test", test)

	//小 转 大 出现异常
	var small int8 = 34
	var large = *(*int32)(unsafe.Pointer(&small))
	fmt.Println("my large value", large)

	var big uint16 = 773 //
	var little = *(*uint8)(unsafe.Pointer(&big))
	fmt.Println("my big value", big>>8, big&0xf)
	fmt.Println("my little value", little)                                                       //5
	fmt.Println("my little value", *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&big)) + 1))) //3

	arr := []int{1, 2, 3, 4}
	ptr := uintptr(unsafe.Pointer(&arr[1]))
	qtr := uintptr(unsafe.Pointer(&arr[0])) + unsafe.Sizeof(arr[0])
	fmt.Println("p q addr value", ptr, strconv.FormatInt(int64(ptr), 16), qtr, strconv.FormatInt(int64(qtr), 16))

	barr := []byte{97, 66, 64, 65}
	str := BytesToStr(barr)
	fmt.Println(str)

	human := Human{
		9223372036854775807, 34, 34, "test",
	}
	//必须考虑内存对齐，Sizeof返回的是多少byte
	fmt.Println("size of", unsafe.Sizeof(&human), unsafe.Sizeof(human), unsafe.Sizeof(human.i), unsafe.Sizeof(human.j), unsafe.Sizeof(human.k), unsafe.Sizeof(human.n))

	//只有*AnyType的变量才能转为unsafe.Pointer，即使是引用类型的变量也不能直接转为unsafe.Pointer
	slice := []int{1, 2, 3}
	//_ = unsafe.Pointer(slice)
	fmt.Printf("addrs %p, %p, %p \n", slice, &slice, unsafe.Pointer(&slice))

	var interf interface{} = &i
	fmt.Printf("addrs %p, %p, %p, %p \n", &i, interf, &interf, unsafe.Pointer(&interf))
}

func Float64bits(f float64) uint64 {
	fmt.Printf("float64 address %p , %T\n", &f, f)
	//转化的目标类型（uint64) 的 size 一定不能比原类型 (float64) 还大（二者size都是8个字节）
	//前后两种类型有等价的 memory layout
	return *(*uint64)(unsafe.Pointer(&f))
}

func BytesToStr(arr []byte) string {
	//避免内存拷贝
	return *(*string)(unsafe.Pointer(&arr))
}
