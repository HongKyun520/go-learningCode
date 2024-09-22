package buildin_types

import "fmt"

func GetAllKeyAndValue() {
	map1 := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for k, v := range map1 {
		fmt.Printf("key: %s, value: %d\n", k, v)
	}
}
