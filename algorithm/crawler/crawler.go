package crawler

import (
	"go_server/pkg/logrus"
	"io/ioutil"
	"net/http"
)

/**
多并发的爬虫项目
*/

func Init() {
	resp, err := http.Get("http://www.zhenai.com/zhenghun")
	if err != nil {
		logrus.Error(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logrus.Error(err)
		return
	}
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Info("获取网页：", all)
}
