package sort

import (
	"math/rand"
	"time"
)

const (
	num      = 10  //测试数组的长度
	rangeNUm = 100 //数组元素大小范围
)

/**
常见的排序算法
*/

//生成随机的数组
func GenerateRand() []int {
	randSeed := rand.New(rand.NewSource(time.Now().Unix() + time.Now().UnixNano()))
	arr := make([]int, num)
	for i := 0; i < num; i++ {
		arr[i] = randSeed.Intn(rangeNUm)
	}
	return arr
}

//比较两个切片是否相同
func isSame(slice1, slice2 []int) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	if (slice1 == nil) != (slice2 == nil) {
		return false
	}
	for i, num := range slice1 {
		if num != slice2[i] {
			return false
		}
	}
	return true
}

//冒泡排序
func Bubble(arr []int) {
	size := len(arr)
	var swapped bool
	for i := size - 1; i > 0; i-- {
		swapped = false
		for j := 0; j < i; j++ {
			if arr[j+1] < arr[j] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
				swapped = true
			}
		}
		if swapped != true {
			break
		}
	}
}

//选择排序
func SelectSort(arr []int) {
	for i := 0; i < len(arr); i++ {
		for j := i; j <= len(arr)-1; j++ {
			if arr[j] < arr[i] {
				arr[j], arr[i] = arr[i], arr[j]
			}
		}
	}
}

//插入排序
func InsertSort(arr []int) {
	for i := 1; i < len(arr); i++ {
		for j := i; j > 0; j-- {
			if arr[j-1] > arr[j] {
				arr[j-1], arr[j] = arr[j], arr[j-1]
			}
		}
	}
}

//快速排序
func QuickSort(arr []int, l, r int) {
	if l < r {
		pivot := arr[r]
		i := l - 1
		for j := l; j < r; j++ {
			if arr[j] <= pivot {
				i++
				arr[j], arr[i] = arr[i], arr[j]
			}
		}
		i++
		arr[r], arr[i] = arr[i], arr[r]
		QuickSort(arr, l, i-1)
		QuickSort(arr, i+1, r)
	}
}

//归并排序--合并

func Merge(arr []int, l, mid, r int) {//问题就在合并的地方
	//分别复制左右两个数组
	n1, n2 := mid-l+1, r-mid
	left, right := make([]int, n1), make([]int, n2)
	copy(left, arr[l:mid+1])
	copy(right, arr[mid+1:r+1])
	i, j := 0, 0
	k := 0
	for ; i < n1 && j < n2; k++ {
		if left[i] < right[j] {
			arr[k] = left[i]
			i++
		} else {
			arr[k] = right[j]
			j++
		}
	}
	for ; i < n1; i++ {
		arr[k] = left[i]
		k++
	}
	for ; j < n2; j++ {
		arr[k] = right[j]
		k++
	}
}

//分治
func MergeSort(arr []int, l, r int) {
	if l < r {
		mid := (l + r) / 2
		MergeSort(arr, l, mid)
		MergeSort(arr, mid+1, r)
		Merge(arr, l, mid, r)
	}
}

//堆调整
func adjust_heap(arr []int, i, size int) {
	if i <= (size-2)/2 {
		//左右子节点
		l, r := 2*i+1, 2*i+2
		m := i
		if l < size && arr[i] > arr[m] {
			m = l
		}
		if r < size && arr[r] > arr[m] {
			m = r
		}
		if m != i {
			arr[m], arr[i] = arr[i], arr[m]
			adjust_heap(arr, m, size)
		}
	}
}

//建堆
func build_heap(arr []int) {
	size := len(arr)
	//从最后一个字节开始向前调整
	for i := (size - 2) / 2; i >= 0; i-- {
		adjust_heap(arr, i, size)
	}
}

//堆排序
func HeapSort(arr []int) {
	size := len(arr)
	build_heap(arr)
	for i := size - 1; i > 0; i-- {
		//顶部arr[0]为当前最大值，调整到数组末尾
		arr[0], arr[i] = arr[i], arr[0]
		adjust_heap(arr, 0, i)
	}
}
