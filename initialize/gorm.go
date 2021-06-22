package initialize

import "github.com/jinzhu/gorm"

/**
初始化，连接数据库
*/

func Gorm() *gorm.DB {
	return GormMysql()
}

func GormMysql() *gorm.DB {

}
