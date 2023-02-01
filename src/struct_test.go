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

type Book struct{}

//值方法，Book为 receiver type
func (b Book) SetPages() {
	fmt.Println("SetPages")
}

//指针方法，*Book为 receiver type
func (b *Book) Pages() {
	fmt.Println("Pages")
}

func TestStruct(t *testing.T) {
	test2 := 234.34
	fmt.Println("test2", test2)

	var var1 student
	var1.Id = "test_id1"
	var1.Age = 10
	var1.Name = "Peter"
	fmt.Println(var1)

	personPtr := new(student)
	//底层会转换为(*personPtr).Id = "test_id1"
	personPtr.Id = "test_id1"
	personPtr.Age = 10
	personPtr.Name = "Peter"
	fmt.Println(personPtr)  //&{Peter 10 test_id1}
	fmt.Println(*personPtr) //{Peter 10 test_id1}

	person := student{
		Name: "test",
		Age:  10,
		Id:   "test_id1",
	}
	person.Id = "test_id2"

	personP := &person
	fmt.Println("show name : ", personP.ShowName())
	fmt.Println("show name : ", person.ShowName())

	testPointer(person)

	fmt.Println(person)

	//
	var bb Book
	Book.SetPages(bb)  // 显式调用
	(*Book).Pages(&bb) // 显式调用

	b := Book{}
	b1 := &b

	b.SetPages()  // SetPages
	b1.SetPages() // SetPages golang会将其解释为*b1.SetPages()

	b.Pages()  // Pages golang会将其解释为&b.Pages()，但是必须是可寻址的值类型: 变量, 可寻址的数组元素，可寻址的结构体字段 ，切片 ，指针引用等，(Book{}).Pages() 会报错，字面量不可寻址
	b1.Pages() // Pages
}
