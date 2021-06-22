package DB

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//连接数据库
func MysqlTool() *gorm.DB {
	mysql, err := gorm.Open("mysql", "root:root@tcp(localhost:3306)/demo?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("connect to DB error")
		panic(err)
	}
	return mysql
}
