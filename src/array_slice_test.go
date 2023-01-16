package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func TestArray(f *testing.T) {
	array := []int{1, 2, 3, 4}
	test1 := array[0]
	testPtr(&array[0])
	testArray(array)

	fmt.Println(test1)

	for _, t := range array {
		fmt.Println("array elements : ", t)
	}

	for i := 0; i < 10; i++ {
		func() { fmt.Println("loop", i) }() // 输出全部为10
	}

	testReSlice()
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

func testReSlice() {
	fix := []int{1, 2, 3}
	slis := fix[1:3]
	fmt.Println("slis len cap", len(slis), cap(slis))
	slis[0] = 10
	//slis[3] = 100 index 错误
	fmt.Printf("slis addr %p\n", slis)
	slis = append(slis, 100) //append完生成新的slice
	fmt.Printf("slis addr %p\n", slis)
	slis[0] = 101 //不再改变原来的fix
	fmt.Println("fix", fix)
	fmt.Println("slis", slis)

	slis2 := make([]int, 3, 5)
	//slis2[3] = 10 index 错误
	fmt.Printf("slis2 addr %p\n", slis2)
	slis2 = append(slis2, 10)
	fmt.Printf("slis2 addr %p\n", slis2)
	fmt.Println("slis2", slis2)

	path := []byte("AAAA/BBBBBBBBB")
	sepIndex := bytes.IndexByte(path, '/') // 4
	println(sepIndex)

	dir1 := path[:sepIndex]
	dir2 := path[sepIndex+1:]
	println("dir1: ", string(dir1)) // AAAA
	println("dir2: ", string(dir2)) // BBBBBBBBB

	dir1 = append(dir1, "suffix"...)        //神奇
	println("current path: ", string(path)) //AAAAsuffixBBBB
}
