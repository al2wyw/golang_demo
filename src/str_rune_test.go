package main

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"unicode/utf8"
	"unsafe"
)

func TestStr(t *testing.T) {

	var tr = new(string)
	*tr = "test" //地址保持不变
	fmt.Println(tr)

	str := "test123456789"
	str1 := "test" + strconv.Itoa(123456789)
	str2 := "test123456789"
	str3 := str2[0:4] //只是把地址偏移量赋给了data指针
	test := "test"

	fmt.Printf("str %s, %s, %s, %p, %p, %p \n", str, str1, str2, &str, &str1, &str2)

	strptr := (*reflect.StringHeader)(unsafe.Pointer(&str))
	str1ptr := (*reflect.StringHeader)(unsafe.Pointer(&str1))
	str2ptr := (*reflect.StringHeader)(unsafe.Pointer(&str2))
	str3ptr := (*reflect.StringHeader)(unsafe.Pointer(&str3))
	testptr := (*reflect.StringHeader)(unsafe.Pointer(&test))

	fmt.Println("a ptr:", unsafe.Pointer(strptr.Data))
	fmt.Println("b ptr:", unsafe.Pointer(str1ptr.Data))
	fmt.Println("c ptr:", unsafe.Pointer(str2ptr.Data)) // c == d
	fmt.Println("d ptr:", unsafe.Pointer(str3ptr.Data))
	fmt.Println("t ptr:", unsafe.Pointer(testptr.Data))

	if str == str1 { //先比较大小，然后比较data指针的值，再逐个byte进行内容对比
		fmt.Println("str is equal")
	}

	strByte := []byte(str) //类型转换发生了内容拷贝
	sliceptr := (*reflect.SliceHeader)(unsafe.Pointer(&strByte))
	fmt.Printf("str and bytes %s, %s, %p, %p, %p \n", str, strByte, &str, strByte, &strByte)
	fmt.Println("slice ptr:", unsafe.Pointer(sliceptr.Data))
	strByte[0] = 'd'
	str11 := string(strByte) //类型转换再次发生了内容拷贝
	str11ptr := (*reflect.StringHeader)(unsafe.Pointer(&str11))
	fmt.Println("str11 ptr:", unsafe.Pointer(str11ptr.Data))

	str = str1 //只是修改了data指针的值
	fmt.Printf("str %s, %s, %p, %p \n", str, str1, &str, &str1)
	fmt.Println("a ptr:", unsafe.Pointer(strptr.Data))
	fmt.Println("b ptr:", unsafe.Pointer(str1ptr.Data))

	testStrTsf()
}

func testStrTsf() {
	s := "我是ABC"
	//s[0] is uint8
	fmt.Println("this is a utf8 str:", s[0], s[6], len(s), utf8.RuneCountInString(s)) // 230, 65, 9, 5

	for _, v := range s { //range的特殊对待，按照rune slice来取值
		//v 全部都是 int32
		fmt.Println("this is a range:", v)
		fmt.Println("this is a utf8 range:", string(v))
	}

	arr := []rune(s)
	arr[0] = rune('他') //Redundant type conversion
	fmt.Println("this is a utf8 str:", string(arr))

}
