package test1_test_test

import "code/server/test/test1"

// 检测单行输出
func ExampleSayHello() {
	test1.SayHello()
	// OutPut: Hello World
}

// 检测多行输出
func ExampleSayGoodbye() {
	test1.SayGoodbye()
	// OutPut:
	// Hello,
	// goodbye
}

// 检测乱序输出
func ExamplePrintNames() {
	test1.PrintNames()
	// Unordered output:
	// Jim
	// Bob
	// Tom
	// Sue
}



