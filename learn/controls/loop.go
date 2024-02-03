package controls

import "fmt"

func Loop2() {

	i := 0

	for i < 10 {
		i++
		println(i)
	}

	// 死循环
	for true {
		i++
		println(i)
	}

	for {
		i++
		println(i)
	}

}

func LoopArr() {

	arr := [3]int{1, 2, 3}
	for index, val := range arr {
		println("下标", index, "值", val)
	}

	for index := range arr {
		println("下标", index, "值", arr[index])
	}

}

// map的遍历是随机的
func LoopMap() {

	m := map[string]int{
		"key1": 100,
		"key2": 102,
	}

	for k, v := range m {
		println(k, v)
	}

	for k := range m {
		println(k, m[k])
	}

}

// 千万不要对迭代参数取地址！！！
// 在内存里面，迭代参数都是放在一个地址上的
func LoopBug() {

	users := []User{
		{
			name: "大明",
		},
		{
			name: "小明",
		},
	}

	m := make(map[string]*User)

	for _, u := range users {
		m[u.name] = &u
	}

	fmt.Printf("%v", m)

}

// 99乘法表
func multiplication() {

	for i := 1; i < 10; i++ {
		for j := 1; j <= i; j++ {
			fmt.Printf("%d*%d=%d ", j, i, j*i)
		}
		fmt.Println("")
	}
}

type User struct {
	name string
}
