package main

import (
	"context"
	"fmt"
	"net/http"
)

const requestIDKey = "rid" // 键

func newContextWithRequestID(ctx context.Context, req *http.Request) context.Context {

	reqID := req.Header.Get("X-Request-ID")

	if reqID == "" {
		reqID = "0"
	}

	return context.WithValue(ctx, requestIDKey, reqID) // 创建一个可以储存 K-V 的 context
}

func requestIDFromContext(ctx context.Context) string { // 获取键对应的值
	return ctx.Value(requestIDKey).(string)
}

func middleWare(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		ctx := newContextWithRequestID(req.Context(), req)

		next.ServeHTTP(w, req.WithContext(ctx))

	})
}

func h(w http.ResponseWriter, req *http.Request) { // 处理函数

	reqID := requestIDFromContext(req.Context()) // 在GOLANG1.7中"net/http"原生支持将Context嵌入到*http.Request中

	fmt.Fprintln(w, "request id: ", reqID)

	return
}

func main() {

	http.Handle("/", middleWare(http.HandlerFunc(h)))

	http.ListenAndServe(":9201", nil)
}
