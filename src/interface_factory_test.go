package main

import (
	"fmt"
	"reflect"
	"testing"
)

type Factory interface {
	Produce() bool
	Consume() bool
}

var _ Factory = (*SoftDrinkFactory)(nil)

type CafeFactory struct {
	ProductName string
}

type DrinkFactory struct {
	ProductName string
}

type SoftDrinkFactory struct {
	DrinkFactory
	SoftType string
}

func (d *DrinkFactory) Produce() bool {
	fmt.Printf("DrinkFactory生产%s成功\n", d.ProductName)
	d.ProductName = "new name"
	return true
}

func (d *DrinkFactory) Consume() bool {
	fmt.Printf("DrinkFactory消费%s成功\n", d.ProductName)
	return true
}

func (d *DrinkFactory) Show() bool {
	fmt.Printf("DrinkFactory show%s成功\n", d.ProductName)
	return true
}

func (d *SoftDrinkFactory) Show() bool {
	fmt.Printf("SoftDrinkFactory show%s成功\n", d.ProductName)
	return true
}

func (d *SoftDrinkFactory) Consume() bool {
	fmt.Printf("SoftDrinkFactory消费%s成功\n", d.ProductName)
	return true
}

func (c *CafeFactory) Produce() bool {
	fmt.Printf("CafeFactory生产%s成功\n", c.ProductName)
	return true
}

func (c *CafeFactory) Consume() bool {
	fmt.Printf("CafeFactory消费%s成功\n", c.ProductName)
	return true
}

func doProduce(factory Factory) bool {
	return factory.Produce()
}

func doConsume(factory Factory) bool {
	return factory.Consume()
}

func TestFactory(tt *testing.T) {
	factory := &SoftDrinkFactory{DrinkFactory{"Drink"}, "Cola"}
	doProduce(factory)
	doConsume(factory)

	var inter interface{} = *factory
	var interPtr interface{} = factory
	testInterPointer(interPtr)
	testInter(inter)

	var t = inter.(SoftDrinkFactory)
	fmt.Println("SoftType", t.SoftType)

	var fact Factory = factory
	fact.Consume()

	dmap := make(map[string]interface{})
	dmap["test"] = SoftDrinkFactory{DrinkFactory{"Drink"}, "Pepsi"}
	sd, ok := dmap["test"].(SoftDrinkFactory)
	if !ok {
		return
	}
	sd.Show()
}

func testInterPointer(inter interface{}) {
	if _, ok := inter.(*SoftDrinkFactory); ok {
		fmt.Println("*SoftDrinkFactory")
	}
	if _, ok := inter.(*DrinkFactory); ok {
		fmt.Println("*DrinkFactory")
	}
	if _, ok := inter.(Factory); ok {
		fmt.Println("*Factory")
	}
}

func testInter(inter interface{}) {
	if _, ok := inter.(SoftDrinkFactory); ok {
		fmt.Println("SoftDrinkFactory")
	}
	if _, ok := inter.(DrinkFactory); ok {
		fmt.Println("DrinkFactory")
	}
	if _, ok := inter.(Factory); ok {
		fmt.Println("Factory")
	}
}

func TestReflect(t *testing.T) {
	s := reflect.Int.String()
	fmt.Println(s) //int

	factory := SoftDrinkFactory{
		DrinkFactory{"Juice-Product"},
		"Juice"}

	//需要使用指针来驱动，槽
	klass := reflect.TypeOf(&factory)
	fmt.Println(klass.MethodByName("Consume"))

	kind_type_value(factory)

	rv := reflect.ValueOf(factory)
	prv := reflect.ValueOf(&factory)
	factory.Consume()
	methodCall(rv, "Consume") //method Consume not found
	methodCall(prv, "Consume")

	reflect.ValueOf(&factory.ProductName).Elem().SetString("test") //field set

	methodCall(prv, "Consume")
}

//kind 和 type 自定义类型必然不同 kind 比 type 更抽象
func kind_type_value(v interface{}) {
	rv := reflect.ValueOf(v)
	fmt.Println(rv.Kind(), rv.Type())
}

func methodCall(rv reflect.Value, method string, in ...reflect.Value) {
	// 通过方法名称获取 反射的方法对象
	mv := rv.MethodByName(method)
	// check mv 是否存在
	if !mv.IsValid() {
		fmt.Printf("mv is zero value, method %s not found\n", method)
		return
	}
	// 调用
	// nil 这里代表参数
	mv.Call(in)
}
