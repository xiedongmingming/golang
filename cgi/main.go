package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cgi"
	"os"
)

func main() {

	fmt.Println(os.Args[1])

	switch os.Args[1] {
	case "1":
		main1()
	case "2":
		main2()
	}

}

func main1() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		handler := new(cgi.Handler)

		handler.Path = "/home/yejianfeng/go/gopath/src/cgi-script/" + r.URL.Path
		log.Println(handler.Path)
		handler.Dir = "/home/yejianfeng/go/gopath/src/cgi-script/"

		handler.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(":8989", nil))
}

func main2() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		handler := new(cgi.Handler)

		handler.Path = "F:/liteide/go/bin/go"
		handler.Dir = "F:/liteide/workspace/src/base/cgi/cgi-script"

		script := "F:/liteide/workspace/src/base/cgi/cgi-script" + r.URL.Path

		args := []string{"run", script}

		handler.Args = append(handler.Args, args...)
		handler.Env = append(handler.Env, "GOPATH=F:/liteide/workspace")
		handler.Env = append(handler.Env, "GOROOT=F:/liteide/go")

		handler.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(":8989", nil))
}
