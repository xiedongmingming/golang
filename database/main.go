package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	_ "reflect"
	_ "os"
	"strings"
	"bytes"
	"os"
	"flag"
	"os/exec"
	"fugetech.com/antifraud/log"
)

const (
	CONST_DB_MAX_OPEN_CONNS int = 100
	CONST_DB_MAX_IDLE_CONNS int = 10
)

type TableMetaData struct {
	Fileds []*FiledMetaData
}
type FiledMetaData struct {
	Filed   string
	Type    string
	Null    string
	Key     string
	Default interface{}
	Extra   string
}

var header = `package model

import (
	"time"
	"fmt"
)

// ORM映射中的MODEL对象
`

func CamelCase(name string) string {

	var buffer bytes.Buffer

	strs := strings.Split(name, "_")

	for _, str := range strs {
		buffer.WriteString(strings.ToUpper(str[0:1]) + str[1:])
	}

	return buffer.String()
}

func main() {

	var dst = flag.String("dst", "", "目标文件名")

	flag.Parse()

	if *dst == "" {

		println("params error")

		flag.PrintDefaults()

		os.Exit(1)
	}

	modelFile, err := os.OpenFile(*dst, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0600)

	if err != nil {
		fmt.Println("打开文件失败: " + err.Error())
		return
	}

	modelFile.WriteString(header)

	dsn := "root:xieming4243054@tcp(180.76.55.12:3306)/adwetec?charset=utf8"

	db, err := sql.Open("mysql", dsn) // 注意:执行OPEN的时候并不会去真正连接数据库

	if err != nil {
		fmt.Println("打开数据库失败: " + err.Error())
		return
	}

	if e := db.Ping(); e != nil {
		fmt.Println("PING数据库失败: " + dsn)
		return
	}

	db.SetMaxOpenConns(CONST_DB_MAX_OPEN_CONNS)
	db.SetMaxIdleConns(CONST_DB_MAX_IDLE_CONNS)

	rows, err := db.Query("show tables;")

	if err != nil {
		fmt.Println("显示数据库表错误: " + err.Error())
		return
	}

	tables := make([]string, 0)

	for rows.Next() {

		var table string

		err := rows.Scan(&table)

		if err != nil {
			fmt.Println("扫描数据库错误: " + err.Error())
			return
		}

		fmt.Println(fmt.Sprintf("%s", table))

		tables = append(tables, table)

	}

	for _, table := range tables {

		if table != "report_feed_adgroup" {
			continue
		}

		tablemd := &TableMetaData{}

		tablemd.Fileds = make([]*FiledMetaData, 0)

		rows, err := db.Query(fmt.Sprintf("desc %s;", table))

		if err != nil {
			fmt.Println("DESC数据库表错误: " + err.Error())
			return
		}

		for rows.Next() {

			var filed, types, null, defaults, extra, key interface{}

			err := rows.Scan(&filed, &types, &null, &key, &defaults, &extra)

			if err != nil {
				fmt.Println("扫描数据库错误: " + err.Error())
				return
			}

			fmt.Println(fmt.Sprintf("%v %v %v %v %v %v", string(filed.([]byte)), string(types.([]byte)), string(null.([]byte)), string(key.([]byte)), defaults, string(extra.([]byte))))

			rowmd := &FiledMetaData{
				Filed:   string(filed.([]byte)),
				Type:    string(types.([]byte)),
				Null:    string(null.([]byte)),
				Key:     string(key.([]byte)),
				Default: defaults,
				Extra:   string(extra.([]byte)),
			}

			tablemd.Fileds = append(tablemd.Fileds, rowmd)

		}

		modelFile.WriteString("// **********************************************************************\n")
		modelFile.WriteString(fmt.Sprintf("type %s struct {\n", CamelCase(table)))

		var jsonStruct bytes.Buffer
		var tojsonFunc bytes.Buffer

		jsonStruct.WriteString(fmt.Sprintf("type %s struct {\n", CamelCase(table)+"Json"))

		tojsonFunc.WriteString(fmt.Sprintf("func (this *%s) ToJson() *%s {\n", CamelCase(table), CamelCase(table)+"Json"))
		tojsonFunc.WriteString(fmt.Sprintf("\n\treturn &%s{\n", CamelCase(table)+"Json"))

		max := 0

		for _, row := range tablemd.Fileds {
			if len(row.Filed) > max {
				max = len(row.Filed)
			}
		}

		for _, row := range tablemd.Fileds {

			name := CamelCase(row.Filed)

			switch row.Type {

			case "bigint(100) unsigned", "bigint(100)", "int(20)":

				if row.Extra != "auto_increment" {

					modelFile.WriteString(fmt.Sprintf("\t%s\tint64\t`xorm:\"%s\"`\n", name, row.Filed))

					jsonStruct.WriteString(fmt.Sprintf("\t%s\tint64\t`json:\"%s\"`\n", name, row.Filed))

					tojsonFunc.WriteString(fmt.Sprintf("\t\t%s:\tthis.%s,\n", name, name))

				} else {

					modelFile.WriteString(fmt.Sprintf("\t%s\tint64\t`xorm:\"%s\"`%s\n", name, row.Filed, "// XORM自动自增长"))

					jsonStruct.WriteString(fmt.Sprintf("\t%s\tint64\t`json:\"%s\"`\n", name, row.Filed))

					tojsonFunc.WriteString(fmt.Sprintf("\t\t%s:\tthis.%s,\n", name, name))

				}

			case "varchar(255)", "varchar(1000)":

				modelFile.WriteString(fmt.Sprintf("\t%s\tstring\t`xorm:\"%s\"`\n", name, row.Filed))

				jsonStruct.WriteString(fmt.Sprintf("\t%s\tstring\t`json:\"%s\"`\n", name, row.Filed))

				tojsonFunc.WriteString(fmt.Sprintf("\t\t%s:\tthis.%s,\n", name, name))

			case "int(10)":

				modelFile.WriteString(fmt.Sprintf("\t%s\tint32\t`xorm:\"%s\"`\n", name, row.Filed))

				jsonStruct.WriteString(fmt.Sprintf("\t%s\tint32\t`json:\"%s\"`\n", name, row.Filed))

				tojsonFunc.WriteString(fmt.Sprintf("\t\t%s:\tthis.%s,\n", name, name))

			case "double(10,2)":

				modelFile.WriteString(fmt.Sprintf("\t%s\tfloat64\t`xorm:\"%s\"`\n", name, row.Filed))

				jsonStruct.WriteString(fmt.Sprintf("\t%s\tfloat64\t`json:\"%s\"`\n", name, row.Filed))

				tojsonFunc.WriteString(fmt.Sprintf("\t\t%s:\tthis.%s,\n", name, name))

			case "datetime":

				modelFile.WriteString(fmt.Sprintf("\t%s\ttime.Time\t`xorm:\"%s\"`\n", name, row.Filed))

				jsonStruct.WriteString(fmt.Sprintf("\t%s\tstring\t`json:\"%s\"`\n", name, row.Filed))

				tojsonFunc.WriteString(fmt.Sprintf("\t\t%s:\tthis.%s.String(),\n", name, name))

			case "tinyint(2)":

				modelFile.WriteString(fmt.Sprintf("\t%s\tint\t`xorm:\"%s\"`\n", name, row.Filed))

				jsonStruct.WriteString(fmt.Sprintf("\t%s\tint\t`json:\"%s\"`\n", name, row.Filed))

				tojsonFunc.WriteString(fmt.Sprintf("\t\t%s:\tthis.%s,\n", name, name))

			case "float(10,2)":

				modelFile.WriteString(fmt.Sprintf("\t%s\tfloat32\t`xorm:\"%s\"`\n", name, row.Filed))

				jsonStruct.WriteString(fmt.Sprintf("\t%s\tfloat32\t`json:\"%s\"`\n", name, row.Filed))

				tojsonFunc.WriteString(fmt.Sprintf("\t\t%s:\tthis.%s,\n", name, name))

			default:

				log.Infof("不符合规定的数据库表字段类型: %s --> %s --> %s", table, row.Filed, row.Type)

			}

		}

		switch table {// 在这里给JSON添加额外的字段
		case "account_mapping":
			jsonStruct.WriteString("\n") // 中间用空行分割
			jsonStruct.WriteString("\tAdvertiserName\tstring\t`json:\"advertiser_name\"`\n")
			jsonStruct.WriteString("\tMediaName\tstring\t`json:\"media_name\"`\n")
			jsonStruct.WriteString("\tBalance\tfloat64\t`json:\"balance\"`\n")
			jsonStruct.WriteString("\tBudget\tfloat64\t`json:\"budget\"`\n")
		case "report_feed_account":
			jsonStruct.WriteString("\n") // 中间用空行分割
			jsonStruct.WriteString("\tBalance\tfloat64\t`json:\"balance\"`\n")
			jsonStruct.WriteString("\tBudget\tfloat64\t`json:\"budget\"`\n")
		}

		modelFile.WriteString("}\n")
		jsonStruct.WriteString("}\n")
		tojsonFunc.WriteString("\t}\n}\n")

		modelFile.WriteString(jsonStruct.String())
		modelFile.WriteString(tojsonFunc.String())

	}

	cmd := exec.Command("gofmt", "model.go") // 从TRACKING服务器上拉取日志存放到TMP目录下

	out, err := cmd.Output()

	if err != nil {
		fmt.Println("gofmt file error: " + err.Error())
		return
	}

	fmt.Println(string(out))

	if modelFile != nil {
		modelFile.Close()
	}

	modelFile, err = os.OpenFile(*dst, os.O_WRONLY|os.O_TRUNC, 0600)

	if err != nil {
		fmt.Println("open file error: " + err.Error())
		return
	}

	modelFile.WriteString(string(out))

	if modelFile != nil {
		modelFile.Close()
	}

}
