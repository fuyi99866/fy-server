package util

import (
	"fmt"
	"strconv"
	"strings"
)

func Ints2string(spit string, arr []int) string {
	str := ""
	for i, v := range arr {
		c := strconv.Itoa(v)
		if i != 0 {
			str = fmt.Sprintf("%s%s%s", str, spit, c)
		} else {
			str = c
		}
	}
	return str
}

func String2ints(spit, s string) []int {
	if s == "" {
		return []int{}
	}
	strArr := strings.Split(s, spit)
	var arrInt []int
	for _, v := range strArr {
		num, _ := strconv.Atoi(v)
		arrInt = append(arrInt, num)
	}
	return arrInt
}
