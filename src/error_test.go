package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

type respData struct {
	Status int             `json:"status"`
	Name   string          `json:"name"`
	Gender bool            `json:"gender"`
	Amount decimal.Decimal `json:"amount"` //Decimal已经自己实现了 UnmarshalJSON 方法
}

func TestError(t *testing.T) {
	err := requestA("", "http://localhost:8000/test")
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	var target *url.Error
	if errors.As(err, &target) {
		fmt.Printf("find url error %s\n", target)
	}
	var targetInf net.Error
	if errors.As(err, &targetInf) {
		fmt.Printf("find net error %s\n", targetInf)
	}
}

func requestA(requestBody string, url string) error {
	client := http.Client{Timeout: time.Millisecond * 600}
	resp, err := client.Post(url, "application/json", bytes.NewReader(([]byte)(requestBody)))
	if err != nil {
		return errors.Wrapf(err, "http request error %s", requestBody)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "http read error %s", requestBody)
	}
	defer resp.Body.Close()

	var model respData
	err = json.Unmarshal(respBody, &model)
	if err != nil {
		return errors.Wrapf(err, "Unmarshal error %s", respBody)
	}
	fmt.Println("http result", model)
	return nil
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
