package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"

	fugeLog "fugetech.com/anti-fraud/logrus"
)

func PrivateIpExp() {

	// 10.0.0.0~10.255.255.255（A类）
	// 172.16.0.0~172.31.255.255（B类）
	// 192.168.0.0~192.168.255.255（C类）

	reg, _ := regexp.Compile(`(^10(\.(2[0-4]\d|25[0-5]|[01]?\d\d?)){3}$)|(^172\.(1[6-9]|2\d|3[0-1]){1}(\.(2[0-4]\d|25[0-5]|[01]?\d\d?)){2}$)|(^192\.168(\.(2[0-4]\d|25[0-5]|[01]?\d\d?)){2}$)`)

	fmt.Println(reg.MatchString("10.0.0.0"))
	fmt.Println(reg.MatchString("10.255.255.255"))
	fmt.Println(reg.MatchString("10.0.0.256"))
	fmt.Println(reg.MatchString("172.16.0.0"))
	fmt.Println(reg.MatchString("172.31.255.255"))
	fmt.Println(reg.MatchString("172.18.255.255"))
	fmt.Println(reg.MatchString("172.18.255.256"))
	fmt.Println(reg.MatchString("192.168.0.0"))
	fmt.Println(reg.MatchString("192.168.255.255"))
	fmt.Println(reg.MatchString("192.168.0.256"))

	if ok, _ := regexp.MatchString(`^172\.(1[6-9]|2\d|3[0-1]){1}(\.(2[0-4]\d|25[0-5]|[01]?\d\d?)){2}$`, "172.31.0.0"); ok {
		fmt.Println("匹配")
	} else {
		fmt.Println("不匹配")
	}

	for i := 0; i < 999; i++ {
		if ok, _ := regexp.MatchString(`^(2[0-4]\d|25[0-5]|[01]?\d\d?)$`, fmt.Sprintf("%d", i)); ok {
			fmt.Print("匹配")
		} else {
			fmt.Print("不匹配")
		}
	}

	if ok, _ := regexp.MatchString(`^10(\.(2[0-4]\d|25[0-5]|[01]?\d\d?)){3}$`, "10.10.255.255"); ok {
		fmt.Println("匹配")
	} else {
		fmt.Println("不匹配")
	}
}
func IpExp(line string) (string, bool) { // 查看指定的参数是不是合法的IP地址

	ip := ipreg.FindString(line)

	return ip, ipreg.MatchString(ip)
}
func IdfaExp(line string) bool {
	return false
}

var logger = fugeLog.NewLogger("info", "file", "./log")

var ipreg = regexp.MustCompile(`((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)_`) // IP地址正则表达式
var idfa = regexp.MustCompile(`\w{8}-\w{4}-\w{4}-\w{4}-\w{12}`)

func main() {

	var alllines, deallines int

	cfg := flag.String("s", "", "source file") // 用于搜集命令行参数
	cnt := flag.Int("n", 100, "goroutine num") // 用于搜集命令行参数

	flag.Parse() // 解析命令行参数

	buses := make(chan string, 100000)

	var group sync.WaitGroup

	// ************************************************************************
	for i := 0; i < *cnt; i++ { // 启动多个线程来处理行

		go func() {

			for {

				line := <-buses

				line = strings.TrimSpace(line)

				ip, isMatched := IpExp(line)

				if !isMatched || !strings.Contains(ip, "10") {
					logger.Info(fmt.Sprintf("%s", line))
				}

				deallines++

				group.Done()
			}
		}()
	}

	// ************************************************************************
	// 处理源文件
	source, err := os.Open(*cfg)

	if err != nil {

		logger.Info(fmt.Sprintf("open source file error: %s %s", *cfg, err.Error()))

		os.Exit(0)
	}

	defer func() {
		if source != nil {
			source.Close()
		}
	}()

	rd := bufio.NewReader(source)

	for { // 逐行读取日志文件

		line, err := rd.ReadString('\n')

		if err != nil {
			if err == io.EOF {
				logger.Info("文件处理完毕: " + *cfg)
				break
			} else {
				logger.Info("文本行读取错误: " + err.Error())
			}
			continue
		}

		buses <- line

		alllines++

		group.Add(1)
	}

	defer func() {
		logger.Info(fmt.Sprintf("alllines: %s deallines: %s", alllines, deallines))
	}()

	// ************************************************************************

	group.Wait() // 等待所有的行被处理完成

	logger.Info("处理完成")

}
