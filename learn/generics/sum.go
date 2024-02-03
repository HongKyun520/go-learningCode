package generics

// 泛型约束 Number接口下定义了泛型的范围
func Sum[T Number](vals ...T) T {

	var res T
	for _, val := range vals {
		res = res + val
	}
	return res
}

// 定义支持的类型 ~表示支持衍生类型
type Number interface {
	~int | ~int64 | ~float64 | ~float32
}
