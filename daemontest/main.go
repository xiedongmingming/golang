package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"syscall"
	"time"

	fugeLog "fugetech.com/anti-fraud/logrus"
)

// *****************************************************************************
// 进程后台运行的几种方法 - nohup/setsid/&
// 如何让命令提交后不受本地关闭终端窗口/网络断开连接的干扰呢?
// 场景:如果只是临时有一个命令需要长时间运行什么方法能最简便的保证它在后台稳定运行呢?
// 当用户注销或者网络断开时终端会收到HUP信号从而关闭其所有子进程.因此我们的解决办法就有两种途径:要么让进程忽略HUP信号要么让进程运行在新的会话里从而成为不属于此终端的子进程
// 1.nohup
// 标准输出和标准错误缺省会被重定向到NOHUP.OUT文件中.一般我们可在结尾加上"&"来将命令同时放入后台运行也可用">filename 2>&1"来更改缺省的重定向文件名
// 2.setsid
// NOHUP无疑能通过忽略HUP信号来使我们的进程避免中途被中断.但如果我们换个角度思考如果我们的进程不属于接受HUP信号的终端的子进程那么自然也就不会受到HUP信号的影响了.
// setsid ping www.ibm.com
// 值得注意的是上例中我们的进程ID(PID)为31094.而它的父ID(PPID)为1(即为INIT进程ID)并不是当前终端的进程ID
// 3.&
// 这里还有一个关于SUBSHELL的小技巧.我们知道将一个或多个命名包含在"()"中就能让这些命令在子SHELL中运行中从而扩展出很多有趣的功能
// 当我们将"&"也放入"()"内之后我们就会发现所提交的作业并不在作业列表中也就是说是无法通过JOBS来查看的.让我们来看看为什么这样就能躲过HUP信号的影响吧.
// (ping www.ibm.com &)
// 新提交的进程的父ID(PPID)为1(即为INIT进程ID)并不是当前终端的进程ID.因此并不属于当前终端的子进程从而也就不会受到当前终端的HUP信号的影响了.
// *****************************************************************************
func daemon(nochdir, noclose int) int { // 实现调用程序的后台运行功能

	var ret, ret2 uintptr

	var err syscall.Errno

	darwin := runtime.GOOS == "darwin"

	if syscall.Getppid() == 1 { // 表示父进程为系统初始化进程

		// logger.Info("already a daemon")

		return 0

	}

	ret, ret2, err = syscall.RawSyscall(syscall.SYS_FORK, 0, 0, 0) // fork off the parent process

	if err != 0 {
		return -1
	}

	if ret2 < 0 { // failure
		os.Exit(-1)
	}

	if darwin && ret2 == 1 { // handle exception for darwin
		ret = 0
	}

	if ret > 0 { // if we got a good pid, then we call exit the parent process.
		os.Exit(0)
	}

	_ = syscall.Umask(0) // change the file mode mask

	// 之后的父进程PPID就变成了INIT进程的ID(1)
	s_ret, s_errno := syscall.Setsid() // 为当前进程设置一个新的会话ID(这样当前进程就不会收到所在终端进程(父进程)HUP信号的影响了)

	if s_errno != nil {

		log.Printf("error: syscall.setsid errno: %d", s_errno)

	}

	if s_ret < 0 {
		return -1
	}

	if nochdir == 0 { // 切换当前工作目录
		os.Chdir("/")
	}

	if noclose == 0 { // 标准文件描述符重定向

		f, e := os.OpenFile("/dev/null", os.O_RDWR, 0)

		if e == nil {

			fd := f.Fd()

			syscall.Dup2(int(fd), int(os.Stdin.Fd()))  // 表示进行标准输入的重定向
			syscall.Dup2(int(fd), int(os.Stdout.Fd())) // 表示进行标准输出的重定向
			syscall.Dup2(int(fd), int(os.Stderr.Fd())) // 表示进行标准错误的重定向

		}
	}

	return 0
}

func main() {

	var logger = fugeLog.NewLogger("info", "console", "./log")

	// daemon(0, 1)

	if syscall.Getppid() == 1 {

		logger.Info("already a daemon")

	} else {

		logger.Info(fmt.Sprintf("ppid: %d", syscall.Getppid()))

	}

	for {

		logger.Info("hello")

		time.Sleep(1 * time.Second)

	}

}
