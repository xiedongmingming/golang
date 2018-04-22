package main

import (
	"fmt"
	"reflect"
	"strconv"
)

func main() {

	var chs1 = make(chan int)
	var chs2 = make(chan float64)
	var chs3 = make(chan string)

	var ch4close = make(chan int)

	defer close(ch4close)

	go func(c chan int, ch4close chan int) {

		for i := 0; i < 5; i++ {
			c <- i
		}

		close(c)

		ch4close <- 1

	}(chs1, ch4close)

	go func(c chan float64, ch4close chan int) {

		for i := 0; i < 5; i++ {
			c <- float64(i) + 0.1
		}

		close(c)

		ch4close <- 1

	}(chs2, ch4close)

	go func(c chan string, ch4close chan int) {

		for i := 0; i < 5; i++ {
			c <- "string:" + strconv.Itoa(i)
		}

		close(c)

		ch4close <- 1

	}(chs3, ch4close)

	var selectCase = make([]reflect.SelectCase, 4)

	selectCase[0].Dir = reflect.SelectRecv
	selectCase[0].Chan = reflect.ValueOf(chs1)

	selectCase[1].Dir = reflect.SelectRecv
	selectCase[1].Chan = reflect.ValueOf(chs2)

	selectCase[2].Dir = reflect.SelectRecv
	selectCase[2].Chan = reflect.ValueOf(chs3)

	selectCase[3].Dir = reflect.SelectRecv
	selectCase[3].Chan = reflect.ValueOf(ch4close)

	done := 0
	finished := 0

	for finished < len(selectCase)-1 {

		chosen, recv, recvOk := reflect.Select(selectCase) // 执行选择

		if recvOk {
			done = done + 1
			switch chosen { // 表示选中的是那个通道
			case 0:
				fmt.Println(chosen, recv.Int()) // 选择的值
			case 1:
				fmt.Println(chosen, recv.Float())
			case 2:
				fmt.Println(chosen, recv.String())
			case 3:
				finished = finished + 1
				done = done - 1
				fmt.Println("finished\t", finished)
			}
		}
	}

	fmt.Println("done", done)

}
