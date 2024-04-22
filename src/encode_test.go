package main

import (
	"fmt"
	"golang_demo/src/encode"
	"reflect"
	"strings"
	"testing"
	"time"
)

type AgeT int

func (a *AgeT) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", *a+10)), nil
}

type SensitiveStr string

func (str SensitiveStr) Encode() ([]byte, error) {
	return []byte(strings.ReplaceAll(string(str), "test", "****")), nil
}

type Reachable interface {
	Reach()
}

type Sub struct {
	Target string `encode:"target"`
}

type EData struct {
	Name   string       `encode:"name"`
	Age    AgeT         `encode:"age"`
	Date   time.Time    `encode:"date,dateformat=2006-01-02 15:04:05"`
	Salary float64      `encode:"salary"`
	Valid  bool         `encode:"valid"`
	IdCard SensitiveStr `encode:"idCard"`
	Encode string       `encode:"encode,encoder=customizedEncoder"`
	Sub
	no string `encode:"no"`
}

func (d EData) Reach() {
	panic("implement me")
}

type EDataPtr *EData

func TestStructEncode(t *testing.T) {
	var data EData
	data.no = "no"
	data.Name = "test"
	data.Age = 10
	data.Valid = true
	data.Date = time.Now()
	data.Salary = 34.132
	data.Encode = "test for encode"
	data.IdCard = "34443test445555"
	data.Sub.Target = "sub"
	b, _ := encode.Encode(data)
	fmt.Println(string(b))
	c, _ := encode.Encode(&data) // MarshalText trigger
	fmt.Println(string(c))

	var reachable Reachable = data
	r, _ := encode.Encode(reachable) //equivalent to b
	fmt.Println(string(r))
	r1, _ := encode.Encode(&reachable) //rare case
	fmt.Println(string(r1))

	var reachableP Reachable = &data
	rp, _ := encode.Encode(reachableP) // MarshalText trigger equivalent to c
	fmt.Println(string(rp))
	rp1, _ := encode.Encode(&reachableP) //rare case
	fmt.Println(string(rp1))

	var dataP *EData = &data
	var prt EDataPtr = dataP
	fmt.Printf("%s, %s\n", reflect.TypeOf(prt).Name(), reflect.TypeOf(dataP).Name())
	fmt.Println(prt)

}

func init() {
	err := encode.RegisterEncoder("customizedEncoder", func(obj interface{}) ([]byte, error) {
		if str, ok := obj.(string); ok {
			return []byte(fmt.Sprintf("/ %s /", str)), nil
		}
		return nil, nil
	})
	if err != nil {
		fmt.Println(err)
	}
}
