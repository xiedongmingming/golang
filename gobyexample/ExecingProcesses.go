package main

import "syscall"
import "os"
import "os/exec"

func main() { //当前进程会被启动的进程替换

	binary, lookErr := exec.LookPath("ls") //获取命令的绝对路径

	if lookErr != nil {
		panic(lookErr)
	}

	args := []string{"ls", "-a", "-l", "-h"}

	env := os.Environ()

	execErr := syscall.Exec(binary, args, env) // 需要环境
	if execErr != nil {
		panic(execErr)
	}
}
