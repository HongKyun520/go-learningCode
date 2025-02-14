package main

import "fmt"

// 常量
const PI = 3.14

// 全局变量的声明和赋值
var name = "gopher"

/**
通过 const 关键字来进行常量的定义。

通过在函数体外部使用 var 关键字来进行全局变量的声明和赋值。

通过 type 关键字来进行结构(struct)和接口(interface)的声明。

通过 func 关键字来进行函数的声明。
*/

func main() {

	naturals := make(chan int)

	squares := make(chan int)

	// 生成自然数
	go func() {
		fmt.Println("naturals 开始	")
		defer close(naturals)
		for i := 0; i < 10; i++ {
			naturals <- i
		}
		fmt.Println("naturals 结束")
	}()

	// 计算自然数的平方
	go func() {
		fmt.Println("squares 开始")
		defer close(squares)
		for natural := range naturals {
			squares <- natural * natural
		}
		fmt.Println("squares 结束")
	}()

	for square := range squares {
		println(square)
	}

	//println("Hello World!")

	//a1 := DeferReturnStructV1()
	//fmt.Println("DeferReturnStructV1:", a1.name) // 输出: Tom
	//
	//a2 := DeferReturnStructV2()
	//fmt.Println("DeferReturnStructV2:", a2.name)

	//DeferClosureLoopV3()
	//buildin_types.GetAllKeyAndValue()

	//i := 32
	//funcs.Defer()
	//funcs.DeferClosureV1()
	//funcs.DeferClosureV2()

	//funcs.Recursive(2)

	//println(learn.Global)
	//println(learn.GlobalConst)

	//learn.Extremum()
	//learn.StringTest()
	//learn.Defer()
	//learn.DeferClosure()
	//learn.DeferClosureV1()
	//a := learn.DeferReturnV2()
	//fmt.Println(a.Name)

	//learn.DeferClosureLoopV1()
	//learn.DeferClosureLoopV2()
	//learn.DeferClosureLoopV3()
	//variable := learn.IfNewVariable(100, 300)
	//println(variable)

	//learn.Array()
	//learn.Slice()
	//learn.SubSlice()
	//learn.ShareSlice()
	//learn.MapTest()

	//funcs.Recursive(2)
	//name := funcs.Closure("康权")
	//fmt.Println(name())

	//getAge := funcs.Clousure1()
	//fmt.Println(getAge())
	//controls.LoopBug()
	//learn.Array()

	//types.UserFish()

	//fmt.Println(generics.Sum[int](1, 2, 3))
	//fmt.Println(generics.Sum[float32](1.0, 2.0, 3.0))
	//fmt.Println(generics.Sum[int](1, 2, 3))
	//learn.StringTest()

	//learn.ByteTest()

	//learn.StructTest()
}

type A struct {
	name string
}

func DeferReturnStructV1() *A {
	a := &A{
		name: "a",
	}
	defer func() {
		a.name = "Tom"
	}()

	return a
}

func DeferReturnStructV2() (a *A) {
	a = &A{
		name: "a",
	}
	defer func() {
		a.name = "Tom"
	}()

	return a
}

func DeferClosureLoopV3() {

	for i := 0; i < 10; i++ {
		j := i

		defer func() {
			println(j)
		}()

	}

}
