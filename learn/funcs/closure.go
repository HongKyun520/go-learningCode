package funcs

import (
	"fmt"
	"strconv"
)

// 闭包，方法 + 它绑定的运行上下文    Go语言的匿名函数
// 闭包如果使用不当可能会引起内存泄露问题，即一个对象被闭包引用的话，它是不会被垃圾回收的

func Closure(name string) func() string {
	return func() string {
		return "hello" + name
	}
}

func Clousure1() func() string {
	age := 10
	fmt.Printf("out %p", &age)
	return func() string {
		age++
		fmt.Printf("in %p", &age)
		return strconv.Itoa(age)
	}
}
