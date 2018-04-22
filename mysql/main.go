package main

import (
	"fmt"
	"log"
	"net/http"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// GOLANG的MYSQL连接池实现:
// GOLANG内部自带了连接池功能.SQL.OPEN函数实际上是返回一个连接池对象而不是单个连接
// 在OPEN的时候并没有去连接数据库只有在执行QUERY、EXEC方法的时候才会去实际连接数据库
// 在一个应用中同样的库连接只需要保存一个SQL.OPEN之后的DB对象就可以了不需要多次OPEN

// 数据库连接池测试

// 1.开启WEB服务
func startHttpServer() {

	http.HandleFunc("/pool", pool)

	err := http.ListenAndServe(":9090", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// 2.DB对象初始化
var db *sql.DB

func init() {

	db, _ = sql.Open("mysql", "rdswetec:FWkKYmmLtEbQzs#kG!vA4f!eYiClhuuu@tcp(192.168.0.18:3306)/adwetec_prod?charset=utf8")

	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)

	db.Ping()

	// 连接池的实现关键在于设置最大连接数和设置最大空闲连接数其中:
	//
	// SetMaxOpenConns: 用于设置最大打开的连接数(默认值为0表示不限制)
	// SetMaxIdleConns: 用于设置闲置的连接数
	//
	// 设置最大的连接数可以避免并发太高导致连接MYSQL出现"TOO MANY CONNECTIONS"的错误
	// 设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用
}

// 3.请求方法
func pool(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query("SELECT * FROM adwetec_user LIMIT 1")
	defer rows.Close()

	checkErr(err)

	columns, _ := rows.Columns()

	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]string)

	for rows.Next() {

		err = rows.Scan(scanArgs...)

		for i, col := range values {

			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}

		}
	}

	fmt.Println(record)

	fmt.Fprintln(w, "finish")
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

// 使用AB进行并发测试访问:
//
// $ ab -c 100 -n 1000 'http://localhost:9090/pool'
// 在数据库中查看连接进程: show processlist;
// 可以看到有100来个进程
//
// 因为避免了重复创建连接所以使用连接池可以很明显的提高性能
func main() {
	startHttpServer()
}

// 小结
// GOLANG这边实现的连接池只提供了SetMaxOpenConns和SetMaxIdleConns方法进行连接池方面的配置
// 在使用的过程中有一个问题就是数据库本身对连接有一个超时时间的设置.如果超时时间到了数据库会单方面断掉连接
// 此时再用连接池内的连接进行访问就会出错
// packets.go:32: unexpected EOF
// packets.go:118: write tcp 192.168.3.90:3306: broken pipe
//
// 上面的错误都是go-sql-drive本身的输出
// 有的时候还会出现bad connection的错误
// 多请求几次后连接池会重新打开新连接这时候就没有问题了
