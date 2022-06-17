package models

import (
	"flag"
	"github.com/jinzhu/gorm"
	"go_server/models/mongodb"
	"go_server/pkg/gredis"
	"go_server/pkg/log"
	"go_server/pkg/logrus"
	"go_server/pkg/setting"
	"go_server/service/tag_service"
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
	log.Init()
	//读取配置文件
	config := flag.String("c", "D:/fy-server/go_server/conf/app.ini", "config file path")
	//解析配置文件
	flag.Parse()
	setting.Init(*config) //根据配置文件初始化配置
	//初始化数据库
	Init()
	gredis.InitRedis()
	//r, err := SelectMutiTableByjoins()
	//logrus.Info("r =  ", err, r.TagName," ", r.TagState," ",r.ArticleTitle," ",r.ArticleDesc," ",r.ArticleContent)


/*	resp, err := util.SendCartMsg("ZhaoGuangChao", "推送消息", "付义+15602959486", "https://prerelease.ubtrobot.com/ADIS")
	if err != nil {
		return
	}
	logrus.Info("resp = ",resp)*/
	/*mongodb.InitDB()*/
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
