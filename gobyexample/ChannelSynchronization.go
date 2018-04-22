package main

import "fmt"
import "time"

func worker(done chan bool) {

	fmt.Print("working...")

	time.Sleep(time.Second)

	fmt.Println("done")

	done <- true
}

func main() { // 用于工作同步

	done := make(chan bool, 1)

	go worker(done)

	<-done
}
