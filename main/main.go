package main

import (
	"GoInAction/learn/generics"
	"fmt"
)

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
	//println("Hello World!")
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

	fmt.Println(generics.Sum[int](1, 2, 3))
	fmt.Println(generics.Sum[float32](1.0, 2.0, 3.0))
	fmt.Println(generics.Sum[int](1, 2, 3))

}
