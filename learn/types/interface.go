package types

import "fmt"

// 接口基本语法：type名字interface{}
// 接口只能包含方法

type User struct {
	Name string
	Age  int
}

type List interface {
	Add(index int, val any)
	Append(val any)
	Delete(val any)
}

// StructTest 结构体的初始化
// Go是没有构造函数的！
// 初始化语法 Struct{}
// 获取指针 &Struct{}  new(Struct)  new可以理解为Go会为你的变量分配内存，并且把内存都置为0
func StructTest() {
	//指针 打印地址
	u1 := &User{}
	fmt.Println(u1)

	u6 := new(User)
	fmt.Println(u6)

	// u2字段中都是零值
	u2 := User{}
	fmt.Printf("%v \n", u2)
	fmt.Printf("%+v \n", u2)

	// u3中的字段也都是零值
	var u3 User
	fmt.Println(u3)

	// 初始化 按照字段赋值
	var u4 User = User{Name: "Tom"}
	fmt.Println(u4)

	// 初始化 按照顺序赋值，必须全部赋值
	var u5 User = User{"Tom", 18}
	fmt.Println(u5)

	// 这个是nil
	var l3Ptr *LinkedList
	fmt.Println(l3Ptr)
}

// 结构体接收器。值传递
func (u User) ChangeName(name string) {
	fmt.Println(name)
	u.Name = name
}

func (u User) ChangeAge(age int) {
	fmt.Println(age)
	u.Age = age
}

// 指针接收器，传递指针
func (u *User) ChangeNameByPointer(name string) {
	fmt.Println(name)
	u.Name = name
}

// 指针接收器，传递指针
func (u *User) ChangeAgeByPointer(age int) {
	fmt.Println(age)
	u.Age = age
}

func ChangeUser() {

	u1 := User{
		Name: "Tom",
		Age:  18,
	}

	// 指针和结构体调用是兼容的
	u1.ChangeName("123")
	u1.ChangeNameByPointer("lkq")

	u2 := &User{
		Name: "Tom",
		Age:  19,
	}

	// 指针和结构体调用是兼容的
	u2.ChangeAge(12)
	u2.ChangeAgeByPointer(13)

}

// 结构体自引用 (设计自身的数据结构)
// 如果在结构体内部还要引用自身，那么只能使用指针.准确来说，在整个调用链上，如果构成循环，那就只能用指针
// 叶子节点
type node struct {
	prev *node
	next *node
}
