package funcs

// 将方法赋值给变量

import (
	"fmt"
	"strings"
)

func Functional4() string {
	fmt.Println("hello, functional 4")
	return "hello"
}

func Functional5(a int) {
	fmt.Println("hello, functional 5")

}

func useFunctional4() {
	func4 := Functional4
	func4()

	func5 := Functional5
	func5(18)
}

// 返回一个，返回string的无参方法
func Functional6() func() string {
	return func() string {
		return "hello"
	}
}

// Func4 函数式编程
func Func4() {
	// 方法本身就可以赋值给某个变量，而变量就可以直接发起调用
	// 使用 := 的前提，就是左边必须有至少一个新变量
	myFunc3 := Func3
	_, _ = myFunc3(1, 2)
}

// Func5 函数式编程入门 - 局部方法
func Func5() {
	// 在方法内部声明一个局部变量，作用域就在方法内
	fn := func(str string) string {
		return "hello," + str
	}

	s := fn("kyun")
	fmt.Println(s)
}

// Func6 函数式编程入门 - 方法作为返回值。方法本身可以作为一个返回值
func Func6() func(str string) string {
	return func(str string) string {
		return "hello," + str
	}
}

func Func6Invoke() {
	func6 := Func6()
	s := func6("kyun")
	fmt.Println(s)
}

// Func7 函数式编程入门 - 匿名方法
// 在方法内部可以声明一个匿名方法 但是需要立即发起调用
// 为什么需要立刻发起调用？因为匿名，即没有名字，不立刻调用的话后面你都没办法调用了
// 匿名方法在defer中用的较多
func Func7() {

	// 新定义了一个方法并调用，因为方法后面有()，赋值给了fn， fn是方法的返回值，是一个string
	hello := func() string {
		return ""
	}()
	println(hello)
}

// 这里跟上面的区别，这里是定义了一个方法，并将其赋值给fn，方法并没有被调用
func functional8() {
	fn := func() string {
		return "hello"
	}
	fn()
}

func Func8(abc string) (string, int) {
	res := strings.Split(abc, "")
	return res[0], len(res)
}
