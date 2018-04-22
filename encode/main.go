package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {

	sText := "中文"
	textQuoted := strconv.QuoteToASCII(sText) // "\u4e2d\u6587"

	fmt.Println(textQuoted)

	textUnquoted := textQuoted[1 : len(textQuoted)-1] // \u4e2d\u6587

	fmt.Println(textUnquoted)

	sUnicodev := strings.Split(textUnquoted, "\\u")

	var context string

	for _, v := range sUnicodev {

		if len(v) < 1 { // 第一个不要
			continue
		}

		temp, err := strconv.ParseInt(v, 16, 32)

		if err != nil {
			panic(err)
		}

		context += fmt.Sprintf("%c", temp)

	}

	fmt.Println(context)

}
