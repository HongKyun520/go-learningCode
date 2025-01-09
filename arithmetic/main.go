package main

import "fmt"

func main() {
	arr := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}
	quickSort(arr, 0, len(arr)-1)
	fmt.Println(arr)
}

func quickSort(arr []int, left, right int) {
	if left < right {
		// 获取分区点
		pivot := partition(arr, left, right)
		// 递归排序左边部分
		quickSort(arr, left, pivot-1)
		// 递归排序右边部分
		quickSort(arr, pivot+1, right)
	}
}

// partition 函数用于分区,返回分区点位置
func partition(arr []int, left, right int) int {
	// 选择最右边元素作为基准值
	pivot := arr[right]
	// i 指向小于基准值区域的最后一个位置
	i := left - 1

	// 遍历区间,将小于基准值的元素放到左边
	for j := left; j < right; j++ {
		if arr[j] <= pivot {
			i++
			// 交换元素
			arr[i], arr[j] = arr[j], arr[i]
		}
	}

	// 将基准值放到正确位置
	arr[i+1], arr[right] = arr[right], arr[i+1]
	return i + 1
}
