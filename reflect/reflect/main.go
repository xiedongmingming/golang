package main

import (
	"fmt"
	"reflect"
)

type Service struct {
	name string
	id   uint
}

func main() {

	var service *Service

	fmt.Println(reflect.ValueOf(&service).String())
	fmt.Println(reflect.ValueOf(&service).Elem().String())

	element := reflect.ValueOf(&service).Elem() // 首先返回的是指针类型之后返回的是指针指向的内容

	fmt.Println(element.Type().String())

	//	if running, ok := n.services[element.Type()]; ok {

	//		element.Set(reflect.ValueOf(running)) // 填充值

	//		return nil
	//	}
}
