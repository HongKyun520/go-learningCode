package types

import "fmt"

// 衍生类型
// 衍生类型和原类型可以相互转换，但仅仅可以转换而已，原类型的方法衍生类型无法访问

type Fish struct {
	Name string
}

func (f Fish) Swim() {
	fmt.Println("fish is swimming")
}

type FakeFish Fish

func UserFish() {

	f1 := Fish{Name: "123"}
	f2 := FakeFish(f1)

	f2.Name = "345"
	// 无法调用方法

	fmt.Println(f1)
	fmt.Println(f2)
}
