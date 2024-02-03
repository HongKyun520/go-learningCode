package component

type Inner struct {
}

func (i Inner) doSomething() {

}

type Outer struct {
	Inner
}

type OuterPtr struct {
	*Inner
}

func UserOuter() {

	var o Outer
	// 可以调用组合的方法
	o.doSomething()
	o.Inner.doSomething()

	var op *OuterPtr
	op.doSomething()

	o1 := Outer{Inner: Inner{}}
	o1.doSomething()

	op1 := OuterPtr{Inner: &Inner{}}
	op1.doSomething()

	// 调用顺序 先看是否自己有这个方法，没有则调用组合结构体内的这个方法
}
