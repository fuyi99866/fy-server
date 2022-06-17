package util

import (
	"bytes"
	"github.com/eventials/go-tus"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type Reader struct {
	io.Reader
	Total   int64
	Current int64
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	r.Current += int64(n)
	logrus.Info("progress = ", r.Current*100/r.Total)
	return
}

//获取文件下载进度
func DownloadFileProgress(url, filename string) {
	r, err := http.Get(url)
	if err != nil {
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	f, err := os.Create(filename)
	if err != nil {
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	reader := &Reader{
		Reader: r.Body,
		Total:  r.ContentLength,
	}

	n, err := io.Copy(f, reader)
	logrus.Info("n ,err = ", n, err)
}

//下载文件
func HttpGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		logrus.Info("HttpGet failed! ", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Info("ReadAll failed! ", err)
		return
	}

	logrus.Info("get response :", string(body))

}

//go-fastdfs断点续传
func UploadByBreakPoint(filepath, uri string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// create the tus client.
	client, err := tus.NewClient(uri, nil) //"http://10.10.17.15:8087/group1/big/upload/"
	if err != nil {
		logrus.Error("创建客户端失败")
		return "", err
	}

	// create an upload from a file.
	upload, err := tus.NewUploadFromFile(f)
	if err != nil {
		logrus.Error("上传失败")
		return "", err
	}

	//输出进度条
	go func() {
		for upload.Progress() < 100 {
			logrus.Info("upload.Progress() = ", upload.Progress()) //文件大小
			time.Sleep(time.Duration(time.Millisecond * 500))
		}
	}()

	// create the uploader.
	uploader, err := client.CreateUpload(upload)
	if err != nil {
		logrus.Error("获取进度失败")
		return "", err
	}

	// start the uploading process.
	err = uploader.Upload()
	if err != nil {
		logrus.Error("Upload failed")
		return "", err
	}

	if upload.Progress() == 100 {
		logrus.Info("upload.Progress() = ", upload.Progress()) //文件大小
		logrus.Info("文件上传成功 ")
	}

	//上传完成后，再通过秒传接口，获取文件信息
	//http://127.0.0.1:8080/upload?md5=430a71f5c5e093a105452819cc48cc9c&output=json

	return uploader.Url(), nil
}

//秒传文件
func HttpGetBySecond(uri string) ([]byte, error) {
	logrus.Info("uri = ", uri)
	md5 := string([]byte(uri)[len(uri)-32 : len(uri)])
	logrus.Info("md5 = ", md5)
	ur := "http://10.10.17.15:8087/group1/upload?md5=" + md5 + "&output=json"
	logrus.Info("ur = ", ur)

	resp, err := http.Get(ur)
	defer resp.Body.Close()
	if err != nil {
		logrus.Info("HttpGet failed! ", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Info("ReadAll failed! ", err)
		return nil, err
	}

	logrus.Info("get response :", string(body))
	return body, nil
}

//上传文件
func HttpPost(filepath, uri string) ([]byte, error) {

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	logrus.Info("file = ", file.Name())
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", uri, body)

	//设置header
	req.Header.Set("Content-Type", writer.FormDataContentType())

	//Do方法发送请求
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("发送POST请求失败: ", err)
		return nil, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("ReadAll failed: ", err)
		return nil, err
	}

	return b, nil
}

//模拟表单发送POST请求
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, path)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", uri, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	return request, err
}
