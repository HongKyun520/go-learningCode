package week1

import (
	"errors"
	"fmt"
)

// 切片作业
var ErrIndexOutOfRange = errors.New("下标超出范围")

// DeleteAt 删除指定位置的元素
// 如果下标不是合法的下标，返回ErrIndexOutOfRange
func DeleteAt[T any](src []T, index int) ([]T, error) {

	// 判断切片的长度
	length := len(src)
	if index < 0 || index >= length {
		return nil, fmt.Errorf("ekit: %w, 下标超出范围, 长度 %d, 下标 %d", ErrIndexOutOfRange, length, index)
	}
	// 将index下的元素删除，然后index后的元素往前移
	for i := index; i+1 < length; i++ {
		src[i] = src[i+1]
	}

	// 返回指定下的切片
	return src[:length-1], nil
}

// 切片的缩容
func Shrink[T any](src []T) []T {
	c, l := cap(src), len(src)
	n, changed := calCapacity(c, l)

	// 不需要缩容，则直接返回原始切片
	if !changed {
		return src
	}

	s := make([]T, 0, n)
	s = append(s, src...)
	return s
}

// 判断容量
func calCapacity(c, l int) (int, bool) {
	// 容量 <= 64 缩不缩无所谓， 因为内存也浪费不了多少
	if c <= 64 {
		return c, false
	}

	if c > 2048 && (c/l >= 2) {
		factor := 0.625
		return int(float32(c) * float32(factor)), true
	}

	// 如果在2048以内，并且元素不足 1/4 ，那么直接缩减为一半
	if c <= 2048 && (c/l >= 4) {
		return c / 2, true
	}

	// 整个实现的核心是希望在后续少触发扩容的前提下，一次性释放尽可能多的内存
	return c, false
}
