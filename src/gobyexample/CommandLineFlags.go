package main

import "flag"
import "fmt"

func main() { //用于解析命令行选项

	//在参数前面使用单横杠和双横杠没有区别

	wordPtr := flag.String("word", "foo", "a string")
	numbPtr := flag.Int("numb", 42, "an int")
	boolPtr := flag.Bool("fork", false, "a bool")

	var svar string
	flag.StringVar(&svar, "svar", "bar", "a string var")

	flag.Parse()

	fmt.Println("word:", *wordPtr)
	fmt.Println("numb:", *numbPtr)
	fmt.Println("fork:", *boolPtr)
	fmt.Println("svar:", svar)
	fmt.Println("tail:", flag.Args()) // -word=opt a1 a2 a3 结尾位置参数(必须放在最后)

	// 位置参数:标志必须出现在位置参数之前
	// 默认提供帮助信息: -h --help
}
