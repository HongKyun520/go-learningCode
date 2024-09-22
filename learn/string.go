package learn

import (
	"fmt"
	"unicode/utf8"
)

// 牢记strings包即可

func StringTest() {

	// 转义
	println("He said:\"hello, go!\"")

	println(`换行换行我
可以换行`)

	println(len("abc"))

	println(len("你好"))

	// 输出中文长度
	println(utf8.RuneCountInString("你好"))

	println(fmt.Sprintf("hello %d", 123))

	//var str string = "lkq.lpx"

}
