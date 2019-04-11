package main

import "fmt"



func main() {
	var age int = 1
	name := "我的年龄是"
	fmt.Println("hello, world")
	fmt.Println("你好啊，世界")
	fmt.Println(name, age)


	var prefixString, nameString, ageString, ageInt = "你好，", "我是章三", "我的年龄是8岁", 8

	fmt.Println(prefixString, nameString, ageString, ageInt)


	songName, singer, lyrics := "椿", "沈以诚", "恍惚间，浸透了回忆"

	fmt.Println("这是一首歌，\n歌名：", songName, "\n歌手：", singer, "\n歌词：", lyrics)


	_, secondNumber, lastString := canOmittedParameters()
	fmt.Println(secondNumber, lastString)

}

// 可省略参数函数体
func canOmittedParameters()(int, int, string) {
	return 12, 14, "哈哈哈"
}