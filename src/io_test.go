package main

import (
	"fmt"
	"golang_demo/src/io_utils"
	"io/ioutil"
	"os"
	"testing"
)

func TestIO(t *testing.T) {
	io_utils.ListDirs()
	file, err := io_utils.GetState("./io_test.go")
	if err != nil {
		fmt.Println("file state error", err)
		return
	}
	fmt.Println("file", file)

	exec, err := os.Executable()
	if err != nil {
		fmt.Println("exec error", err)
		return
	}
	fmt.Println("exec", exec)

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("work dir error", err)
		return
	}
	fmt.Println("work dir", wd)

	ret, err := ioutil.ReadFile("../note/goland.txt")
	if err != nil {
		fmt.Println("read file error", err)
		return
	}
	fmt.Println(string(ret[:100]))
}
