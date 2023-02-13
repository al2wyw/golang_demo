package main

import (
	"fmt"
	"testing"
)

type Person struct {
	Name   string
	Age    byte
	IsDead bool
}

func TestMap(t *testing.T) {
	var myNilDic map[string]Person
	fmt.Println(myNilDic)
	for _, val := range myNilDic {
		fmt.Println(val)
	}
	fmt.Println(myNilDic["test"])
	if val, ok := myNilDic["test"]; ok {
		fmt.Println("myNilDic ok", val, ok)
	}
	// myNilDic["test"] = Person{} //panic

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
	dict1["test1"] = struct {
		Name string
	}{"peter"}

	fmt.Println("dict1", dict1["test"])
	fmt.Println("dict1", dict1["test1"])
	name := dict1["test"].Name
	fmt.Println("dict1 name", name)

	dict3 := make(map[string]*student)
	dict3["test"] = &student{Name: "Peter", Age: 18, Id: "test1"}
	fmt.Println("show name : ", dict3["test"].ShowName())
	//不能直接使用struct，不然改不了值
	for name := range dict3 {
		if dict3[name].Name == "Peter" {
			dict3[name].Age = 19
		}
	}

	if target, exist := dict3["test"]; exist {
		fmt.Println("target ", target)
	}

	fmt.Println("show dict3 name : ", dict3["test"])

	dict4 := make(map[Person]string)

	dict4[Person{
		"peter", 18, false,
	}] = "test1"

	dict4[Person{
		"ken", 48, true,
	}] = "test2"

	anne := Person{
		"anne", 38, true,
	}
	dict4[anne] = "test3"

	fmt.Println("show anne's name : ", dict4[anne])
}
