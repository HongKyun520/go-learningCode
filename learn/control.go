package learn

// 控制结构

func IfOnly(age int) string {
	if age >= 18 {
		return "成年人"
	}
	return "小孩"
}

func IfElse(age int) string {
	if age >= 18 {
		return "成年人"
	} else {
		return "小孩"
	}
}

func IfElseIf(age int) string {
	if age >= 18 {
		return "成年"
	} else if age >= 12 {
		return "青少年"
	} else {
		return "小孩"
	}
}

// IfNewVariable Go的if-else支持一种新的写法，可以在if-else块里面定义一个新的局部变量
func IfNewVariable(start int, end int) string {

	if distance := end - start; distance > 100 {
		println(distance)
		return "距离太远了"
	} else {
		println(distance)
		return "距离比较远"
	}

}

// for后面不跟任何条件，相当于for true 死循环
func Loop3() {

	for {
		println("hello")
	}

}

// switch switch语句和别的语言类似
// default分支可以有也可以没有。 switch的值必须可比较
func Switch(status int) string {
	switch status {
	case 0:
		return "初始化"
	case 1:
		return "运行中"
	default:
		return "未知状态"
	}
}

func SwitchBool(age int) {
	switch {
	case age >= 18:
		println("123")
	case age > 12:
		println("1233")

	case age < 10:
		println("kkkkkk")
	}

}
