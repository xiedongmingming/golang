package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

type Person struct {
	Name    string
	Age     int
	Emails  []string
	Company string
	Role    string
}
type OnlineUser struct {
	User      []*Person
	LoginTime string
}

func Handler(w http.ResponseWriter, r *http.Request) {

	dumx := Person{
		Name:    "zoro",
		Age:     27,
		Emails:  []string{"dg@gmail.com", "dk@hotmail.com"},
		Company: "Omron",
		Role:    "SE",
	}

	chxd := Person{
		Name:   "chxd",
		Age:    27,
		Emails: []string{"test@gmail.com", "d@hotmail.com"},
	}

	onlineUser := OnlineUser{
		User: []*Person{&dumx, &chxd},
	}

	// t := template.New("Person template")
	// t, err := t.Parse(templ)

	t, err := template.ParseFiles("tmpl.html")

	checkError(err)

	err = t.Execute(w, onlineUser)

	checkError(err)
}

func checkError(err error) {

	if err != nil {

		fmt.Println("fatal error ", err.Error())

		os.Exit(1)

	}
}

func main() { //利用GOLANG的TEMPLATE模板包进行WEB开发

	http.HandleFunc("/", Handler)

	http.ListenAndServe(":8888", nil) //开始监听

}
