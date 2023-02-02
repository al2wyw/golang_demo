package main

import (
	"fmt"
	"github.com/shopspring/decimal"
	"math/big"
	"strings"
	"testing"
)

func TestDecimal(t *testing.T) {

	text := big.NewInt(int64(1000180)).Text(36)
	text = strings.Repeat("0", 8-len(text)) + text
	fmt.Println("36 text", text)

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
