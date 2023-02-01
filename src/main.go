package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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
}

type response struct {
	url string
	err error
	res []byte
}

func request(ctx context.Context, url string, ch chan<- *response) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		ch <- &response{url, err, nil}
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- &response{url, err, nil}
		return
	}
	defer resp.Body.Close() // always close

	if resp.StatusCode != http.StatusOK {
		ch <- &response{url, errors.New(resp.Status), nil}
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- &response{url, errors.New("read body failed"), nil}
		return
	}

	ch <- &response{url, nil, body}
}

func compute(a, b int, c string) (int, error) {

	var value interface{} = c
	switch value.(type) {
	case string:
		fmt.Println("string type")
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
