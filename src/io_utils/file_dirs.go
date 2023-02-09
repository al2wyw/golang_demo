package io_utils

import (
	"fmt"
	"os"
)

func ListDirs() {
	dirs, err := os.ReadDir("./")
	if err != nil {
		fmt.Println("read dirs error", err)
		return
	}
	for _, dir := range dirs {
		fmt.Println(dir)
	}
}

func GetState(file string) (string, error) {
	stat, err := os.Stat(file)
	if err != nil {
		return "", err
	}
	return stat.Mode().String(), nil
}
