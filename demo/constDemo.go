package main 

import "fmt"

// 常量作为枚举类型
const (
	NetCodeSuccess = 200
	NetCodeFailure = 500
	NetCodeNotFound = 404
)

// 自增, 计数加1
const (
	GenderUnknown = iota
	GenderMale
	GenderFemale
)

func main() {
	const commonVal int = 404

	fmt.Println("无资源代码：", commonVal)

	fmt.Println("成功", NetCodeSuccess, "\n失败", NetCodeFailure, "\n找不到资源", NetCodeNotFound)


	fmt.Println("自增枚举使用", "\n未知性别", GenderUnknown, "\n男性", GenderMale, "\n女性", GenderFemale)

	for i := 0; i < 10; i++ {
		
	}
}

