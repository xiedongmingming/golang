package main

import (
	"fmt"
	"sync"
	"time"
)

// 锁:
// 1、互斥锁: sync.Mutex
// 零值表示了未被锁定的互斥量(开箱即用锁)
// 当锁定一个已被锁定的互斥锁时那么进行重复锁定操作的协程将会被阻塞
// 当解锁一个已处于解锁状态的互斥锁进行解锁操作的时候就会引发一个运行时恐慌
// 2、读写锁: sync.RWMutex
// 在读写锁管辖的范围内它允许任意个读操作同时进行.但是在同一时刻它只允许有一个写操作在进行
// 并且在某一个写操作被进行的过程中读操作的进行也是不被允许的
// 零值就已经是立即可用的读写锁了
// 只有当所有的读锁被解锁之后才会去唤醒试图去写锁定的协程
// 对一个未被读锁定的读写锁进行读锁不会引发运行时恐防

func main() {

	var mutex sync.Mutex

	fmt.Println("lock the lock. (g0)")
	mutex.Lock()
	fmt.Println("the lock is locked. (g0)")

	for i := 1; i <= 3; i++ {
		go func(i int) {
			fmt.Printf("lock the lock. (g%d)\n", i)
			mutex.Lock()
			fmt.Printf("the lock is locked. (g%d)\n", i)
		}(i)
	}

	time.Sleep(time.Second)
	fmt.Println("unlock the lock. (g0)")
	mutex.Unlock()
	fmt.Println("the lock is unlocked. (g0)")

	time.Sleep(time.Second)

	var rwmutex sync.RWMutex

	rwmutex.Lock()
	rwmutex.UUnlock()

	rwmutex.RLock()
	rwmutex.RUnlock()

	locker := rwmutex.RLocker() // 返回一个实现了指定接口的值: sync.Locker
	locker.Lock()
	locker.Unlock()
}
