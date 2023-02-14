package sync_tests

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestWait(t *testing.T) {
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
	ch <- "cmd.2"
	close(ch)
	wg.Wait()
}

var compactCount int64 = 0
var lock = new(sync.Mutex)

func TestMutex(t *testing.T) {
	wg := new(sync.WaitGroup)
	loop := 100
	wg.Add(loop)

	fmt.Println(runtime.GOMAXPROCS(0), runtime.NumCPU())
	for i := 0; i < loop; i++ {
		go func() {
			//lock.Lock()
			defer wg.Done() //; lock.Unlock()
			compactCount++
			fmt.Println("compactCount", compactCount) //没有锁顺序错乱但是结果是对的，没有出现互相覆盖， wg有同步能力???
		}()
	}

	wg.Wait()
}

func TestSchedule(t *testing.T) {
	runtime.GOMAXPROCS(1)
	go func() {
		for {
		}
	}()
	time.Sleep(1 * time.Second)
	// can not run ???
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
		}
	}()

	time.Sleep(3 * time.Second)
}
