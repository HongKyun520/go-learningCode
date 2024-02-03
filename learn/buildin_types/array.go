package buildin_types

import "fmt"

func Array() {

	// int[] a = {1, 2, 3}
	// [5]int{1,2,3,4,5}

	a1 := [3]int{9, 8, 7}
	fmt.Printf("a1: %v, len=%d, cap=%d  \n", a1, len(a1), cap(a1))

}
