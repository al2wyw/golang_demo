package main

import (
	"fmt"
	"github.com/goinggo/mapstructure"
	"github.com/jinzhu/copier"
	"math"
	"reflect"
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

// （1）如果方法需要修改receiver，那么必须使用指针方法；
// （2）如果receiver是一个很大的结构体，考虑到效率，应该使用指针方法；
// （3）一致性。如果一些方法必须是指针receiver，那么其它方法也应该使用指针receiver；
// （4）对于一些基本类型、切片、或者小的结构体，使用value receiver效率会更高一些。
type Book struct{}

// 值方法，Book为 receiver type
func (b Book) SetPages() {
	fmt.Println("SetPages")
}

// 指针方法，*Book为 receiver type
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

	b.Pages()    // Pages golang会将其解释为&b.Pages()，但是必须是可寻址的值类型: 变量, 可寻址的数组(数组赋值给变量，通过变量和数组下标访问的元素)元素，可寻址的结构体(结构体赋值给变量)字段 ，切片的元素 ，指针引用等(map元素不可寻址)
	bptr.Pages() // Pages
	//Book{}.Pages()

	fmt.Println("book slice test start")
	books := make([]Book, 3)
	books[0].Pages()

	bookDict := make(map[string]Book)
	bookDict["k1"] = Book{}
	//bookDict["k1"].Pages()
}

type S struct {
	A int
	b int
}

func (s *S) canPrintPtr() {
	fmt.Println("pointer print", *s)
}

func (s S) canPrint() {
	fmt.Println("print", s)
}

// TestElemAddressable  只要值是可寻址的，就可以在值上直接调用指针方法， map的元素是不可寻址的(可以把元素声明成指针元素), slice的元素可寻址
func TestElemAddressable(t *testing.T) {
	dict := map[string]S{"A": {A: 1, b: 2}}
	//dict["A"].canPrintPtr()
	dict["A"].canPrint()

	slice := []S{{A: 1, b: 2}}
	slice[0].canPrintPtr()
	slice[0].canPrint()

	array := [1]S{{A: 1, b: 2}}
	array[0].canPrintPtr()
	array[0].canPrint()

	//字面量有点奇怪，不要理会:
	//S{A : 1, b: 2}.canPrintPtr()
	(&S{A: 1, b: 2}).canPrintPtr()
	S{A: 1, b: 2}.canPrint()
	//map[string]S{"A": {A: 1, b: 2}}["A"].canPrintPtr()
	map[string]S{"A": {A: 1, b: 2}}["A"].canPrint()
	[]S{{A: 1, b: 2}}[0].canPrintPtr()
	[]S{{A: 1, b: 2}}[0].canPrint()
	//[1]S{{A: 1, b: 2}}[0].canPrintPtr()
	[1]S{{A: 1, b: 2}}[0].canPrint()
}

// CanSet = CanAddr + exported struct field (struct的非导出字段FieldByName不可set)
// CanAddr = an element of A slice, an element of an addressable array, A field of an addressable struct, or the result of dereference A pointer
func TestReflectAddressable(t *testing.T) {
	var x = 8
	canSet(x) //false
	canSet(&x)
	var y = &x
	canSet(y)
	fmt.Println(x)

	s := S{A: 1, b: 2}
	canSet(&s.A)
	canSet(s.A) //false
	fmt.Println(s)

	ret := newCanAddr(reflect.TypeOf(S{}))
	fmt.Println(ret, ret.(*S).A)

	slice := []int{1, 2, 3}
	canAddr(slice) //true
	array := [3]int{1, 2, 3}
	canAddr(array)  //false
	canAddr(&array) //addressable array true
	dict := map[string]int{"test": 3}
	canAddr(dict) //false
	canAddr(s)    //false
	canAddr(&s)   //addressable struct true
}

func canSet(value interface{}) {
	fmt.Println("can set ......")
	rv := reflect.ValueOf(value)
	fmt.Println(rv.CanAddr(), rv.CanSet(), rv.Type())
	rv = reflect.Indirect(rv) //comment掉这一行全部case都是false
	fmt.Println(rv.CanAddr(), rv.CanSet(), rv.Type())
	if rv.CanSet() {
		rv.Set(reflect.ValueOf(23))
		fmt.Println(rv.Interface())
	}
}

func canAddr(value interface{}) {
	fmt.Println("can addr ......")
	rv := reflect.Indirect(reflect.ValueOf(value))

	var rvIndex reflect.Value
	switch rv.Type().Kind() {
	case reflect.Slice, reflect.Array:
		rvIndex = rv.Index(0)
	case reflect.Map:
		rvIndex = rv.MapIndex(reflect.ValueOf("test"))
	case reflect.Struct:
		rvIndex = rv.FieldByName("A")
	}
	fmt.Println(rvIndex.CanAddr(), rvIndex.CanSet(), rvIndex.Type())
	if rvIndex.CanAddr() {
		rvIndex.Set(reflect.ValueOf(33))
		//rvIndex.Addr().Type() == reflect.PointerTo(rvIndex.Type())
		fmt.Println(rvIndex.Addr().Type(), reflect.PointerTo(rvIndex.Type()))
	}
	fmt.Println(rv.Interface())
}

func newCanAddr(rt reflect.Type) interface{} {
	fmt.Println("new can addr ......")
	if rt.Kind() != reflect.Struct {
		return nil
	}
	ret := reflect.New(rt)
	fmt.Println(ret.CanAddr(), ret.CanSet(), ret.Type())
	rv := reflect.Indirect(ret)
	rv = rv.FieldByName("A")
	if rv.CanSet() {
		rv.Set(reflect.ValueOf(23))
	}
	return ret.Interface()
}

// 多重继承
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

// Go有匿名字段特性
type Train struct {
	Car //强依赖Car struct类型，如果只是依赖Show()和Set(age int)两个方法，可以改为interface 参考Vehicle
	*Car2
	createTime time.Time
	//age int   正常写法，Go的特性可以写成
	int
}

// 给Train加方法，t指定接受变量的名字，变量可以叫this，t，p
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

// 如果接口不想被外部包实现，可以增加一个私有方法，但是可以用匿名嵌套的方式破解
type Vehicle interface {
	Engine()
	//private()
}

type EngineImpl struct {
	Name string
}

// Vehicle如果有两个方法，EngineImpl只实现一个方法的话不能算实现了接口，必须两个方法都实现
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

	canImplementVehicle(EngineImpl{"test engine"})        //false
	canImplementVehicle(&EngineImpl{"test engine"})       //normal addr
	canImplementVehicle(DefaultEngine{"special engine"})  //normal
	canImplementVehicle(&DefaultEngine{"special engine"}) //normal indirect addr
}

func canImplementVehicle(value interface{}) {
	fmt.Println("canImplementVehicle ......")
	vehicleType := reflect.TypeOf((*Vehicle)(nil)).Elem()

	rv := reflect.ValueOf(value)
	if rv.Type().Implements(vehicleType) {
		fmt.Println("normal", rv.Type())
	}

	rv = reflect.Indirect(rv)
	if rv.CanAddr() {
		if rv.Type().Implements(vehicleType) {
			fmt.Println("indirect", rv.Type())
		}
		if rv.Addr().Type().Implements(vehicleType) {
			fmt.Println("addr", rv.Addr().Type())
		}
	}
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
	car1 := Car{Name: "car1"}
	var car2 Car2
	err := copier.Copy(&car2, car1)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(car2)

	car := &Car{}
	source := map[string]interface{}{
		"age":  1,
		"name": "value",
	}
	err = mapstructure.Decode(source, car)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(*car)
	source["age"] = 2
	err = copier.Copy(car, source) // not supported
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(*car)
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
