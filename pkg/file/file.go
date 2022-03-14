package file

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"strconv"
)

//获取文件大小
func GetSize(f multipart.File) (int, error) {
	content, err := ioutil.ReadAll(f)
	return len(content), err
}

//获取文件后缀
func GetExt(fileName string) string {
	return path.Ext(fileName)
}

//检查文件是否存在
func CheckExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

//检查文件权限
func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

//如果不存在则新建文件夹
func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist == true {
		if err := MkDir(src); err != nil {
			return err
		}
	}
	return nil
}

//新建文件夹
func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

//检查文件是否存在
func CheckNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

//打开文件
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return f, nil
}

//多次尝试打开文件
func MustOpen(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}
	src := dir + "/" + filePath
	perm := CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err = IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}
	f, err := Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return f, nil
}

var infile *string = flag.String("i", "unsorted.dat", "File contains values for sorting")
var outfile *string = flag.String("o", "sorted.dat", "File to receive sorted values")

//读取文件中的数据
func ReadValues(infile string) (values []int, err error) {
	file, err := os.Open(infile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	br := bufio.NewReader(file)
	values = make([]int, 0)
	for {
		line, isPrefix, err1 := br.ReadLine()
		if err1 != nil {
			if err1 != io.EOF {
				return nil, err
			}
			break
		}
		if isPrefix {
			return nil, errors.New("too long")
		}
		str := string(line) //转换字符数组为字符串
		value, err1 := strconv.Atoi(str)
		if err1 != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return
}

//将数字写入到文件
func WriteValues(values []int, outfile string) error {
	file, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, value := range values {
		str := strconv.Itoa(value)
		file.WriteString(str + "\n")
	}
	return nil
}
