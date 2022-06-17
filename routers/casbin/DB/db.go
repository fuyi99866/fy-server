package DB

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"go_server/pkg/setting"
)

//连接数据库
func MysqlTool() *gorm.DB {
	var err error
	var dataPath string
	if setting.DatabaseSetting.Type == "mysql" {
		dataPath = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			setting.DatabaseSetting.User,
			setting.DatabaseSetting.Password,
			setting.DatabaseSetting.Host,
			setting.DatabaseSetting.Name)
	} else if setting.DatabaseSetting.Type == "sqlite3" {
		logrus.Info("dataType = ", setting.DatabaseSetting.Type)
		dataPath = "data/test.db"
	}
	db, err := gorm.Open(setting.DatabaseSetting.Type, dataPath)
	if err != nil {
		logrus.Fatal("无法连接数据库... err: %/v", err)
	}
	return db
}
