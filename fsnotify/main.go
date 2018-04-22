package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/howeyc/fsnotify"
)

func usage() {

	fmt.Println("")
	fmt.Printf("usage: %s montior-directory file-max-bytes-limit\n", os.Args[0])
	fmt.Println("for example:")
	fmt.Printf("%s /opt/logs 1024\n", os.Args[0])

	os.Exit(1)

}

func isDir(dirname string) bool {

	fhandler, err := os.Stat(dirname)

	if !(err == nil || os.IsExist(err)) {
		return false
	} else {
		return fhandler.IsDir()
	}
}

func isFile(filename string) bool {

	fhandler, err := os.Stat(filename)

	if !(err == nil || os.IsExist(err)) { // 表示文件不存在的情况
		return false
	} else if fhandler.IsDir() {
		return false
	}

	return true
}

func emptiedFile(filename string) bool { // 清空文件

	FN, err := os.Create(filename) // 新建文件时若文件已经存在则内容会被截断
	defer FN.Close()

	if err != nil {
		return false
	}

	fmt.Fprint(FN, "")

	return true
}
func getFileByteSize(filename string) (bool, int64) { // 获取文件内容的字节数

	if !isFile(filename) {
		return false, 0
	}

	fhandler, _ := os.Stat(filename)

	return true, fhandler.Size()
}

func main() {

	var maxByte int64 = 1024 * 1024

	if len(os.Args) < 2 {
		usage()
	}

	if len(os.Args) >= 3 {

		maxByte_, err := strconv.Atoi(os.Args[2])

		if err != nil {

			log.SetPrefix("[error] ")
			log.Println(os.Args[2], "is not a legitimate int number")

			usage()
		}

		maxByte = int64(maxByte_)
	}

	dirpath := os.Args[1] // 被监控的目录

	if !isDir(dirpath) {

		log.SetPrefix("[error] ")
		log.Println(dirpath, "is not a legitimate directory")

		usage()
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	go func() {

		for {

			select {

			case ev := <-watcher.Event:

				if ev.IsModify() { // 事件是修改类型

					_, size := getFileByteSize(ev.Name) // 参数表示文件名

					log.Println("event: ", ev, ", byte:", size)

					if size >= maxByte {

						if !emptiedFile(ev.Name) {

							log.SetPrefix("[error] ")
							log.Printf("%s : can not empty file\n", ev.Name)

						}
					}
				}

			case err := <-watcher.Error:

				log.Println("error:", err)

			}

		}

	}()

	err = watcher.Watch(dirpath) // 表示监控指定的目录
	if err != nil {
		log.Fatal(err)
	}

	<-done

	watcher.Close()
}

func follow(filename string) error {

	file, _ := os.Open(filename)

	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	_ = watcher.Watch(filename)

	r := bufio.NewReader(file)

	for {

		by, err := r.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}

		fmt.Print(string(by))
		if err != io.EOF {
			continue
		}

		if err = waitForChange(watcher); err != nil {
			return err
		}
	}

}

func waitForChange(w *fsnotify.Watcher) error {

	for {

		select {
		case event := <-w.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				return nil
			}
		case err := <-w.Errors:
			return err
		}
	}
}
