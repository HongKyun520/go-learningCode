package practice

import "fmt"

// 1、Go 是强类型语言，能不能设计一个方法，可以计算任意数字类型切片的和的方法？
// func SumInt64([]int64)
// func SumInt32([]int32)
// 调用SumInt64()必须传对类型
func sumSlice(a Adder) float64 {
	return a.Add()
}

// 2、获得 map 的所有 key、所有 value
func getMapAllValue() {

	m1 := make(map[string]string, 16)

	m1["key1"] = "123"
	m1["key2"] = "456"
	for k, v := range m1 {
		fmt.Println("key:", k, "value:", v)
	}
}

type Adder interface {
	Add() float64
}

type Integer32Slice []int32

type Integer64Slice []int64

func (s Integer32Slice) Add() float64 {
	var sum int32 = 0
	for _, num := range s {
		sum += num
	}
	return float64(sum)
}

func (s Integer64Slice) Add() float64 {
	var sum int64 = 0
	for _, num := range s {
		sum += num
	}
	return float64(sum)
}
