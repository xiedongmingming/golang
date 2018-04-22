package main

import (
	"os"
	"os/signal"
	"syscall"

	fugeLog "fugetech.com/anti-fraud/logrus"
)

var logger = fugeLog.NewLogger("info", "console", "./log")

func main() {

	Signaler()
}

func Signaler() { // 捕获信号量

	ch := make(chan os.Signal, 10)

	signal.Notify(ch,
		syscall.SIGABRT,   // Signal(0x6)
		syscall.SIGALRM,   // Signal(0xe)
		syscall.SIGBUS,    // Signal(0x7)
		syscall.SIGCHLD,   // Signal(0x11)
		syscall.SIGCLD,    // Signal(0x11)
		syscall.SIGCONT,   // Signal(0x12)
		syscall.SIGFPE,    // Signal(0x8)
		syscall.SIGHUP,    // Signal(0x1)
		syscall.SIGILL,    // Signal(0x4)
		syscall.SIGINT,    // Signal(0x2)
		syscall.SIGIO,     // Signal(0x1d)
		syscall.SIGIOT,    // Signal(0x6)
		syscall.SIGKILL,   // Signal(0x9)
		syscall.SIGPIPE,   // Signal(0xd)
		syscall.SIGPOLL,   // Signal(0x1d)
		syscall.SIGPROF,   // Signal(0x1b)
		syscall.SIGPWR,    // Signal(0x1e)
		syscall.SIGQUIT,   // Signal(0x3)
		syscall.SIGSEGV,   // Signal(0xb)
		syscall.SIGSTKFLT, // Signal(0x10)
		syscall.SIGSTOP,   // Signal(0x13)
		syscall.SIGSYS,    // Signal(0x1f)
		syscall.SIGTERM,   // Signal(0xf)
		syscall.SIGTRAP,   // Signal(0x5)
		syscall.SIGTSTP,   // Signal(0x14)
		syscall.SIGTTIN,   // Signal(0x15)
		syscall.SIGTTOU,   // Signal(0x16)
		syscall.SIGUNUSED, // Signal(0x1f)
		syscall.SIGURG,    // Signal(0x17)
		syscall.SIGUSR1,   // Signal(0xa)
		syscall.SIGUSR2,   // Signal(0xc)
		syscall.SIGVTALRM, // Signal(0x1a)
		syscall.SIGWINCH,  // Signal(0x1c)
		syscall.SIGXCPU,   // Signal(0x18)
		syscall.SIGXFSZ,   // Signal(0x19)
	)

	for {

		switch sig := <-ch; sig {

		case syscall.SIGINT: // 重新加载配置文件: ctrl+c kill -2 pid

			logger.Info("syscall.SIGINT")

			signal.Stop(ch)

			return

		case syscall.SIGKILL: // 这个信号无法捕捉

			logger.Info("syscall.SIGKILL")

			signal.Stop(ch)

			return

		case syscall.SIGUSR1: // 重新加载配置文件: kill -10 pid

			logger.Info("syscall.SIGUSR1")

			signal.Stop(ch)

			return

		case syscall.SIGTERM: // 重新加载配置文件: kill pid

			logger.Info("syscall.SIGTERM")

			signal.Stop(ch)

			return

		default:

			logger.Info("default: " + sig.String())

		}
	}
}
