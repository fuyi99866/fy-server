package sort

import (
	"fmt"
	"testing"
)

func TestBubble(t *testing.T) {
	arr := GenerateRand()
	Bubble(arr)
	for _, v := range arr {
		fmt.Println(v, " ")
	}
}

func TestInsertSort(t *testing.T) {
	arr := GenerateRand()
	InsertSort(arr)
	for _, v := range arr {
		fmt.Println(v, " ")
	}
}

func TestHeapSort(t *testing.T) {
	arr := GenerateRand()
	HeapSort(arr)
	for _, v := range arr {
		fmt.Println(v, " ")
	}
}

func TestMergeSort(t *testing.T) {
	arr := GenerateRand()
	for _, v := range arr {
		fmt.Println(v, " ")
	}
	fmt.Println("-----------------------------")
	MergeSort(arr, 0, len(arr)-1)
	for _, v := range arr {
		fmt.Println(v, " ")
	}
}

func TestQuickSort(t *testing.T) {
	arr := GenerateRand()
	l := 0
	r := len(arr) - 1
	QuickSort(arr, l, r)
	for _, v := range arr {
		fmt.Println(v, " ")
	}
}

func TestSelectSort(t *testing.T) {
	arr := GenerateRand()
	SelectSort(arr)
	for _, v := range arr {
		fmt.Println(v, " ")
	}
}
