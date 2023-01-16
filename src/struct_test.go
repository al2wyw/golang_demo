package main

import (
	"fmt"
	"math"
	"strconv"
	"testing"
)

type student struct {
	Name string `json:"name"`
	Age  int    `json:"age,omitempty"`
	Id   string `json:"id"`
}

var count = 0

func (s *student) ShowName() string {
	if count == math.MaxInt {
		count = 0
	}
	count++
	return s.Name + strconv.Itoa(count)
}

func testPointer(p student) {
	p.Name = "test again"
	fmt.Println(p)
}

func TestStruct(t *testing.T) {
	test2 := 234.34
	fmt.Println("test2", test2)
	person := student{
		Name: "test",
		Age:  10,
		Id:   "test_id1",
	}

	fmt.Println("show name : ", (&person).ShowName())
	fmt.Println("show name : ", person.ShowName())

	testPointer(person)

	fmt.Println(person)
}
