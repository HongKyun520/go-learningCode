package learn

import (
	"database/sql"
	"fmt"
)

// 闭包

type MyStruct struct {
	Name string
	age  int
}

// Closure 闭包，方法 + 它绑定的运行上下文。 闭包如果使用不当可能会引起内存泄露的问题。即一个对象被闭包引用的话，它是不会被垃圾回收的
func Closure(str string) func() string {
	return func() string {
		return "hello" + str
	}
}

// YourName 不定参数 不定参数是指最后一个参数，可以传入任意多个值，注意必须是最后一个参数才可以声明为不定参数
// 不定参数在方法内部可以被当成切片来使用
func YourName(name string, alias ...string) {
	if len(alias) > 0 {
		println(alias[0])
	} else {
		println(name)
	}
}

func YourNameInvoke() {
	YourName("test")
	YourName("test", "01", "02")
	YourName("test", "01", "02", "03")
	YourName("test", "01", "02", "04", "05")
}

// defer go语言机制，允许在方法返回前一刻。执行一段逻辑
// defer类似斩，先定义后执行，后定义先执行   一个方法不能超过8个defer

func Defer() {
	defer func() {
		println("第一个 defer ...")
	}()

	defer func() {
		println("第二个 defer ...")
	}()
}

// DeferClosure defer与闭包  结果 i = 1
// 确定值原则
// 作为参数传入的：定义 defer 的时候就确定了
// 作为闭包引入的：执行 defer 对应的方法的时候才确定
func DeferClosure() {
	i := 0
	defer func() {
		println(i)
	}()
	i++
}

func DeferClosureV1() {
	i := 0
	defer func(val int) {
		println(i)
	}(i)
	i++

	println("外部", i)

}

// DeferReturn 如果是带名字的返回值，那么可以修改这个返回值，否则不能修改。
// DeferReturn defer修改返回值 这里defer并没有修改到a的值，a = 0
func DeferReturn() int {
	a := 0
	defer func() {
		a = 1
	}()
	return a
}

// DeferReturnV1 return a = 2
func DeferReturnV1() (a int) {
	a = 0
	defer func() {
		a = 2
	}()

	return a
}

func DeferReturnV2() *MyStruct {
	a := &MyStruct{
		Name: "kyun",
		age:  12,
	}

	defer func() {
		a.Name = "lkq"
	}()

	return a
}

// DeferClosureLoopV1 闭包与循环
// 10 10 10 10  for循环内的形参i的地址是不变的
func DeferClosureLoopV1() {

	for i := 0; i < 10; i++ {
		defer func() {
			fmt.Printf("%p", &i)
			println("print_", i)
		}()
	}
}

// DeferClosureLoopV2
// 9 8 7 6 .....   每一次循环，i当作入参复制了一份传了进去
func DeferClosureLoopV2() {
	for i := 0; i < 10; i++ {
		defer func(val int) {
			fmt.Printf("%p", &val)
			println("print_", val)
		}(i)
	}
}

// 9 8 7 6 ....   i赋值给了新的变量j
func DeferClosureLoopV3() {
	for i := 0; i < 10; i++ {
		j := i
		defer func() {
			fmt.Printf("%p", &j)
			println("print_", j)
		}()
	}

}

// 释放数据库资源
func Query() {
	db, _ := sql.Open("", "")
	defer db.Close()

	db.Query("SELECT")
}
