package main

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
)

func main() {

}

type response struct {
	url string
	err error
	res []byte
}

func request(ctx context.Context, url string, ch chan<- response) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		ch <- response{url, err, nil}
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- response{url, err, nil}
		return
	}
	defer resp.Body.Close() // always close

	if resp.StatusCode != http.StatusOK {
		ch <- response{url, errors.New(resp.Status), nil}
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- response{url, errors.New("read body failed"), nil}
		return
	}

	ch <- response{url, nil, body}
}

func compute(a, b int, c string) (int, error) {
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
