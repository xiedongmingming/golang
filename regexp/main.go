package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func Test000() {
	re, _ := regexp.Compile(`a=(\d+),b=(\d+)`)
	c := re.ReplaceAllString("test regexp a=1234,b=5678. test regexp replace a=8765,b=3210 ", "c=$2,d=$1")

	fmt.Println(c)

	//这个测试一个字符串是否符合一个表达式。
	match, _ := regexp.MatchString("p([a-z]+)ch", "peach")
	fmt.Println(match)
	//上面我们是直接使用字符串，但是对于一些其他的正则任务，你需要使用 Compile 一个优化的 Regexp 结构体。
	r, _ := regexp.Compile("p([a-z]+)ch")
	//这个结构体有很多方法。这里是类似我们前面看到的一个匹配测试。
	fmt.Println(r.MatchString("peach"))
	//这是查找匹配字符串的。
	fmt.Println(r.FindString("peach punch"))
	//这个也是查找第一次匹配的字符串的，但是返回的匹配开始和结束位置索引，而不是匹配的内容。
	fmt.Println(r.FindStringIndex("peach punch"))
	//Submatch 返回完全匹配和局部匹配的字符串。例如，这里会返回 p([a-z]+)ch 和 `([a-z]+) 的信息。
	fmt.Println(r.FindStringSubmatch("peach punch"))
	//类似的，这个会返回完全匹配和局部匹配的索引位置。
	fmt.Println(r.FindStringSubmatchIndex("peach punch"))
	//带 All 的这个函数返回所有的匹配项，而不仅仅是首次匹配项。例如查找匹配表达式的所有项。
	fmt.Println(r.FindAllString("peach punch pinch", -1))
	//All 同样可以对应到上面的所有函数。
	fmt.Println(r.FindAllStringSubmatchIndex(
		"peach punch pinch", -1))
	//这个函数提供一个正整数来限制匹配次数。
	fmt.Println(r.FindAllString("peach punch pinch", 2))
	//上面的例子中，我们使用了字符串作为参数，并使用了如 MatchString 这样的方法。我们也可以提供 []byte参数并将 String 从函数命中去掉。
	fmt.Println(r.Match([]byte("peach")))
	//创建正则表示式常量时，可以使用 Compile 的变体MustCompile 。因为 Compile 返回两个值，不能用语常量。
	r = regexp.MustCompile("p([a-z]+)ch")
	fmt.Println(r)
	//regexp 包也可以用来替换部分字符串为其他值。
	fmt.Println(r.ReplaceAllString("a peach", "<fruit>"))
	//Func 变量允许传递匹配内容到一个给定的函数中，
	in := []byte("a peach")
	out := r.ReplaceAllFunc(in, bytes.ToUpper)
	fmt.Println(string(out))

	text := "Hello 世界！123 go."

	reg := regexp.MustCompile(`[a-z]+`)             // 查找连续的小写字母
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["ello" "o"]

	reg = regexp.MustCompile(`[^a-z]+`)             // 查找连续的非小写字母
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["H" " 世界！123 G" "."]

	reg = regexp.MustCompile(`[\w]+`)               // 查找连续的单词字母
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["Hello" "123" "Go"]

	reg = regexp.MustCompile(`[^\w\s]+`)            // 查找连续的非单词字母、非空白字符
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["世界！" "."]

	reg = regexp.MustCompile(`[[:upper:]]+`)        // 查找连续的大写字母
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["H" "G"]

	reg = regexp.MustCompile(`[[:^ascii:]]+`)       // 查找连续的非 ASCII 字符
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["世界！"]

	reg = regexp.MustCompile(`[\pP]+`)              // 查找连续的标点符号
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["！" "."]

	reg = regexp.MustCompile(`[\PP]+`)              // 查找连续的非标点符号字符
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["Hello 世界" "123 Go"]

	reg = regexp.MustCompile(`[\p{Han}]+`)          // 查找连续的汉字
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["世界"]

	reg = regexp.MustCompile(`[\P{Han}]+`)          // 查找连续的非汉字字符
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["Hello " "！123 Go."]

	reg = regexp.MustCompile(`Hello|Go`)            // 查找 Hello 或 Go
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["Hello" "Go"]

	reg = regexp.MustCompile(`^H.*\s`)              // 查找行首以 H 开头，以空格结尾的字符串
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["Hello 世界！123 "]

	reg = regexp.MustCompile(`(?U)^H.*\s`)          // 查找行首以 H 开头，以空白结尾的字符串（非贪婪模式）
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["Hello "]

	reg = regexp.MustCompile(`(?i:^hello).*Go`)     // 查找以 hello 开头（忽略大小写），以 Go 结尾的字符串
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["Hello 世界！123 Go"]

	reg = regexp.MustCompile(`\QGo.\E`)             // 查找 Go.
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["Go."]

	reg = regexp.MustCompile(`(?U)^.* `)            // 查找从行首开始，以空格结尾的字符串（非贪婪模式）
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["Hello "]

	reg = regexp.MustCompile(` [^ ]*$`)             // 查找以空格开头，到行尾结束，中间不包含空格字符串
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // [" Go."]

	reg = regexp.MustCompile(`(?U)\b.+\b`)          // 查找“单词边界”之间的字符串
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["Hello" " 世界！" "123" " " "Go"]

	reg = regexp.MustCompile(`[^ ]{1,4}o`)          // 查找连续 1 次到 4 次的非空格字符，并以 o 结尾的字符串
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["Hello" "Go"]

	reg = regexp.MustCompile(`(?:Hell|G)o`)         // 查找 Hello 或 Go
	fmt.Printf("%q\n", reg.FindAllString(text, -1)) // ["Hello" "Go"]

	reg = regexp.MustCompile(`(?PHell|G)o`)                   // 查找 Hello 或 Go，替换为 Hellooo、Gooo
	fmt.Printf("%q\n", reg.ReplaceAllString(text, "${n}ooo")) // "Hellooo 世界！123 Gooo."

	reg = regexp.MustCompile(`(Hello)(.*)(Go)`)              // 交换 Hello 和 Go
	fmt.Printf("%q\n", reg.ReplaceAllString(text, "$3$2$1")) // "Go 世界！123 Hello."

	// 特殊字符的查找
	reg = regexp.MustCompile("[\\f\\t\\n\\r\\v\\123\\x7F\\x{10FFFF}\\\\\\^\\$\\.\\*\\+\\?\\{\\}\\(\\)\\[\\]\\|]")
	fmt.Printf("%q\n", reg.ReplaceAllString("\f\t\n\r\v\123\x7F\U0010FFFF\\^$.*+?{}()[]|", "-"))

}

//为了方便提取，我们会把正则表达式中要提取的数据使用命名方式来书写正则表达式。这个技术在Go语言中如何实现，可以看下面这篇博客：

//Using the Go Regexp Package
//http://blog.kamilkisiel.net/blog/2012/07/05/using-the-go-regexp-package/

//简单期间，这里复制其中几个例子的代码：

//我们期望在字符串  1000abcd123  中找出前后两个数字。

//例子1：匹配到这个字符串的例子
//package main

//import(
//    "fmt"
//    "regexp"
//)

//var digitsRegexp = regexp.MustCompile(`(\d+)\D+(\d+)`)

//func main(){
//    someString:="1000abcd123"
//    fmt.Println(digitsRegexp.FindStringSubmatch(someString))
//}
//上面代码输出：
//[1000abcd123 1000 123]

//例子2：使用带命名的正则表达式
//package main

//import(
//    "fmt"
//    "regexp"
//)

//var myExp=regexp.MustCompile(`(?P<first>\d+)\.(\d+).(?P<second>\d+)`)

//func main(){
//    fmt.Printf("%+v",myExp.FindStringSubmatch("1234.5678.9"))
//}

//上面代码输出，所有匹配到的都输出了：

//[1234.5678.9 1234 5678 9]

//这里的Named capturing groups  (?P<name>) 方式命名正则表达式是 python、Go语言特有的， java、c# 是 (?<name>) 命名方式。

//参考： http://www.crifan.com/detailed_explanation_about_python_regular_express_match_named_group/

//http://osdir.com/ml/go-language-discuss/2012-07/msg00605.html

//Go 支持的正则表达式语法可以参看：https://code.google.com/p/re2/wiki/Syntax

//例子3：对正则表达式类扩展一个获得所有命名信息的方法，并使用它。
//package main

//import(
//    "fmt"
//    "regexp"
//)

////embed regexp.Regexp in a new type so we can extend it
//type myRegexp struct{
//    *regexp.Regexp
//}

////add a new method to our new regular expression type
//func(r *myRegexp)FindStringSubmatchMap(s string) map[string]string{
//    captures:=make(map[string]string)

//    match:=r.FindStringSubmatch(s)
//    if match==nil{
//        return captures
//    }

//    for i,name:=range r.SubexpNames(){
//        //Ignore the whole regexp match and unnamed groups
//        if i==0||name==""{
//            continue
//        }

//        captures[name]=match[i]

//    }
//    return captures
//}

////an example regular expression
//var myExp=myRegexp{regexp.MustCompile(`(?P<first>\d+)\.(\d+).(?P<second>\d+)`)}

//func main(){
//    mmap:=myExp.FindStringSubmatchMap("1234.5678.9")
//    ww:=mmap["first"]
//    fmt.Println(mmap)
//    fmt.Println(ww)
//}
//上面代码的输出结果：
//map[first:1234 second:9]

//1234

//例子4，抓取限号信息，并记录到一个Map中。

//package main

//import(
//    "fmt"
//    iconv "github.com/djimenez/iconv-go"
//    "io/ioutil"
//    "net/http"
//    "os"
//    "regexp"
//)

//// embed regexp.Regexp in a new type so we can extend it
//type myRegexp struct{
//    *regexp.Regexp
//}

//// add a new method to our new regular expression type
//func(r *myRegexp)FindStringSubmatchMap(s string)[](map[string]string){
//    captures:=make([](map[string]string),0)

//    matches:=r.FindAllStringSubmatch(s,-1)

//    if matches==nil{
//        return captures
//    }

//    names:=r.SubexpNames()

//    for _,match:=range matches{

//        cmap:=make(map[string]string)

//        for pos,val:=range match{
//            name:=names[pos]
//            if name==""{
//                continue
//            }

//            /*
//                fmt.Println("+++++++++")
//                fmt.Println(name)
//                fmt.Println(val)
//            */
//            cmap[name]=val
//        }

//        captures=append(captures,cmap)

//    }

//    return captures
//}

//// 抓取限号信息的正则表达式
//var myExp=myRegexp{regexp.MustCompile(`自(?P<byear>[\d]{4})年(?P<bmonth>[\d]{1,2})月(?P<bday>[\d]{1,2})日至(?P<eyear>[\d]{4})年(?P<emonth>[\d]{1,2})月(?P<eday>[\d]{1,2})日，星期一至星期五限行机动车车牌尾号分别为：(?P<n11>[\d])和(?P<n12>[\d])、(?P<n21>[\d])和(?P<n22>[\d])、(?P<n31>[\d])和(?P<n32>[\d])、(?P<n41>[\d])和(?P<n42>[\d])、(?P<n51>[\d])和(?P<n52>[\d])`)}

//func ErrorAndExit(err error){
//    fmt.Fprintln(os.Stderr,err)
//    os.Exit(1)
//}

//func main(){
//    response,err:=http.Get("http://www.bjjtgl.gov.cn/zhuanti/10weihao/index.html")
//    defer response.Body.Close()

//    if err!=nil{
//        ErrorAndExit(err)
//    }

//    input,err:=ioutil.ReadAll(response.Body)
//    if err!=nil{
//        ErrorAndExit(err)
//    }
//    body :=make([]byte,len(input))
//    iconv.Convert(input,body,"gb2312","utf-8")

//    mmap:=myExp.FindStringSubmatchMap(string(body))

//    fmt.Println(mmap)
//}

//上述代码输出：

//[map[n32:0 n22:9 emonth:7 n11:3 n41:1 n21:4 n52:7 bmonth:4 n51:2 bday:9 n42:6 byear:2012 eday:7 eyear:2012 n12:8 n31:5]
//map[emonth:10 n41:5 n52:6 n31:4 byear:2012 n51:1 eyear:2012 n32:9 bmonth:7 n22:8 bday:8 n11:2 eday:6 n42:0 n21:3 n12:7]
//map[bday:7 n51:5 n22:7 n31:3 eday:5 n32:8 byear:2012 bmonth:10 emonth:1 eyear:2013 n11:1 n12:6 n52:0 n21:2 n42:9 n41:4]
//map[eyear:2013 byear:2013 n22:6 eday:10 bmonth:1 n41:3 n32:7 n31:2 n21:1 n11:5 bday:6 n12:0 n51:4 n42:8 emonth:4 n52:9]]

//这段代码首先下载北京市交管局的网页；然后把这个gb2312的页面转换成utf-8编码，然后用正则表达式提取其中的限号信息。

//func main() {

//	//	var digitsRegexp = regexp.MustCompile(`\w{8}-\w{4}-\w{4}-\w{4}-\w{12}`)

//	//	someString := "time=\"2017-07-23T00:35:23+08:00\" level=info msg=\"[[SendGetHttp: client do error Get https://api.rtbasia.com/ipscore/query?key=T1FUGETECH4212111176&did=000CB534-65FC-4120-8CC0-3199630C4E75&r=1&fm=7: dial tcp 175.6.228.154:443: getsockopt: connection refused]]\" "

//	//	fmt.Println(digitsRegexp.FindStringSubmatch(someString))
//	//	fmt.Println(digitsRegexp.FindString(someString))

//	var advertiserIDRegex = regexp.MustCompile(`^(dmp-\d{1,10}|\d{1,10})$`)
//	var campaignIDRegex = regexp.MustCompile(`^\d{0,10}$`)
//	var orderIDRegex = regexp.MustCompile(`^\d{0,10}$`)

//	if !advertiserIDRegex.MatchString(advertiserID) {
//		return false
//	}
//}

var idfaExp = regexp.MustCompile(`\w{8}-\w{4}-\w{4}-\w{4}-\w{12}`) // IDFA的正则表达式
var imeiExp = regexp.MustCompile(`\d{15}`)                         // IMEI的正则表达式

type Result struct {
	Err  int      `json: "error"`
	Data []string `json: "data"`
}

func main() {

	data := "{\"errno\":0,\"data\":[\"d9545e34339afd76b6f669a45239f35f\"]}"

	var result Result

	err := json.Unmarshal([]byte(data), &result)

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	} else {

		if result.Err == 0 && result.Data[0] == "d9545e34339afd76b6f669a45239f35f" {
			fmt.Println("沉默用户")
		} else {
			fmt.Println("非沉默用户")
		}
	}

	line := "(,BD2AF2CA714F4ABA154F37FE72A6B7F7,1D0EF26E0FECCCE89D67E9D7EC08C27F5E62030C)"
	line = strings.TrimPrefix(line, "(") // 文本格式: (idfa,imei,mac)
	line = strings.TrimSuffix(line, ")")

	arr := strings.Split(line, ",")

	fmt.Println(arr[0])
	fmt.Println(arr[1])
	fmt.Println(arr[2])

	idfa := "2D34EF1-8927-4438-93E7-6D87020364D6"

	if ok := idfaExp.MatchString(idfa); ok {
		fmt.Println("匹配")
	} else {
		fmt.Println("不匹配")
	}

	imei := "134567894324568"

	if ok := imeiExp.MatchString(imei); ok {
		fmt.Println("匹配")
	} else {
		fmt.Println("不匹配")
	}
}
