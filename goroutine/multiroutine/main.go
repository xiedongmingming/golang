package main

import (
	"bufio"

	"flag"
	"fmt"
	"io"

	"math/rand"

	"os"
	"path/filepath"

	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	fugeLog "fugetech.com/anti-fraud/logrus"

	"github.com/BurntSushi/toml"
)

// ***************************************************************************
// 配置参数
type Params struct {
	Source       string `toml:"source"`       // 源文件
	Destination  string `toml:"destination"`  // 存放报告的目录
	ChannelSize  int    `toml:"channelsize"`  //
	GoroutineNum int    `toml:"goroutinenum"` //
	LogFile      string `toml:"logfile"`      //
}

func NewParams() *Params {
	return &Params{ // 在此提供默认值
		Source:       "",
		Destination:  "./report",
		ChannelSize:  10000000,
		GoroutineNum: 11,
		LogFile:      "./log",
	}
}
func (this *Params) InitFromFile(path string) {
	if _, err := toml.DecodeFile(path, this); err != nil {
		panic(fmt.Sprintf("can't decode configure file: %s", path))
	}
}

// ***************************************************************************
type Service struct {
	Params         *Params
	Logger         *fugeLog.Logger
	DestFile       *os.File //结果数据位置
	Buses          []chan string
	Group          *sync.WaitGroup
	PerSecMutex    *sync.Mutex
	LinesPerSecond int64
}

func NewService() *Service {
	return &Service{
		PerSecMutex: new(sync.Mutex),
		Group:       new(sync.WaitGroup),
	}
}
func (this *Service) DealLine(ch chan string) {

	for {

		line := <-ch

		line = strings.TrimSpace(line)

		fields := strings.Split(line, " ")

		if len(fields) == 2 {

			score, _ := strconv.Atoi(fields[1])

			this.DestFile.WriteString(fmt.Sprintf("%2d %s\n", score, fields[0]))

		} else {
			this.DestFile.WriteString(fmt.Sprintf("没有分数 %s\n", line))
		}

		this.LinesPerSecond++

		this.Group.Done()

	}
}

// ***************************************************************************
// 程序入口
func main() {

	var service = NewService()

	service.Params = NewParams()

	cfg := flag.String("conf", "", "specify configuration path") // 用于搜集命令行参数

	flag.Parse() // 解析命令行参数

	if len(*cfg) != 0 {
		service.Params.InitFromFile(*cfg) // 从配置文件中获取启动参数
	}

	service.Logger = fugeLog.NewLogger("info", "file", service.Params.LogFile)

	rand.Seed(time.Now().Unix())

	runtime.GOMAXPROCS(runtime.NumCPU()) // 利用CPU多核来处理HTTP请求

	if _, err := os.Stat(filepath.Dir(service.Params.Destination)); os.IsNotExist(err) { // 创建报告文件位置

		if err := os.MkdirAll(filepath.Dir(service.Params.Destination), 0700); err != nil {

			service.Logger.Info(fmt.Sprintf("mkdirall dir %s error %s", service.Params.Destination, err.Error()))

			os.Exit(0)
		}
	}

	f, err := os.OpenFile(fmt.Sprintf(service.Params.Destination+"file"), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0600)

	if err != nil {

		service.Logger.Info("open file error: " + err.Error())

		os.Exit(0)
	}

	service.DestFile = f // 保存记录的文件(根据分数进行分档)

	go func() {

		for {

			<-time.Tick(time.Second)

			service.Logger.Info(fmt.Sprintf("lines per second: %d gorountine: %d", service.LinesPerSecond, runtime.NumGoroutine()))

			service.PerSecMutex.Lock()
			service.LinesPerSecond = 0
			service.PerSecMutex.Unlock()
		}
	}()

	service.Buses = make([]chan string, service.Params.GoroutineNum)

	for i := 0; i < service.Params.GoroutineNum; i++ {
		service.Buses[i] = make(chan string, service.Params.ChannelSize) // 通道(缓冲文件中的每一行)
	}

	for i := 0; i < service.Params.GoroutineNum; i++ { // 启动多个协程
		go service.DealLine(service.Buses[i])
	}

	// ***********************************************************************
	// 处理源文件
	sourceFile, err := os.Open(service.Params.Source)

	if err != nil {

		service.Logger.Info(fmt.Sprintf("open source file error: %s %s", service.Params.Source, err.Error()))

		os.Exit(0)
	}

	defer func() {

		if sourceFile != nil {
			sourceFile.Close()
		}

	}()

	rd := bufio.NewReader(sourceFile)

	for { // 逐行读取日志文件

		line, err := rd.ReadString('\n')

		if err != nil {
			if err == io.EOF {
				service.Logger.Info("文件处理完毕: " + service.Params.Source)
				break
			} else {
				service.Logger.Info("文本行读取错误: " + err.Error())
			}
			continue
		}

		service.Group.Add(1)

		fields := strings.Split(line, " ")

		if len(fields) == 2 {

			score, _ := strconv.Atoi(fields[1])

			service.Buses[score/10] <- line

		} else {
			service.Buses[service.Params.GoroutineNum-1] <- line
		}

	}
	// ***********************************************************************

	service.Group.Wait() // 等待所有的行都被处理完毕

	service.DestFile.Close()

	service.Logger.Info("处理结束")

}
