package sync_tests

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

type data struct {
	name string
}

var count = 0

func NewData() interface{} {
	count++
	return &data{
		"test" + strconv.Itoa(count),
	}
}

var dataPool = sync.Pool{
	New: NewData,
}

func TestPool(t *testing.T) {
	for loop := 1; loop < 3; loop++ {
		go func(i int) {
			for {
				myData := dataPool.Get().(*data) //随机获取，对象实例不会和协程绑定
				fmt.Printf("this is %d, %s, %p \n", i, myData.name, myData)
				dataPool.Put(myData)
				time.Sleep(100 * time.Millisecond)
			}
		}(loop)
	}

	time.Sleep(3000 * time.Millisecond)
}
