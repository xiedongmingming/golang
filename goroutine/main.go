package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {

	name := "Eric"

	go func() {
		fmt.Printf("Hello, %s.\n", name)
	}()

	go println("go")
	go func() {
		println("go func")
	}()

	time.Sleep(1)

	runtime.Gosched() // 让其它协程有机会运行
}
