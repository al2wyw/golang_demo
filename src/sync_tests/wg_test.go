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
var ret float64
var loopO = 1000000

func TestMutex(t *testing.T) {
	Concurr(1) //只有一个cpu可以运行协程, 不需要考虑竞态， 但是协程互相切换时应该需要考虑脏读(python的多线程问题) ???
	fmt.Println("================")
	ConcurrV2(1)
	fmt.Println("================")
	//Concurr(2) //多个cpu同时运行协程，会出现竞态
}

func Concurr(concurrency int) {
	wg := new(sync.WaitGroup)
	loop := loopO
	wg.Add(loop)

	fmt.Println(runtime.GOMAXPROCS(concurrency), runtime.NumCPU())
	for i := 0; i < loop; i++ {
		go func() {
			defer wg.Done()
			ret = 10.34 / 2.34
			compactCount++
			//fmt.Println("compactCount", compactCount) //concurrency > 1 时 没有锁顺序错乱但是结果是对的，没有出现互相覆盖， wg有同步能力???
		}()
	}

	wg.Wait()
	fmt.Println("compactCount", compactCount)
}

func ConcurrV2(concurrency int) {
	wg := new(sync.WaitGroup)

	loop := 100000
	numG := loopO / loop
	wg.Add(numG)
	fmt.Println(runtime.GOMAXPROCS(concurrency), runtime.NumCPU())
	for i := 0; i < numG; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < loop; i++ {
				ret = 10.34 / 2.34
				compactCount++
			}
		}()
	}

	wg.Wait()
	fmt.Println("compactCount", compactCount)
}

func TestSchedule(t *testing.T) {
	runtime.GOMAXPROCS(1)

	for i := 0; i < 100; i++ {
		go func() {
			for {

			}
		}()
	}
	time.Sleep(3 * time.Second)

	// can run， 抢占式调度: G运行时间超过forcePreemptNS(10ms)，则抢占这个P
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
		}
	}()

	time.Sleep(1 * time.Second)
}
