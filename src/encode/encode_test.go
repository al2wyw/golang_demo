package encode

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type AgeT int

func (a *AgeT) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", *a+10)), nil
}

type Reachable interface {
	Reach()
}

type Sub struct {
	Target string `encode:"target"`
}

type Data struct {
	Name   string    `encode:"name"`
	Age    AgeT      `encode:"age"`
	Date   time.Time `encode:"date,dateformat=2006-01-02 15:04:05"`
	Salary float64   `encode:"salary"`
	Valid  bool      `encode:"valid"`
	Sub
	no string `encode:"no"`
}

func (d Data) Reach() {
	panic("implement me")
}

type DataPtr *Data

func TestStructEncode(t *testing.T) {
	var data Data
	data.no = "no"
	data.Name = "test"
	data.Age = 10
	data.Valid = true
	data.Date = time.Now()
	data.Salary = 34.132
	data.Sub.Target = "sub"
	b, _ := Encode(data)
	fmt.Println(string(b))
	c, _ := Encode(&data) // MarshalText trigger
	fmt.Println(string(c))

	var reachable Reachable = data
	r, _ := Encode(reachable) //equivalent to b
	fmt.Println(string(r))
	r1, _ := Encode(&reachable) //rare case
	fmt.Println(string(r1))

	var reachableP Reachable = &data
	rp, _ := Encode(reachableP) // MarshalText trigger equivalent to c
	fmt.Println(string(rp))
	rp1, _ := Encode(&reachableP) //rare case
	fmt.Println(string(rp1))

	var dataP *Data = &data
	var prt DataPtr = dataP
	fmt.Printf("%s, %s\n", reflect.TypeOf(prt).Name(), reflect.TypeOf(dataP).Name())
	fmt.Println(prt)

}
