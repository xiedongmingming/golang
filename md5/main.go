package main

import (
	"fmt"
	"regexp"
	"strings"

	"crypto/sha1"
	"crypto/md5"

	"encoding/hex"
	"os"
	"io"
)

var md5Regexp = regexp.MustCompile(`[a-f0-9]{32}`)

func isMD5(str string) bool {
	return md5Regexp.MatchString(str)
}

// 如果是md5，小写返回
// 如果不是，md5返回
func handleMobileID(id string) string {
	if isMD5(id) {
		return strings.ToLower(id)
	}
	sum := md5.Sum([]byte(id))
	return hex.EncodeToString(sum[:])
}
func HashFile() { // 对文件内容计算SHA1

	TestFile := "123.txt"

	infile, inerr := os.Open(TestFile)

	if inerr == nil {

		md5h := md5.New()
		io.Copy(md5h, infile)
		fmt.Printf("%x %s\n", md5h.Sum([]byte("")), TestFile)

		sha1h := sha1.New()
		io.Copy(sha1h, infile)
		fmt.Printf("%x %s\n", sha1h.Sum([]byte("")), TestFile)

	} else {
		fmt.Println(inerr)
		os.Exit(1)
	}
}
func main() {

	fmt.Println(handleMobileID("861258036731851"))

	// 哈希函数
	TestString := "hi, pandaman!"

	Md5Inst := md5.New()
	Md5Inst.Write([]byte(TestString))
	Result := Md5Inst.Sum([]byte(""))
	fmt.Printf("%x\n\n", Result)

	Sha1Inst := sha1.New()
	Sha1Inst.Write([]byte(TestString))
	Result = Sha1Inst.Sum([]byte(""))
	fmt.Printf("%x\n\n", Result)
}
