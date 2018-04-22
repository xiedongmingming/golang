package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {

	//	os.Args
	//	os.Chdir()
	//	os.Chmod()
	//	os.Chown()
	//	os.Chtimes()
	//	os.Clearenv()
	//	os.Create()
	//	os.DevNull
	//	os.Environ()
	//	os.ErrExist
	//	os.ErrInvalid
	//	os.ErrInvalid
	//	os.ErrPermission
	//	os.Exit()
	//	os.Expand()
	//	os.ExpandEnv()
	//	os.File
	//	os.FileInfo
	//	os.FileMode
	//	os.FileMode
	//	os.FindProcess()
	//	os.Getegid()
	//	os.Getenv()
	//	os.Geteuid()
	//	os.Getgid()
	//	os.Getgroups()
	//	os.Getpagesize()
	//	os.Getpid()
	//	os.Getppid()
	//	os.Getuid()
	//	os.Getwd()
	//	os.Hostname()
	//	os.Interrupt
	//	os.IsExist()
	//	os.IsNotExist()
	//	os.IsPathSeparator()
	//	os.IsPermission()
	//	os.Kill
	//	os.Lchown()
	//	os.Link()
	//	os.LinkError
	//	os.LookupEnv()
	//	os.Lstat()
	//	os.Mkdir()
	//	os.MkdirAll()
	//	os.ModeAppend
	//	os.NewFile()
	//	os.NewSyscallError()
	//	os.Open()
	//	os.OpenFile()
	//	os.Pipe()
	//	os.Readlink()
	//	os.Remove()
	//	os.RemoveAll()
	//	os.Rename()
	//	os.SameFile()
	//	os.Setenv()
	//	os.StartProcess()
	//	os.Stat()
	//	os.Symlink()
	//	os.TempDir()
	//	os.Truncate()
	//	os.Unsetenv()

	// os.TempDir --　默认用于临时文件的目录

	fmt.Println(os.TempDir()) // 获取当前系统的临时文件路径

	//**************************************************************************
	// ioutil.Discard

	cnt, err := io.WriteString(ioutil.Discard, "wokaoa")
	if err != nil {
		return
	}

	fmt.Println(cnt) // 输出6

	//**************************************************************************
	// ioutil.NopCloser()
	// func NopCloser(r io.Reader) io.ReadCloser 该接口仅提供CLOSE方法

	//**************************************************************************
	// ioutil.ReadAll()
	// func ReadAll(r io.Reader) ([]byte, error)

	var dir string
	var data []byte

	dir, err = filepath.Abs(filepath.Dir(os.Args[0])) // 获取程序运行的目录路径
	if err != nil {
		return
	}

	data, err = ioutil.ReadFile(dir + "/main.go")
	if err != nil {
		return
	}
	reader := bytes.NewReader(data)

	data, err = ioutil.ReadAll(reader)

	// fmt.Printf("%s", string(data))

	//**************************************************************************
	// ioutil.ReadDir()
	// func ReadDir(dirname string) ([]os.FileInfo, error) // 返回一个有序的、子目录信息的列表

	//**************************************************************************
	// ioutil.ReadFile()
	// func ReadFile(filename string) ([]byte, error)

	dir, err = filepath.Abs(filepath.Dir(os.Args[0])) // 获取程序运行的目录路径
	if err != nil {
		return
	}

	data, err = ioutil.ReadFile(dir + "/main.go")
	if err != nil {
		return
	}

	// fmt.Printf("%s", string(data))

	//**************************************************************************
	// ioutil.TempDir()
	// func TempDir(dir, prefix string) (name string, err error) // 线程安全

	//**************************************************************************
	// ioutil.TempFile()
	// func TempFile(dir, prefix string) (f *os.File, err error) // 线程安全--调用本函数的程序有责任在不需要临时文件时摧毁它

	//**************************************************************************
	// ioutil.WriteFile()
	// func WriteFile(filename string, data []byte, perm os.FileMode) error
}
