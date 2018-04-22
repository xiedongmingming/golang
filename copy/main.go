package main

import (
	"fmt"
)

type Address [20]byte

func main() {

	var etherbase Address

	etherbase[0] = 1

	base := etherbase // 内存拷贝

	fmt.Println(etherbase[0])
	fmt.Println(base[0])

	etherbase[0] = 2

	fmt.Println(etherbase[0])
	fmt.Println(base[0])

}
