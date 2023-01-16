package main

import (
	"fmt"
	"testing"
)

func TestMap(t *testing.T) {
	dict := make(map[string]string)
	dict["test"] = "1"
	dict["value"] = "2"
	fmt.Println("dict", dict)
	if val, ok := dict["ok"]; ok {
		fmt.Println("dict ok", val, ok)
	}

	dict1 := make(map[string]struct {
		Name string
	})
	dict1["test"] = struct {
		Name string
	}{"rdy"}
	fmt.Println("dict1", dict1["test"])
	name := dict1["test"].Name
	fmt.Println("dict1 name", name)
	dict3 := make(map[string]student)
	dict3["test"] = student{Name: "Peter", Age: 18, Id: "test1"}
	//fmt.Println("show name : " ,dict3["test"].ShowName())
}
