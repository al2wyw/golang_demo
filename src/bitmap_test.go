package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"testing"
)

func TestBitmap(t *testing.T) {

	rb1 := roaring.BitmapOf(1, 2, 3, 4, 5, 100, 1000)
	rb2 := roaring.BitmapOf(200000010, 200000020, 200000030, 100000060, 100000070, 100000080, 100000090, 100000100, 200000040, 200000050, 100000110, 100000120, 9020032, 9020033, 9020036, 9020037, 9010029, 9010030, 9010031, 9010032, 9010033, 9010040, 9010041, 9010042, 9010043)

	fmt.Println(rb1.String())
	fmt.Println(rb2.String())

	rb2.Add(400000020)

	i := rb2.Iterator()
	for i.HasNext() {
		fmt.Println(i.Next())
	}

	buf := new(bytes.Buffer)
	rb2.WriteTo(buf) // we omit error handling

	fmt.Println(len(buf.Bytes()))
	fmt.Println(base64.StdEncoding.EncodeToString(buf.Bytes()))

	newrb := roaring.New()
	newrb.ReadFrom(buf)

	if rb2.Equals(newrb) {
		fmt.Println("They are the same")
	}
}
