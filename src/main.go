package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

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
	fmt.Println(ret)

	http.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		time.Sleep(time.Second)
		writer.WriteHeader(200)
		size, err := fmt.Fprint(writer, `{"name": "test", "status": 200, "birth": "2019-01-01", "gender": true, "amount": 200.02}`)
		if err != nil {
			fmt.Println("error", err)
		}
		fmt.Println("write byte", size)
	})
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println(err)
	}
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
