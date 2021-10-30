package crawler

import (
	"go_server/pkg/logger"
	"io/ioutil"
	"net/http"
)

/**
多并发的爬虫项目
*/

func Init() {
	resp, err := http.Get("https://www.zhenai.com/zhenhun")
	if err != nil {
		logger.Error(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logger.Error(err)
		return
	}
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("获取网页：", all)
}
