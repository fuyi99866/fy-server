package util

import (
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"strings"
	"time"
)

func UUID() string {
	return uuid.NewV4().String()
}

func UUIDShort() string {
	return strings.ReplaceAll(uuid.NewV4().String(), "-", "")
}

//获取随机字符串
func GetRandomString(n int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < n; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
