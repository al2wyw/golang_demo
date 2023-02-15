package main

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

//自动继承
type Base1 interface {
	Produce() bool
	Safe() bool
}

type Base2 interface {
	Produce() bool
	Safe() bool
	Write() bool
}

type MyImp int

func (m MyImp) Produce() bool {
	panic("implement me")
}

func (m MyImp) Safe() bool {
	panic("implement me")
}

func (m MyImp) Write() bool {
	panic("implement me")
}

func TestAutoExtend(t *testing.T) {
	var u MyImp = 10
	var base1 Base1 = u
	var base2 Base2 = u
	fmt.Println(base1, base2)

	base1 = base2
	fmt.Println(base1)
}

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
	SoftType string `json:"age,omitempty"`
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

// Method reflect method test
func (d SoftDrinkFactory) Method() bool {
	fmt.Printf("SoftDrinkFactory方法调用%s成功\n", d.ProductName)
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

type Duck interface {
	Quack()
}

type Cat struct{}

// Cat implements Duck
func (c *Cat) Quack() {
	fmt.Println("meow")
}

func TestFactory(tt *testing.T) {
	factory := &SoftDrinkFactory{DrinkFactory{"Drink"}, "Cola"}
	doProduce(factory)
	doConsume(factory)

	testFactory(factory)
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

	var duck Duck = &Cat{}
	duck.Quack()

	fmt.Printf("duck type %T", duck)
}

func testFactory(inter Factory) {
	if _, ok := inter.(*SoftDrinkFactory); ok {
		fmt.Println("*SoftDrinkFactory")
	}
	if _, ok := inter.(*DrinkFactory); ok {
		fmt.Println("*DrinkFactory")
	}
	if _, ok := inter.(Factory); ok {
		fmt.Println("Factory")
	}
}

func testInterPointer(inter interface{}) {
	if _, ok := inter.(*SoftDrinkFactory); ok {
		fmt.Println("*SoftDrinkFactory")
	}
	if _, ok := inter.(*DrinkFactory); ok {
		fmt.Println("*DrinkFactory")
	}
	if _, ok := inter.(Factory); ok {
		fmt.Println("Factory")
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
	fmt.Printf("factory addr %p\n", &factory)

	//指针类型可以获取到指针方法和值方法，不可以获取到字段；值类型可以获取到值方法和字段
	class := reflect.TypeOf(factory)
	fmt.Println(class.FieldByName("SoftType"))
	fmt.Println(class.MethodByName("Method"))
	fmt.Println(class.MethodByName("Consume"))

	for i := 0; i < class.NumMethod(); i++ {
		m := class.Method(i)
		fmt.Printf("%s: %v\n", m.Name, m.Type)
	}

	//需要使用指针来驱动，因为Consume是指针方法
	klass := reflect.TypeOf(&factory)
	//fmt.Println(klass.FieldByName("SoftType")) //panic
	fmt.Println(klass.MethodByName("Method"))
	fmt.Println(klass.MethodByName("Consume"))

	for i := 0; i < klass.NumMethod(); i++ {
		m := klass.Method(i)
		fmt.Printf("%s: %v\n", m.Name, m.Type)
	}

	kind_type_value(factory)
	DoFiledAndMethod(factory)

	rv := reflect.ValueOf(factory)
	prv := reflect.ValueOf(&factory)

	factory.Consume()
	methodCall(rv, "Method")
	methodCall(rv, "Consume") //method Consume not found
	methodCall(prv, "Method") //这个怎么会调的到???
	methodCall(prv, "Consume")

	//必须是指针，才是addressable的
	reflect.ValueOf(&(factory.ProductName)).Elem().SetString("test") //field set

	methodCall(prv, "Consume")

	ptr := uintptr(unsafe.Pointer(&factory)) + unsafe.Offsetof(factory.ProductName)
	test := (*string)(unsafe.Pointer(ptr)) //这个回转pinter为什么告警，但是这样就不会告警: test := (*string)(unsafe.Pointer(uintptr(unsafe.Pointer(&factory)) + unsafe.Offsetof(factory.ProductName)))
	//栈扩容缩容导致ptr的地址值无效 ???

	*test = "test again"

	methodCall(prv, "Consume")
}

func TestInterface(t *testing.T) {
	var softFactory = SoftDrinkFactory{
		DrinkFactory{"Juice-Product"},
		"Juice"}
	var factory Factory = &softFactory

	fmt.Printf("addr %p, %p, %p \n", factory, &factory, unsafe.Pointer(&factory))
	fmt.Println(reflect.TypeOf(softFactory))
	fmt.Println(reflect.TypeOf(factory))

	//参考 调用方	接收方 的规则表格//
	var errorInfPtrType = reflect.TypeOf((*error)(nil))         //kind ptr
	var errorInfType = errorInfPtrType.Elem()                   //kind interface
	var errorT = reflect.TypeOf((error)(nil))                   // nil
	var errorStructPtrType = reflect.TypeOf(errors.New("test")) //kind ptr
	var errorStructType = errorStructPtrType.Elem()             //kind struct New return a ptr of struct
	if errorT != nil && errorT.AssignableTo(errorInfType) {
		fmt.Println("nothing")
	}
	if errorStructPtrType.AssignableTo(errorInfType) {
		fmt.Println("interface = &struct{}")
	}
	if errorStructType.AssignableTo(errorInfType) {
		fmt.Println("interface = struct{}")
	}
	//参考 调用方	接收方 的规则表格//
}

//kind 和 type 自定义类型必然不同 kind 比 type 更抽象
func kind_type_value(v interface{}) {
	rv := reflect.ValueOf(v)
	fmt.Println(rv.Kind(), rv.Type())
}

func DoFiledAndMethod(input interface{}) {

	getType := reflect.TypeOf(input)
	fmt.Println("get Type is :", getType.Name())

	getValue := reflect.ValueOf(input)
	fmt.Println("get all Fields is:", getValue)

	// 获取方法字段
	// 1. 先获取interface的reflect.Type，然后通过NumField进行遍历
	// 2. 再通过reflect.Type的Field获取其Field
	// 3. 最后通过Field的Interface()得到对应的value
	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i)
		value := getValue.Field(i).Interface()
		fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)

		tagValue := field.Tag.Get("json")
		if tagValue == "" {
			fmt.Printf("can not get tag value \n")
		}
	}

	// 获取方法
	// 1. 先获取interface的reflect.Type，然后通过.NumMethod进行遍历
	for i := 0; i < getType.NumMethod(); i++ {
		m := getType.Method(i)
		fmt.Printf("%s: %v\n", m.Name, m.Type)
	}
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
