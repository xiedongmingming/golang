package main

import (
	"fmt"
	"reflect"
	"io"
)

// *******************************************************************
type MyReader struct {
	Name string
}

func (r MyReader) Read(p []byte) (n int, err error) {
	return 0, nil
}

// *******************************************************************
type Bird struct {
	Name           string
	LifeExpectance int
}

func (b *Bird) Fly() {
	fmt.Println("i am flying")
}

// *******************************************************************
func main() {

	// 对所有接口进行反射都可以得到一个包含TYPE和VALUE的信息结构
	// TYPE:  被反射的这个变量本身的类型信息
	// VALUE: 该变量实例本身的信息

	// *******************************************************
	fmt.Println("=========================================")

	var reader io.Reader

	reader = &MyReader{"a.txt"}

	fmt.Println("类型名称: ", reflect.TypeOf(reader))  // 返回的是参数的动态类型 *main.MyReader
	fmt.Println("值型字符: ", reflect.ValueOf(reader)) // 返回新结构体(已初始化) &{a.txt}

	fmt.Println("=========================================")

	// 1. 获取类型信息
	var x float64 = 3.4

	fmt.Println("类型: ", reflect.TypeOf(x))
	fmt.Println("类型: ", reflect.TypeOf(x).Kind())

	fmt.Println("=========================================")

	sparrow := &Bird{"sparrow", 3}

	s := reflect.ValueOf(sparrow).Elem()

	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {

		f := s.Field(i)

		fmt.Printf("%d: %s %s = %v\n", i, typeOfT.Field(i).Name, f.Type(), f.Interface())

	}

}
