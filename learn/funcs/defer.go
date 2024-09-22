package funcs

// 延迟调用
// 执行顺序，第二个 到 第一个  defer有自己的栈，出函数时出栈执行
func Defer() {

	defer func() {
		println("第一个 defer")
	}()

	defer func() {
		println("第二个 defer")
	}()

}

func DeferClosureV1() {

	i := 0

	// 闭包的方式进行defer， defer的参数是函数执行时传入的，不是函数定义时传入的，所以i值为1
	defer func() {
		println(i)
	}()

	i = 1
}

func DeferClosureV2() {

	i := 0

	// 传值的方式进行defer， defer的参数是定义时传入的，所以i值为0
	defer func(val int) {
		println(val)
	}(i)

	i = 1
}

// 不带名字的返回值，使用defer无法修改返回值
func DeferReturn() int {
	a := 0
	defer func() {
		a = 1
	}()

	return a
}

// 带名字的返回值，使用defer可以修改返回值
func DeferReturnNamed() (a int) {
	a = 0
	defer func() {
		a = 1
	}()

	return a
}

// 使用defer修改结构体的返回值
func DeferReturnStructV1() *A {
	a := &A{
		name: "a",
	}
	defer func() {
		a.name = "Tom"
	}()

	return a
}

// 使用defer修改结构体的返回值
func DeferReturnStructV2() (a *A) {
	a = &A{
		name: "a",
	}
	defer func() {
		a.name = "Tom"
	}()

	return a
}

type A struct {
	name string
}
