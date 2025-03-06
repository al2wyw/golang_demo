package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type MyInf interface {
	Action(a int32) string
}

type MyData struct {
	Name   string
	Value  int32
	Amount int32
}

func (m MyData) Action(a int32) string {
	return strconv.Itoa(int(a))
}

type ByteSlice []byte

func (p *ByteSlice) Append(data []byte) {
	slice := *p
	// Body as above, without the return.
	*p = slice
}

func main() {

	test := make(ByteSlice, 10)
	test.Append(make([]byte, 10))

	test = append(test, byte(43))

	ret, _ := compute(10, 10, "+")

	str := readStr()
	var data = MyData{str, 10, 10}
	my := data
	var inf MyInf = my

	fmt.Println(ret, inf)
}

func readStr() string {
	path, _ := os.Getwd()
	bys, _ := ioutil.ReadFile(path + "/test.txt")
	return string(bys)
}

func compute(a, b int, c string) (int, error) {

	var value interface{} = c
	switch value.(type) {
	case string:
		fmt.Println("string type")
	case *string:
		fmt.Println("string ptr type")
	case int:
		fmt.Println("int type")
	default:
		fmt.Println("error type")
	}

	switch c {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		return a / b, nil
	default:
		return -1, errors.New("操作符不合法")
	}
}
