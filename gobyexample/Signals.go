package main

import "fmt"
import "os"
import "os/signal"
import "syscall"

func main() {

	sigs := make(chan os.Signal, 1) //信号通道(只有一个值)

	done := make(chan bool, 1) //布尔值通道(只有一个值)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM) // 表示程序关注的信号--UNIX信号

	go func() {

		sig := <-sigs

		fmt.Println()

		fmt.Println(sig)
		fmt.Println(sig.String())

		done <- true
	}()

	fmt.Println("awaiting signal")
	<-done //表示从通道中读取值
	fmt.Println("exiting")
}
