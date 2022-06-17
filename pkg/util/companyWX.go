package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	corpid  = "ww9af9a3bb3d5f994c"                          //企业ID
	agentId = "1000002"                                     //应用ID
	secret  = "CyhVJllMGfvXj-Fl9TsuXdjPqpWgAP8FFlPXsDq2EWg" //Secret
)

//企业微信应用消息提醒方法如下
func SendCartMsg(Tousers interface{}, title, description, url string) (map[string]interface{}, error) {
	btntex := "详情" //可以自定义卡片底下的文字

	qyurl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", corpid, secret)
	data, err := httpGetJson(qyurl)
	logrus.Info("data = ", data)
	if err != nil {
		logrus.Error(err)
		return data, err
	}
	errcode := data["errcode"].(float64)
	if errcode != 0 {
		logrus.Error(errcode)
		return make(map[string]interface{}), nil
	}
	access_token := data["access_token"]
	//卡片内容，不同类型消息通知以下内容需修改（可参考企业微信开发文档）
	req := map[string]interface{}{
		"touser":  Tousers,
		"msgtype": "textcard",
		"agentid": agentId,
		"textcard": map[string]interface{}{
			"title":       title,
			"description": description,
			"url":         url,
			"btntext":     btntex,
		},
	}
	sendurl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", access_token)
	data, err = httpPostJson(sendurl, req)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return data, nil
}

//封装了http请求的GET和POST方法，其返回值是map[string]interface{},方便我们使用
func httpGetJson(url string) (map[string]interface{}, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func httpPostJson(url string, data map[string]interface{}) (map[string]interface{}, error) {
	res, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(res))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	logrus.Info("resp = ", resp.Body)
	var data2 map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data2)
	if err != nil {
		return nil, err
	}
	return data2, nil
}
