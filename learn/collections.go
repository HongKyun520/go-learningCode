package learn

import "fmt"

// 内置类型
// 数组、切片、map、channel

// [cap]type 初始化语法
// 1、初始化要指定长度
// 2、可以直接初始化
// 3、arr[i]的形式访问元素
// 4、len和cap操作用于获取数组长度

func Array() {
	a1 := [3]int{7, 8, 9}
	fmt.Printf("a1: %v, len: %d, cap: %d", a1, len(a1), cap(a1))

	a2 := [3]int{8, 9}
	fmt.Printf("a2: %v, len: %d, cap: %d", a2, len(a2), cap(a2))

	a3 := [3]int{}
	fmt.Printf("a3: %v, len: %d, cap: %d", a3, len(a3), cap(a3))

	// 数组不支持append操作
	// 按照下标索引，如果编译器能判断出来下表越界，那么就编译错误，如果不能，那么运行时候摆错，出现panic
	fmt.Printf("a[1]:%d", a1[2])

}

// Slice 切片，语法：[]type
// 1、直接初始化
// 2、make初始化：make([]type, length, capacity)
// 3、arr[i]的形式访问元素
// 4、append追加元素
// 5、len获取元素数量，cap获取切片容量
// 6、推荐写法 s1 := make([]type, 0, capacity)
// 在初始化切片的时候需要预估容量

func Slice() {
	// 初始化4个容量的切片，不需要指定容量
	s1 := []int{1, 2, 3, 4}
	fmt.Printf("s1:%v  len:%d, cap:%d \n", s1, len(s1), cap(s1))

	//使用make，初始化一个 3个元素，容量为4的切片
	s2 := make([]int, 3, 4)
	fmt.Printf("s2:%v, len:%d, cap:%d \n", s2, len(s2), cap(s2))

	// 追加一个元素 不触发扩容
	s2 = append(s2, 7)
	fmt.Printf("s2:%v, len:%d, cap:%d \n", s2, len(s2), cap(s2))

	// 再追加一个元素，触发扩容了
	s2 = append(s2, 8)
	fmt.Printf("s2:%v, len:%d, cap:%d \n", s2, len(s2), cap(s2))

	// make如果只穿入一个参数，表示创建了一个4个元素的切片
	// 以为：{}     实际上：{0,0,0,0}
	s3 := make([]int, 4)
	fmt.Printf("s3:%v, len:%d, cap:%d \n", s3, len(s3), cap(s3))

	// 按照下表进行索引
	fmt.Printf("s3[2]:%d", s3[2])

	// 超出下标返回，panic
	fmt.Printf("s3[2]:%d", s3[99])
}

// 子切片
// 数组和切片都可以通过[start:end]的形式来获取子切片 跟Java的String.subString()API类似  左闭右开

func SubSlice() {

	s1 := []int{2, 4, 6, 8, 10}
	s2 := s1[1:3]
	fmt.Printf("s2:%v, len:%d, cap:%d \n", s2, len(s2), cap(s2))

	s3 := s1[2:]
	fmt.Printf("s3:%v, len:%d, cap:%d \n", s3, len(s3), cap(s3))

	s4 := s1[:3]
	fmt.Printf("s4:%v, len:%d, cap:%d \n", s4, len(s4), cap(s4))

}

// 共享内存问题
// 核心：共享数组    子切片和切片究竟会不会互相影响，就抓住一点：它们是不是还共享数组？
// 如果结构未变化，则肯定共享，结构变化了，就可能不是共享的了。扩容了结构会发生变化，使用切片时需要关注扩容问题
func ShareSlice() {

	s1 := []int{1, 2, 3, 4}
	s2 := s1[2:]
	fmt.Printf("s1:%v, len:%d, cap:%d \n", s1, len(s1), cap(s1))
	fmt.Printf("s2:%v, len:%d, cap:%d \n", s2, len(s2), cap(s2))

	s2[0] = 99
	fmt.Printf("s1:%v, len:%d, cap:%d \n", s1, len(s1), cap(s1))
	fmt.Printf("s2:%v, len:%d, cap:%d \n", s2, len(s2), cap(s2))

	// 上面结果，s2的切片修改同步到了s1的切片的值，原因在于共享内存，s2未进行扩容，跟s1共享了内存

	// 下面，s2进行了扩容，使得修改互不影响
	s2 = append(s2, 199)
	fmt.Printf("s1:%v, len:%d, cap:%d \n", s1, len(s1), cap(s1))
	fmt.Printf("s2:%v, len:%d, cap:%d \n", s2, len(s2), cap(s2))

	s2[1] = 1999
	fmt.Printf("s1:%v, len:%d, cap:%d \n", s1, len(s1), cap(s1))
	fmt.Printf("s2:%v, len:%d, cap:%d \n", s2, len(s2), cap(s2))

}

// map  key-value数据结构
// 初始化：1、make方法，记得预估容量   2、直接初始化元素
// 赋值：使用中括号
// map的遍历是随机的，即遍历两次，输出结果是不一样的
func MapTest() {

	m1 := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	m2 := make(map[string]string, 2)
	m2["key2"] = "value2"

	// 读取元素，会有两个返回值
	val, ok := m1["key1"]
	if ok {
		println("第一步:", val)
	}

	val = m1["key2"]
	println("第二步:", val)

	println(len(m2))
	for k, v := range m1 {
		println(k, v)
	}

	for k := range m1 {
		println(k)
	}

	delete(m1, "key1")
}

// 一般实践
func Input(arr []int) {

	// 一般不要修改传入的array

}

func Map() {

	m1 := map[string]string{
		"key1": "123",
	}

	fmt.Println(m1)

	// 容量默认是16
	m2 := make(map[string]string, 8)

	fmt.Println(m2)

}
