package main

import (
	"fmt"
	"log"
	"net/rpc"
	"net/http"
	"net/rpc/jsonrpc"
	"base/rpc/common"

	"adwetec.com/tools/utils"
)

func JsonRpcClient()  {

	rpcClient, e := jsonrpc.Dial("tcp", "localhost:135")

	if e != nil {

		log.Println("error dail rpc server:", e)

		http.Error(w, e.Error(), http.StatusInternalServerError)

		return
	}

	var reply map[string]interface{}

	e = rpcClient.Call("Mpmsg.Handle", msg, &reply)

	if e != nil {

		log.Println("error call rpc method:", e)

		http.Error(w, e.Error(), http.StatusInternalServerError)

		return
	}
}

func main() {

	// 无论是调用RPC客户端的同步或者是异步方法都必须指定要调用的服务及其方法名称以及一个客户端传入参数的引用还有一个用于接收处理结果参数的指针

	client, err := rpc.DialHTTP("tcp", "localhost"+":1234")

	if err != nil {
		log.Fatal("dialing:", err)
	}

	args := &common.Args{199, 8}

	var reply int

	// 同步
	err = client.Call("Arith.Multiply", args, &reply)

	if err != nil {
		log.Fatal("arith error:", err)
	}

	fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)

	// 异步:
	quotient := new(common.Quotient)

	divCall := client.Go("Arith.Divide", args, &quotient, nil)

	replyCall := <-divCall.Done

	if replyCall.Error != nil {
		fmt.Println("RPC调用错误")
	} else {
		fmt.Println(utils.String(replyCall.Reply))
	}
}
