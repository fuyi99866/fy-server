package initialize

import (
	gormadapter "github.com/casbin/gorm-adapter"
	"github.com/jinzhu/gorm"
)

var Casbin = new(casbin)

type casbin struct {
}

//创建需要初始化的管理员权限
var carbines = []gormadapter.CasbinRule{
	{PType: "p", V0: "admin", V1: "/swagger/*any", V2: "GET"},
	{PType: "p", V0: "admin", V1: "/auth", V2: "POST"},
	{PType: "p", V0: "admin", V1: "/ws", V2: "GET"},
	{PType: "p", V0: "admin", V1: "/alive", V2: "GET"},
	{PType: "p", V0: "admin", V1: "/user/test ", V2: "GET"},
	{PType: "p", V0: "admin", V1: "/user/:name ", V2: "GET"},
	{PType: "p", V0: "admin", V1: "/user", V2: "GET"},
	{PType: "p", V0: "admin", V1: "/user", V2: "POST"},
	{PType: "p", V0: "admin", V1: "/policy", V2: "POST"},
	{PType: "p", V0: "admin", V1: "/policy", V2: "DELETE"},
	{PType: "p", V0: "admin", V1: "/policy", V2: "GET"},
	{PType: "p", V0: "admin", V1: "/product", V2: "GET"},
}
var db *gorm.DB

func (c *casbin) Init() error {
	db.AutoMigrate(gormadapter.CasbinRule{})
	return db.Transaction(func(tx *gorm.DB) error {
		if tx.Find(&[]gormadapter.CasbinRule{}).RowsAffected == 12 {
			return nil
		}

		if err := tx.Create(&carbines).Error; err != nil { // 遇到错误时回滚事务
			return err
		}
		return nil
	})
}
