package funcs

import "fmt"

func YourName(name string, aliases ...string) {
	fmt.Println(name, aliases)
}

func CallYourName() {

	YourName("小米")
	YourName("小米", "小明")
	YourName("小米", "雄安秘密")

	// 切片
	aliases := []string{"大明", "小明"}
	YourName("大明", aliases...)

}
