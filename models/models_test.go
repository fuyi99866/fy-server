package models

import (
	"flag"
	"go_server/pkg/gredis"
	"go_server/pkg/logger"
	"go_server/pkg/setting"
	"regexp"
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

	//TODO 开始测试数据库操作

	re:=regexp.MustCompile(`([0-9]+)B-([0-9]+)F`)
	match:=re.FindAllStringSubmatch("10B-11F",-1)
	logger.Info("match = ",match[0][1],match[0][2])
}

/*db.Table("go_service_info").Select("go_service_info.serviceId as service_id, go_service_info.serviceName as service_name,
go_system_info.systemId as system_id, go_system_info.systemName as system_name")
.Joins("left join go_system_info on go_service_info.systemId = go_system_info.systemId").Scan(&results)*/

//通过joins联合查询
func SelectMutiTableByjoins() (*Result, error) {
	var (
		rst Result
	)
	//article.title as article_title,article.desc as article_desc,article.content as article_content,tag.name as tag_name,tag.state as tag_state
	//db.Joins("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "jinzhu@example.org").Find(&user)
	err := db.Table("article").Select("article.title as article_title,article.desc as article_desc,article.content as article_content,tag.name as tag_name,tag.state as tag_state").Joins("left join tag on tag.id = article.tag_id").Scan(&rst).Error
	return &rst, err
}
