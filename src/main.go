package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"sync"
	"time"
)

type empty int

func (em *empty) doNothing(v interface{}) {
	if v == nil {
		return
	}
	fmt.Printf("doNothing %p 成功\n", &v)
}

type Func interface {
	doFun(v interface{})
}

type FuncHandle func(interface{})

func (fun FuncHandle) doFunc(v interface{}) {
	fun(v)
}

func FuncCaller(v interface{}, fun func(interface{})) {
	h := FuncHandle(fun)
	h.doFunc(v)
}

type Factory interface {
	Produce() bool
	Consume() bool
}

var _ Factory = (*DrinkFactory)(nil)

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

var test func(int, int, string) (int, error)

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

func main() {

	var em empty = 1
	var val int = 1
	em++

	em.doNothing(val)
	FuncCaller(val, em.doNothing)
	//fmt.Println(em == val) //em 的值是int，但是类型不是int

	var bigVal int64 = 2
	fmt.Println(bigVal)
	//bigVal = val

	var data *byte
	var in interface{}

	fmt.Println(data, data == nil) // <nil> true
	fmt.Println(in, in == nil)     // <nil> true

	in = data
	fmt.Println(in, in == nil) // <nil> false

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

	array := []int{1, 2, 3, 4}
	test1 := array[0]
	testPtr(&array[0])
	testArray(array)

	for _, t := range array {
		fmt.Println("array elements : ", t)
	}

	for i := 0; i < 10; i++ {
		func() { fmt.Println("loop", i) }() // 输出全部为10
	}

	testReSlice()

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

	test2 := 234.34
	fmt.Println("test2", test2)
	person := student{
		Name: "test",
		Age:  10,
		Id:   "test_id1",
	}

	fmt.Println("show name : ", (&person).ShowName())
	fmt.Println("show name : ", person.ShowName())

	myName := string("123")
	fmt.Println(array, reflect.TypeOf(array))
	fmt.Println(person, reflect.TypeOf(person))
	fmt.Println(myName)

	fmt.Println(test1)
	testPointer(person)
	fmt.Println(person)

	testTypeAssert(1)
	testTypeAssert("test")

	testDecimal()
	testStrTsf()

	testDefer()
	testJson()

	testWaitG()
	testBufferedChan()

	var duck Duck = &Cat{}
	duck.Quack()
}

func testDecimal() error {
	dec, err := decimal.NewFromString("234.5466")
	if err != nil {
		return err
	}

	ret := dec.Add(decimal.NewFromInt(43)).Round(2)

	fmt.Println("decimal ret", ret)
	return nil
}

type Duck interface {
	Quack()
}

type Cat struct{}

func (c *Cat) Quack() {
	fmt.Println("meow")
}

type response struct {
	url string
	err error
	res []byte
}

func request(ctx context.Context, url string, ch chan<- response) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		ch <- response{url, err, nil}
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- response{url, err, nil}
		return
	}
	defer resp.Body.Close() // always close

	if resp.StatusCode != http.StatusOK {
		ch <- response{url, errors.New(resp.Status), nil}
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- response{url, errors.New("read body failed"), nil}
		return
	}

	ch <- response{url, nil, body}
}

func testDefer() {
	var i = 1
	//预先设置好引用的变量
	defer fmt.Println("result: ", func() int { return i * 2 }()) //2
	i++
	fmt.Println("main result: ", i) //2
}

func testReSlice() {
	fix := []int{1, 2, 3}
	slis := fix[1:3]
	fmt.Println("slis len cap", len(slis), cap(slis))
	slis[0] = 10
	//slis[3] = 100 index 错误
	fmt.Printf("slis addr %p\n", slis)
	slis = append(slis, 100) //append完生成新的slice
	fmt.Printf("slis addr %p\n", slis)
	slis[0] = 101 //不再改变原来的fix
	fmt.Println("fix", fix)
	fmt.Println("slis", slis)

	slis2 := make([]int, 3, 5)
	//slis2[3] = 10 index 错误
	fmt.Printf("slis2 addr %p\n", slis2)
	slis2 = append(slis2, 10)
	fmt.Printf("slis2 addr %p\n", slis2)
	fmt.Println("slis2", slis2)

	path := []byte("AAAA/BBBBBBBBB")
	sepIndex := bytes.IndexByte(path, '/') // 4
	println(sepIndex)

	dir1 := path[:sepIndex]
	dir2 := path[sepIndex+1:]
	println("dir1: ", string(dir1)) // AAAA
	println("dir2: ", string(dir2)) // BBBBBBBBB

	dir1 = append(dir1, "suffix"...)        //神奇
	println("current path: ", string(path)) //AAAAsuffixBBBB
}

func testPip() {
	pip := make(chan int, 100)
	pip <- 100

	go func() {
		count := 10
		for i := 0; i < count; i++ {
			val, ok := <-pip
			time.Sleep(1000 * time.Millisecond)
			fmt.Println("wake up start to work", val, ok)
		}
	}()
	close(pip)
	time.Sleep(60 * time.Second)
}

func testTypeAssert(object interface{}) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic error", err)
		}
	}()

	t, ok := object.(string)
	fmt.Println("this is type assert", ok)

	if t == "test" {
		fmt.Println("this is a test")
	}
}

func testStrTsf() {
	s := "我是ABC"
	//s[0] is uint8
	fmt.Println("this is a utf8 str:", s[0], s[6])

	for _, v := range s {
		//v is int32
		fmt.Println("this is a range:", v)
		fmt.Println("this is a utf8 range:", string(v))
	}

	arr := []rune(s)
	arr[0] = rune('他') //Redundant type conversion
	fmt.Println("this is a utf8 str:", string(arr))

}

func testJson() {
	var data = []byte(`{"status": 200, "birth": "2019-01-01", "gender": true, "amount": 200.02}`)
	res := make(map[string]interface{})
	if err := json.Unmarshal(data, &res); err != nil {
		fmt.Println("json deserialize error", err)
		return
	}

	fmt.Println("json deserialize", res)
}

func testBufferedChan() {
	ch := make(chan int, 10)
	for i := 0; i < 10; i++ {
		ch <- i
	}

	go func() {
		for {
			i, ok := <-ch
			if ok {
				fmt.Println("this is value from chan", i)
			} else {
				fmt.Println("chan is false")
				return
			}
		}
	}()

	time.Sleep(100 * time.Millisecond)
	close(ch)
	time.Sleep(100 * time.Millisecond)
}

func testWaitG() {
	ch := make(chan string)
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		defer wg.Done()
		//ch不close, range就不会结束
		for m := range ch {
			fmt.Println("Processed:", m)
			time.Sleep(100 * time.Millisecond) // 模拟需要长时间运行的操作
		}
	}()

	ch <- "cmd.1"
	ch <- "cmd.2" // 不会被接收处理
	close(ch)
	wg.Wait()
}

func testSelect() {
	ch := make(chan int)
	timeout := make(chan bool)
	go func() {
		for {
			select {
			case v := <-ch:
				fmt.Println("read from ch", v)
			case <-timeout:
				fmt.Println("read from timeout")
				//default 一直执行
				//default:
				//	fmt.Println("this is default")
			}
		}
	}()
	_ = <-time.After(10 * time.Second)
	timeout <- true
	fmt.Println("testSelect end")
}

func testPtr(a *int) {
	*a = 10
}

func testPointer(p student) {
	p.Name = "test again"
	fmt.Println(p)
}

func testArray(p []int) {
	p[0] = 100
	fmt.Println("test array", p[0])
	var array []student
	array = append(array, student{Name: "test"})

	data, err := json.Marshal(array)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(string(data))
}

func compute(a, b int, c string) (int, error) {
	switch c {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		return a / b, nil
	default:
		return -1, errors.New("操作符不合法")
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
