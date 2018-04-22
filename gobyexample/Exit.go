package main

import "fmt"
import "os"

func main() { //当通过GORUN运行时会显示状态码(运行二进制时需要通过$?)

	defer fmt.Println("!") // 不会被执行

	os.Exit(3)

}
