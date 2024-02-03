package generics

// T 类型参数，名字叫做T，约束是any，等于没有约束  泛型
type List[T any] interface {
	Add(idx int, t T)
	Append(t T)
}

func UseList() {

	var l List[int]
	l.Append(123)

}

// 结构体也可以使用泛型
type LinkedList[t any] struct {
	head *node[t]
	val  t
}

type node[t any] struct {
	val t
}
