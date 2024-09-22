package learn

var Global = "全局变量"

var internal = "局部变量"

const internalConst = "包内可访问"

const GlobalConst = "全局变量，包外可访问"

var (
	GlobalVariable   = "1"
	internalVariable = "2"
)

const (
	StatusA = iota
	StatusB
	StatusC
	StatusD
	StatusE

	// 主动赋值可以终端
	StatusF = 100
	StatusG
)

const (
	DayA = iota << 1
	DayB
	DayC
	DayD
	DayE
)

func variableTest() {
	// 局部变量
	var a int = 123
	println(a)

}
