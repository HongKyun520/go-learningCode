package funcs

func YourName(name string, aliases ...string) {

}

func callYourName() {

	YourName("小米")
	YourName("小米", "小明")
	YourName("小米", "雄安秘密")

	// 切片
	aliases := []string{"大明", "小明"}
	YourName("大明", aliases...)

}
