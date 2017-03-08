package main

import (
	"fmt"
	"math/rand"
	"time"
)

//*********************************************************************************
//如果一个源码文件被声明为属于MAIN代码包且该文件代码中包含无参数声明和结果声明的MAIN函数则它就是命
//令源码文件.可以通过GORUN命令直接启动运行
//*********************************************************************************

type _httpWorker struct { //结构体加下划线
	cli  *http.Client
	clis *http.Client
}

func main() {

	fmt.Println("Hello World")

	var months = [...]string{ //数组的表示法
		"January",
		"February",
		"March",
		"April",
		"May",
		"June",
		"July",
		"August",
		"September",
		"October",
		"November",
		"December",
	}

	rand.Seed(time.Now().Unix()) //初始化随机种子

	reader := bufio.NewReader(os.Stdin)
	data, _, e := reader.ReadLine()
	if e != nil {
		break
	}
}
