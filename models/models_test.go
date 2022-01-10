package models

import (
	"flag"
	"go_server/pkg/gredis"
	"go_server/pkg/logger"
	"go_server/pkg/setting"
	"testing"
)

type Result struct {
	ArticleTitle   string `json:"article_title"`
	ArticleDesc    string `json:"article_desc"`
	ArticleContent string `json:"article_content"`
	TagName        string `json:"tag_name"`
	TagState       string `json:"tag_state"`
}

func TestAll(t *testing.T) {
	//读取配置文件
	config := flag.String("c", "D://fy-server//go_server//conf/app.ini", "config file path")
	//解析配置文件
	flag.Parse()
	setting.Init(*config) //根据配置文件初始化配置

	//初始化日志系统
	logger.InitLog1(setting.AppSetting.LogLever, "./logs/go_server.log") //初始化日志库 ,使用zap库

	//初始化数据库
	Init()
	gredis.InitRedis()
	r, err := SelectMutiTableByjoins()
	logger.Info("r =  ", err," ",r.ArticleTitle," ", r.ArticleContent," ", r.ArticleDesc," ", r.TagName," ", r.TagState)
}

//通过joins联合查询
func SelectMutiTableByjoins() (*Result, error) {
	//rst := Result{}
	var rst Result
	err := db.Model(Article{}).Select("article.title as article_title,article.desc as article_desc," +
		"article.content as article_content,tag.name as tag_name,tag.state as tag_state").
		Joins("Left JOIN tag ON tag.id = article.tag_id").Scan(&rst).Error
	return &rst, err
}
