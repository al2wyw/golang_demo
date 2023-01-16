package main

import (
	"fmt"
	"github.com/shopspring/decimal"
	"testing"
)

func TestDecimal(t *testing.T) {
	testDecimal()
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
