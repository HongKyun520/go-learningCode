package learn

import "fmt"

// 牢记bytes包

func ByteTest() {

	var b byte = 'a'
	fmt.Println(b) // 打印结果 67 打印的是ASCII码
	fmt.Println(fmt.Sprintf("%v", b))

	// []byte和string是可以互相转换的
	var str string = "this is string"
	var bs []byte = []byte(str)
	println(bs)

}
