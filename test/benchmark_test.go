package test

import (
	"fmt"
	"testing"
)

//基准测试函数
func BenchmarkAddNum(b *testing.B) {
	//重置计时器
	//b.ResetTimer()

	//停止计时器
	//b.StopTimer()

	//开始计时器
	//b.StartTimer()

	var n int
	for i := 0; i < b.N; i++ {
		//fmt.Sprintf("%d", i)
		n++
	}
}

//测试内存
func Benchmark_Alloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("%d", i)
	}
}

//控制计时器
func Benchmark_Add_TimerControl(b *testing.B) {
	//重置计时器
	b.ResetTimer()

	//停止计时器
	b.StopTimer()

	//开始计时器
	b.StartTimer()

	var n int
	for i := 0; i < b.N; i++ {
		n++
	}
}
