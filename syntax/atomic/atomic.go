package main

import "sync/atomic"

func main() {

	// 声明一个int32类型的变量val,初始值为12
	var val int32 = 12
	// 原子操作:读取val的值
	val = atomic.LoadInt32(&val)
	println(val)

	// 原子操作:将val的值设置为13
	atomic.StoreInt32(&val, 13)
	println(val)

	// 原子操作:将val的值加1,并返回新值
	newVal := atomic.AddInt32(&val, 1)
	println(newVal)

	// CAS(Compare And Swap)操作:
	// 比较val的值是否为13,如果是则将其设置为14
	// 返回是否交换成功
	result := atomic.CompareAndSwapInt32(&val, 13, 14)
	println(result)

}
