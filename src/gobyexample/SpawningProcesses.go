package main

import "fmt"
import "io/ioutil"
import "os/exec"

func main() { //通过GO程序启动外部进程

	dateCmd := exec.Command("date") // 对外部进程的包装

	dateOut, err := dateCmd.Output()
	if err != nil {
		panic(err)
	}

	fmt.Println("> date")
	fmt.Println(string(dateOut))

	grepCmd := exec.Command("grep", "hello") //第二个参数表示第一个参数命令对应的参数

	grepIn, _ := grepCmd.StdinPipe()   // 命令的输入
	grepOut, _ := grepCmd.StdoutPipe() // 命令的输出--还有错误
	grepCmd.Start()                    // 启动命令
	grepIn.Write([]byte("hello grep\ngoodbye grep"))
	grepIn.Close() // 关闭输入
	grepBytes, _ := ioutil.ReadAll(grepOut)
	grepCmd.Wait() //等到进程结束

	fmt.Println("> grep hello")
	fmt.Println(string(grepBytes))

	lsCmd := exec.Command("bash", "-c", "ls -a -l -h") // 通过这种方式来将命令和参数放在一起
	lsOut, err := lsCmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println("> ls -a -l -h")
	fmt.Println(string(lsOut))
}
