package main

import (
	"fmt"
	"github.com/jinzhu/copier"
	"math"
	"strconv"
	"testing"
	"time"
	"unsafe"
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
}

//（1）如果方法需要修改receiver，那么必须使用指针方法；
//（2）如果receiver是一个很大的结构体，考虑到效率，应该使用指针方法；
//（3）一致性。如果一些方法必须是指针receiver，那么其它方法也应该使用指针receiver；
//（4）对于一些基本类型、切片、或者小的结构体，使用value receiver效率会更高一些。
type Book struct{}

//值方法，Book为 receiver type
func (b Book) SetPages() {
	fmt.Println("SetPages")
}

//指针方法，*Book为 receiver type
func (b *Book) Pages() {
	fmt.Println("Pages")
}

func TestMethodCall(t *testing.T) {
	//
	var bb Book
	// 显式调用:
	Book.SetPages(bb)
	//Book.Pages(bb) //invalid method expression Book.Pages (needs pointer receiver: (*Book).Pages)
	(*Book).SetPages(&bb) //go会自动生成指针方法
	(*Book).Pages(&bb)

	b := Book{}
	bptr := &b

	b.SetPages()    // SetPages
	bptr.SetPages() // SetPages golang会将其解释为*bptr.SetPages()

	b.Pages()    // Pages golang会将其解释为&b.Pages()，但是必须是可寻址的值类型: 变量, 可寻址的数组(数组赋值给变量)元素，可寻址的结构体(结构体赋值给变量)字段 ，切片的元素 ，指针引用等(map元素不可寻址)，(Book{}).Pages() 会报错，字面量不可寻址
	bptr.Pages() // Pages
	//Book{}.Pages()

	fmt.Println("book slice test start")
	books := make([]Book, 3)
	books[0].Pages()

	bookDict := make(map[string]Book)
	bookDict["k1"] = Book{}
	//bookDict["k1"].Pages()
}

//多重继承
type Car struct {
	Name string
	Age  int
}

func (c *Car) Set(age int) {
	fmt.Println("car set called")
	c.Age = age
}

func (c *Car) Show() {
	fmt.Println("car is shown", c)
}

type Car2 struct {
	Name string
}

func (c *Car2) Action() {
	fmt.Println("car is action", c.Name)
}

//Go有匿名字段特性
type Train struct {
	Car //强依赖Car struct类型，如果只是依赖Show()和Set(age int)两个方法，可以改为interface 参考Vehicle
	*Car2
	createTime time.Time
	//age int   正常写法，Go的特性可以写成
	int
}

//给Train加方法，t指定接受变量的名字，变量可以叫this，t，p
func (t *Train) Set(age int) {
	fmt.Println("train set called")
	t.int = age
}

func TestDerived(t *testing.T) {
	var train Train
	train.int = 10 //这里用的匿名字段写法
	train.Set(1000)

	train.Car2 = new(Car2) //一定要先初始化，因为是指针
	train.Car2.Name = "test for car2"
	train.Action()

	train.Car.Set(1)
	train.Car.Name = "test" //这里Name必须得指定结构体

	train.Show()

	fmt.Println(train)
}

//如果接口不想被外部包实现，可以增加一个私有方法，但是可以用匿名嵌套的方式破解
type Vehicle interface {
	Engine()
	//private()
}

type EngineImpl struct {
	Name string
}

//Vehicle如果有两个方法，EngineImpl只实现一个方法的话不能算实现了接口，必须两个方法都实现
func (c *EngineImpl) Engine() {
	fmt.Println(c.Name, " has engine")
}

type DefaultEngine struct {
	Name string
}

func (c DefaultEngine) Engine() {
	fmt.Println(c.Name, " has default engine")
}

type Bus struct {
	Vehicle
	Pay float32
}

func (c *Bus) Aboard() {
	fmt.Println("welcome aboard")
}

func TestMethodReceiver(t *testing.T) {
	var engine = &EngineImpl{"normal engine"}
	var bus = Bus{engine, 323.22}

	bus.Engine()
	//bus.Vehicle = EngineImpl{"test engine"} //不行
	bus.Vehicle = &DefaultEngine{"engine"}
	bus.Engine()
	bus.Vehicle = DefaultEngine{"special engine"}
	bus.Engine()
	bus.Aboard()
}

/*
var veh Vehicle = Engine{}: 调用方是值
var veh Vehicle = &Engine{}: 调用方是指针
调用方	接收方	能否编译
值		值		true
值		指针		false
指针		值		true
指针		指针		true
指针		指针和值	true
值		指针和值	false
调用方是值时，只要接收方有指针方法那不允许编译，赋值语句报错
*/

func TestCopier(t *testing.T) {
	car1 := Car{Name: "car"}
	var car2 Car2
	err := copier.Copy(&car2, &car1)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(car2)
}

type Nameable interface {
	MyName() string
}

type fruit struct {
	Name string
}

func (f *fruit) MyName() string {
	return f.Name + " from fruit"
}

func (f *fruit) getFruit() Nameable {
	return f
}

type apple struct {
	fruit
	Price float64
}

func (f *apple) MyName() string {
	return f.Name + " from apple"
}

type orange struct {
	fruit
	Color string
}

func (f *orange) MyName() string {
	return f.Name + " from orange"
}

func TestStructReplace(t *testing.T) {
	ap := apple{fruit{Name: "apple"}, 234.44}
	or := orange{fruit{"orange"}, "red"}

	fmt.Println(ap.getFruit().MyName(), ap.MyName()) //from fruit from apple !!!

	printFruitPtr((*fruit)(unsafe.Pointer(&ap)))
	printFruitPtr((*fruit)(unsafe.Pointer(&or)))

	//printFruit(ap) //invalid
	printFruit(fruit{"fruit"})

	printFruitInt(&ap) //from apple
	printFruitInt(&or) //from orange
}

func printFruitPtr(tar *fruit) {
	fmt.Println(tar.Name, tar.MyName()) //from fruit
}

func printFruit(tar fruit) {
	fmt.Println(tar.Name, tar.MyName()) //from fruit
}

func printFruitInt(tar Nameable) {
	fmt.Println(tar, tar.MyName())
}
