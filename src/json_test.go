package main

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"testing"
	"time"
)

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
	Status int             `json:"status"`
	Name   string          `json:"name"`
	Gender bool            `json:"gender"`
	Amount decimal.Decimal `json:"amount"` //Decimal已经自己实现了 UnmarshalJSON 方法
	Birth  Time            `json:"birth"`
}

func TestJson(t *testing.T) {
	testJson()
}

func testJson() {
	var data = []byte(`{"status": 200, "birth": "2019-01-01", "gender": true, "amount": 200.02}`)
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
	fmt.Println("json deserialize", person)
}
