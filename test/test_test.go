package test

import (
	"testing"
)

// 检测单行输出
func TestSayhello(t *testing.T) {
	t.Log("hello 1111")
}

// 检测多行输出
func TestPrintNames(t *testing.T) {
	t.Log("say goodbye")
}


// 检测乱序输出
func TestPrintNames2(t *testing.T) {
	t.Log("print names")
}


