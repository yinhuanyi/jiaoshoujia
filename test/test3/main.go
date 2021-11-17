/**
 * @Author: Robby
 * @File name: main.go
 * @Create date: 2021-10-20
 * @Function:
 **/

package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/cch123/goroutineid"
)

type ReentrantLock struct {
	lock      *sync.Mutex
	cond      *sync.Cond
	recursion int32 // 计数器
	id        int64
}

func NewReentrantLock() sync.Locker {
	res := &ReentrantLock{
		lock:      new(sync.Mutex),
		recursion: 0,
		id:        0,
	}
	res.cond = sync.NewCond(res.lock)
	return res
}

// Lock 可重入锁加锁
func (rt *ReentrantLock) Lock() {
	id := goroutineid.GetGoID()
	fmt.Println(id)

	rt.lock.Lock()

	defer rt.lock.Unlock()

	// 如果当前goroutine拿到这把锁，那么直接recursion+1就好
	if rt.id == id {
		rt.recursion++
		return
	}

	// 如果是其他的goroutine拿到这把锁，那么需要阻塞等待
	for rt.recursion != 0 {
		fmt.Println("阻塞等待")
		rt.cond.Wait()
	}

	// 让可重入锁设置为其他的goroutine的id
	rt.id = id
	rt.recursion = 1
}

// Unlock 可重入锁解锁
func (rt *ReentrantLock) Unlock() {
	rt.lock.Lock()
	defer rt.lock.Unlock()

	// goroutine解锁的时候，如果发现recursion为0，或者不是自己的锁，直接return，什么都不做，或者panic也可以
	if rt.recursion == 0 || rt.id != goroutineid.GetGoID() {
		panic(fmt.Sprintf("the wrong call host: (%d); current_id: %d; recursion: %d", rt.id, goroutineid.GetGoID(), rt.recursion))
		//fmt.Printf("the wrong call host: (%d); current_id: %d; recursion: %d\n", rt.host, goroutineid.GetGoID(), rt.recursion)
		//return
	}

	rt.recursion--

	// 如果发现recursion为0了，说明锁已经可以被释放了
	if rt.recursion == 0 {
		rt.cond.Signal() // 通知
	}

}

func main() {
	lock := NewReentrantLock()

	for i := 0; i < 500; i++ {
		go func() {
			lock.Lock()
		}()

	}

	time.Sleep(time.Second)
	for i := 0; i < 500; i++ {
		go func() {
			lock.Unlock()
		}()

	}

}
