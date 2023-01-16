package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

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
}
