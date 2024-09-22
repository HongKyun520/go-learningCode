package funcs

import "fmt"

// 递归  使用不当报栈溢出
func Recursive(n int) {

	fmt.Println(n)
	if n > 10 {
		return
	}

	Recursive(n + 1)
}
