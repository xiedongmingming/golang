package main

import (
	"fmt"
	"time"
	"github.com/BurntSushi/toml"
	"encoding/json"
	"bytes"
	"sync"
	"path/filepath"
)

// 常用配置文件有:
// JSON INI YAML TOML
// 其中这些文件对应的GOLANG处理库如下:
//
// encoding/json	-- 标准库中的包(可以处理JSON配置文件缺点是不能加注释)
// gcfg 			-- 处理INI配置文件
// toml				-- 处理TOML配置文件
// viper			-- 处理JSON\TOML\YAML\HCL以及JAVA属性配置文件

// 参考地址:
// https://github.com/toml-lang/toml

type TomlConfig struct {
	Title   string
	Owner   OwnerInfo
	DB      Database `toml:"database"`
	Servers map[string]Server
	Clients Clients
}

func (this *TomlConfig) String() string {

	b, err := json.Marshal(*this)

	if err != nil {
		return fmt.Sprintf("%+v", *this)
	}

	var out bytes.Buffer

	err = json.Indent(&out, b, "", "    ")

	if err != nil {
		return fmt.Sprintf("%+v", *this)
	}

	return out.String()
}

type OwnerInfo struct {
	Name string
	Org  string `toml:"organization"`
	Bio  string
	DOB  time.Time
}
type Database struct {
	Server  string
	Ports   []int
	ConnMax int `toml:"connection_max"`
	Enabled bool
}
type Server struct {
	IP string
	DC string
}
type Clients struct {
	Data  [][]interface{}
	Hosts []string
}

func main() {

	var config TomlConfig

	filePath := "./conf.toml"

	if _, err := toml.DecodeFile(filePath, &config); err != nil {
		panic(err)
	}

	fmt.Println(config.String())

}

var (
	cfg  *TomlConfig
	once sync.Once
)

func Config() *TomlConfig { // 配置的单例模式

	once.Do(func() {

		filePath, err := filepath.Abs("./conf.toml")

		if err != nil {
			panic(err)
		}

		fmt.Printf("parse toml file once. filepath: %s\n", filePath)

		if _, err := toml.DecodeFile(filePath, &cfg); err != nil {
			panic(err)
		}

	})

	return cfg
}
