package main

import "os"
import "fmt"

func main() { //只是简单获取命令行后面的参数(原样获取)

	argsWithProg := os.Args // 第一个元素为程序名称
	argsWithoutProg := os.Args[1:]

	arg := os.Args[3]

	fmt.Println(os.Args[0]) // 相对于执行该命令的路径
	fmt.Println(argsWithProg)
	fmt.Println(argsWithoutProg)
	fmt.Println(arg)
}
