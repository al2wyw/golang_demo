package main

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/shopspring/decimal"
	"testing"
	"time"
)

/*
如果使用map[string]interface{}来反序列化:
bool 代表 JSON booleans,
float64 代表 JSON numbers,
string 代表 JSON strings,
nil 代表 JSON null
*/

type Time time.Time

const (
	timeFormart = "2006-01-02"
)

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormart+`"`, string(data), time.Local)
	*t = Time(now)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormart)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormart)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	return time.Time(t).Format(timeFormart)
}

// Data 也可以利用Data来实现 UnmarshalJSON 方法来实现特殊的反序列化逻辑
type Data struct {
	Status  int             `json:"status"  valid:"required,range(1|10)~invalid pageSize"`
	Name    string          `json:"name"`
	Gender  bool            `json:"gender"`
	Amount  decimal.Decimal `json:"amount"` //Decimal已经自己实现了 UnmarshalJSON 方法
	Birth   Time            `json:"birth"`
	Address *Address        `json:"address,omitempty"`
	Values  []int           `json:"values"`
}

type Address struct {
	District string `json:"district"`
	Street   string `json:"street"`
}

func TestJson(t *testing.T) {
	testJson()
}

func testJson() {
	var data = []byte(`{"status": 0, "birth": "2019-01-01", "gender": true, "amount": 200.02}`)
	res := make(map[string]interface{})
	if err := json.Unmarshal(data, &res); err != nil {
		fmt.Println("json deserialize error", err)
		return
	}

	fmt.Println("json deserialize", res)

	var person Data
	if err := json.Unmarshal(data, &person); err != nil {
		fmt.Println("json deserialize error", err)
		return
	}
	ok, err := govalidator.ValidateStruct(person)
	if !ok || err != nil {
		fmt.Println("error person", person)
	}

	str, _ := json.Marshal(person)
	fmt.Println("json serialize", string(str))

	person.Values = []int{}
	str, _ = json.Marshal(person)
	fmt.Println("json serialize", string(str))

	str, _ = json.Marshal(struct{}{})
	fmt.Println("json serialize", string(str))

	str, _ = json.Marshal(nil)
	fmt.Println("json serialize", string(str))
}

type CommonReq struct {
	Interface interface{} `json:"interface,omitempty"`
	Timestamp int64       `json:"timestamp"`
	RequestId string      `json:"requestId"`
}

type Interface1 struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type Interface2 struct {
	Code  string `json:"code"`
	Value int64  `json:"value"`
}

// CommonResp 反序列化的第一种实现方法
type CommonResp struct {
	Data          interface{} `json:"data,omitempty"`
	Version       string      `json:"version"`
	ReturnCode    int         `json:"returnCode"`
	ReturnMessage string      `json:"returnMessage,omitempty"`
}

type Data1 struct {
	Name   string `json:"name"`
	Gender bool   `json:"gender"`
}

type Data2 struct {
	IdList []int64 `json:"idList"`
}

// CommonBaseResp 反序列化的第二种实现方法
type CommonBaseResp struct {
	Version       string `json:"version"`
	ReturnCode    int    `json:"returnCode"`
	ReturnMessage string `json:"returnMessage,omitempty"`
}

type Data1Resp struct {
	CommonBaseResp
	Data Data1 `json:"data"`
}

type Data2Resp struct {
	CommonBaseResp
	Data Data2 `json:"data"`
}

func TestCommonStructJson(t *testing.T) {
	req := &CommonReq{}
	req.RequestId = "dsfdf-344-dfgg-3444"
	req.Timestamp = time.Now().Unix()
	req.Interface = &Interface1{
		Name: "Peter",
		Id:   "id card number",
	}
	marshal(req)
	req.Interface = &Interface2{
		Code:  "valid code",
		Value: 3455555,
	}
	marshal(req)

	var data1 = []byte(`{"returnCode": 200, "ReturnMessage": "OK", "Version": "1.0", "data": {"name":"Peter","gender":true}}`)
	var data2 = []byte(`{"returnCode": 200, "ReturnMessage": "OK", "Version": "1.0", "data": {"idList":[213423,421343,4333]}}`)
	resp := &CommonResp{}
	dat1 := &Data1{}
	unmarshal(data1, resp)
	fmt.Println("json deserialize", dat1.Name)

	resp.Data = dat1
	unmarshal(data1, resp)
	fmt.Println("json deserialize", dat1.Name)

	dat2 := &Data2{}
	resp.Data = dat2
	unmarshal(data2, resp)
	fmt.Println("json deserialize", dat2.IdList)

	data1Resp := &Data1Resp{}
	unmarshal(data1, data1Resp)
	fmt.Println("json deserialize", data1Resp.Data.Name)

	data2Resp := &Data2Resp{}
	unmarshal(data2, data2Resp)
	fmt.Println("json deserialize", data2Resp.Data.IdList)
}

func marshal(req interface{}) {
	str, _ := json.Marshal(req)
	fmt.Println("json serialize", string(str))
}

func unmarshal(req []byte, resp interface{}) {
	err := json.Unmarshal(req, resp)
	if err != nil {
		fmt.Println("json deserialize error", err)
	}
	fmt.Println("json deserialize", resp)
}
