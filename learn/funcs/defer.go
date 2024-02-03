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
