package test

import (
	"fmt"
)

// SayHello 打印一行字符串
func Sayhello() {
	fmt.Println("hello world")
}

// SayGoodbye 打印两行字符串
func SayGoodbye() {
	fmt.Println("hello")
	fmt.Println("goodbye")
}

// PrintNames 打印学生姓名
func PrintNames() {
	student := make(map[int]string, 4)
	student[1] = "Jim"
	student[2] = "Bob"
	student[3] = "Tom"
	student[4] = "Sue"
	for _, value := range student {
		fmt.Println(value)
	}
}

//计算累加函数
func AddNum(n int) (result int) {
	for i := 0; i < n; i++ {
		result = result + i
	}
	return
}
