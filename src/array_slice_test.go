package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func TestArray(f *testing.T) {

	balance := [...]float32{1000.0, 2.0, 3.4, 7.0, 50.0} //balance是数组
	fmt.Println("balance", balance)

	array := []int{1, 2, 3, 4}
	test1 := array[0]
	testPtr(&array[0])
	testArray(array)
	testAppend(array)
	testAppendV1(&array)

	fmt.Println("test1", test1)

	for _, t := range array {
		fmt.Println("array elements : ", t)
	}

	//test append
	slis2 := make([]int, 3, 5)
	//slis2[3] = 10 index 错误
	fmt.Printf("slis2 addr %p\n", slis2)
	slis2 = append(slis2, 10)            //enough space
	fmt.Printf("slis2 addr %p\n", slis2) //the same as above
	fmt.Println("slis2", slis2)
	slis2 = append(slis2, 11, 12)
	fmt.Println("slis2 len cap", len(slis2), cap(slis2)) //cap * 2

	var slis3 []int
	// slis3[0] = 10 //nil slice 不可以直接通过下标访问
	slis3 = append(slis3, 1) //nil slice 可以直接append

	testReSlice()
}

func testAppend(p []int) {
	//original slice not changed
	p = append(p, 10)
	fmt.Println("test append", p)
}

func testAppendV1(p *[]int) {
	//original slice changed
	*p = append(*p, 10)
	fmt.Println("test append", *p)
}

func testArray(p []int) {
	p[0] = 100
	fmt.Println("test array", p[0])
	var array []student
	array = append(array, student{Name: "test"})

	data, err := json.Marshal(array)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(string(data))
}

func testPtr(a *int) {
	*a = 10
}

/**
slice := []int{1, 2, 3, 4, 5}
newSlice := slice[0:3:3]
第一个值：截取的起始索引
第二个值：截取的终止索引（不包括该值）
第三个值：用来计算切片的容量，可以省略，默认和长度一样
容量 = 第三个值 - 第一个值
长度 = 第二个值 - 第一个值
*/
func testReSlice() {
	fix := []int{1, 2, 3}

	slisTest := fix[0:1]
	fmt.Println("slisTest len cap", len(slisTest), cap(slisTest))

	slis := fix[1:3]
	fmt.Printf("slis address %p, %p \n", fix, slis) //并没有重新分配内存
	fmt.Println("slis len cap", len(slis), cap(slis))
	slis[0] = 10
	//slis[3] = 100 index 错误
	fmt.Println("fix", fix) //fix[1] 也变成了10
	fmt.Printf("slis addr %p\n", slis)
	slis = append(slis, 100) //append完生成新的slice
	fmt.Printf("slis addr %p\n", slis)
	slis[0] = 101 //不再改变原来的fix
	fmt.Println("fix", fix)
	fmt.Println("slis", slis)

	path := []byte("AAAA/BBBBBBBBB")
	sepIndex := bytes.IndexByte(path, '/') // 4
	println(sepIndex)

	dir1 := path[:sepIndex]         // len 4 cap 14
	dir2 := path[sepIndex+1:]       // len 9 cap 9
	println("dir1: ", string(dir1)) // AAAA
	println("dir2: ", string(dir2)) // BBBBBBBBB

	dir1 = append(dir1, "suffix"...)        //神奇
	println("current path: ", string(path)) //AAAAsuffixBBBB
}
