/**
 * @Author: Robby
 * @File name: main.go
 * @Create date: 2021-10-20
 * @Function:
 **/

package main

import "sync"

type ReentrantLock struct {
	lock      *sync.Mutex
	cond      *sync.Cond
	recursion int32
	host      int64
}

func main() {

}
