package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"time"
)

const (
	SecretKey = "welcome to wangshubo's blog"
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

func StartServer() {

	http.HandleFunc("/login", LoginHandler)

	http.Handle("/resource", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(ProtectedHandler)),
	))

	log.Println("now listening...")

	http.ListenAndServe(":8080", nil)
}

func main() {


	StartServer()

}

func JsonResponse(response interface{}, w http.ResponseWriter) {

	json, err := json.Marshal(response)

	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {

	response := Response{"gained access to protected resource"}

	JsonResponse(response, w)

}

func LoginHandler(w http.ResponseWriter, r *http.Request) { // 登录请求处理

	var user UserCredentials

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {

		w.WriteHeader(http.StatusForbidden)

		fmt.Fprint(w, "error in request")

		return
	}

	if strings.ToLower(user.Username) != "someone" || user.Password != "p@ssword" {

		w.WriteHeader(http.StatusForbidden)

		fmt.Println("error logging in")

		fmt.Fprint(w, "invalid credentials")

		return
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := make(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Minute * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()

	token.Claims = claims

	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)

		fmt.Fprintln(w, "error extracting the key")

		fatal(err)
	}

	tokenString, err := token.SignedString([]byte(SecretKey))

	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)

		fmt.Fprintln(w, "error while signing the token")

		fatal(err)
	}

	response := Token{tokenString}

	JsonResponse(response, w)

}

func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) { // 进行TOKEN的验证

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

	if err == nil {

		if token.Valid {

			fmt.Println("TOKEN验证通过 ... ")

			next(w, r) // 调用下一个处理器

		} else {

			w.WriteHeader(http.StatusUnauthorized)

			fmt.Fprint(w, "token is not valid")

		}

	} else {

		w.WriteHeader(http.StatusUnauthorized)

		fmt.Fprint(w, "unauthorized access to this resource")

	}

}
