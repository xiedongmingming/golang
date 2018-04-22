package main

import (
	"os"
	"net"
	"log"
	"time"
	"net/rpc"
	"net/http"
	"base/rpc/common"

	"net/rpc/jsonrpc"

	"github.com/streadway/amqp"
)



func JsonRpcServer()  {

	lf, _ := os.OpenFile("/var/log/"+time.Now().Format("2006-01-02T15:04:05")+".log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

	defer lf.Close()

	log.SetOutput(lf)

	mpmsg := new(common.Mpmsg)

	rpc.Register(mpmsg)

	addr, _ := net.ResolveTCPAddr("tcp", "localhost:135")

	ln, e := net.ListenTCP("tcp", addr)

	if e != nil {
		panic(e)
	}

	for {

		conn, e := ln.Accept()

		if e != nil {
			continue
		}

		go jsonrpc.ServeConn(conn)
	}
}

func main() {

	arith := new(common.Arith)

	rpc.Register(arith) // 注册服务
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", "localhost:1234")

	if e != nil {
		log.Fatal("listen error:", e)
	}

	http.Serve(l, nil)
}
