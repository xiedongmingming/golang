package main

import (
	_ "fmt"
	_ "io/ioutil"
	_ "os/exec"
)

func main() {

	select {}

	//  ********************************************************
	//	date := exec.Command("date")
	//	output, err := date.Output()
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println("> date")
	//	fmt.Println(string(output))
	//  ********************************************************

	//  ********************************************************
	//	grepcmd := exec.Command("grep", "hello") // 要运行的命令

	//	// 这里我们明确的获取输入/输出管道
	//	in, _ := grepcmd.StdinPipe()
	//	out, _ := grepcmd.StdoutPipe() // StderrPipe

	//	grepcmd.Start() // 启动程序

	//	in.Write([]byte("hello grep\ngoodbye grep\ndddd gggg\nhello sss"))

	//	in.Close()

	//	grepbytes, _ := ioutil.ReadAll(out)
	//	grepcmd.Wait() // 等待程序运行结束(我们忽略了错误检测)
	//	fmt.Println("> grep hello")

	//	fmt.Println(string(grepbytes))
	//  ********************************************************

	// 注意:当我们需要提供一个明确的命令和参数数组来生成命令(和能够只需要提供一行命令行字符串相比)你想使用通过一个字符串生成一个完整的命令那么你可以使用BASH命令的-C选项
	//	lsCmd := exec.Command("bash", "-c", "ls -a -l -h")
	//	lsOut, err := lsCmd.Output()
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println("> ls -a -l -h")
	//	fmt.Println(string(lsOut))

}
