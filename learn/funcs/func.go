package funcs

// Func0 单返回值
func Func0(name string) string {
	return "hello world"
}

// Func1 多返回值
func Func1(a, b, c int) (string, error) {
	return "", nil
}

// Func2 带名字的返回值，相当于在函数中创建一个局部变量，可以直接返回
func Func2(a int, b int) (str string, err error) {
	str = "result"
	return
	//return str, err
}

func Func3(a, b int) (str string, err error) {
	res := "result"
	return res, nil
}
