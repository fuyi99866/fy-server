package models

import (
	gormadapter "github.com/casbin/gorm-adapter"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

var Casbin = new(casbin)

type casbin struct{}

var casbinrule = []gormadapter.CasbinRule{
	{PType: "p", V0: "admin", V1: "/swagger/*any", V2: "GET"},               //GET    /swagger/*any
	{PType: "p", V0: "admin", V1: "/auth", V2: "POST"},                      //POST   /auth
	{PType: "p", V0: "admin", V1: "/upload", V2: "POST"},                    //POST   /upload
	{PType: "p", V0: "admin", V1: "/ws", V2: "GET"},                         //GET    /ws
	{PType: "p", V0: "admin", V1: "/alive", V2: "GET"},                      //GET    /alive
	{PType: "p", V0: "admin", V1: "/user/test", V2: "GET"},                  //GET    /user/test
	{PType: "p", V0: "admin", V1: "/user/:name", V2: "GET"},                 //GET    /user/:name
	{PType: "p", V0: "admin", V1: "/user", V2: "GET"},                       //GET    /user
	{PType: "p", V0: "admin", V1: "/user", V2: "POST"},                      //POST   /user
	{PType: "p", V0: "admin", V1: "/user/delete", V2: "DELETE"},             //DELETE /user/delete
	{PType: "p", V0: "admin", V1: "/user/update", V2: "POST"},               //POST   /user/update
	{PType: "p", V0: "admin", V1: "/user/get", V2: "GET"},                   //GET    /user/get
	{PType: "p", V0: "admin", V1: "/policy", V2: "POST"},                    //POST   /policy
	{PType: "p", V0: "admin", V1: "/policy", V2: "DELETE"},                  //DELETE /policy
	{PType: "p", V0: "admin", V1: "/policy", V2: "GET"},                     //GET    /policy
	{PType: "p", V0: "admin", V1: "/product", V2: "GET"},                    //GET    /product
	{PType: "p", V0: "admin", V1: "/authority/add", V2: "POST"},             //POST   /authority/add
	{PType: "p", V0: "admin", V1: "/authority/update", V2: "POST"},          //POST   /authority/update
	{PType: "p", V0: "admin", V1: "/authority/set", V2: "POST"},             //POST   /authority/set
	{PType: "p", V0: "admin", V1: "/authority/delete", V2: "DELETE"},        //DELETE /authority/delete
	{PType: "p", V0: "admin", V1: "/tags/get", V2: "GET"},                   //GET    /tags/get
	{PType: "p", V0: "admin", V1: "/tags/get", V2: "POST"},                  //POST   /tags/get
	{PType: "p", V0: "admin", V1: "/tags/get", V2: "PUT"},                   //PUT    /tags/get
	{PType: "p", V0: "admin", V1: "/tags/delete", V2: "DELETE"},             //DELETE /tags/delete
	{PType: "p", V0: "admin", V1: "/articles/getOne", V2: "GET"},            //GET    /articles/getOne
	{PType: "p", V0: "admin", V1: "/articles/getAll", V2: "GET"},            //GET    /articles/getAll
	{PType: "p", V0: "admin", V1: "/articles/add", V2: "POST"},              //POST   /articles/add
	{PType: "p", V0: "admin", V1: "/articles/put", V2: "PUT"},               //PUT    /articles/put
	{PType: "p", V0: "admin", V1: "/articles/delete", V2: "DELETE"},         //DELETE /articles/delete
	{PType: "p", V0: "admin", V1: "/cmd/set", V2: "POST"},                   //POST /articles/delete
	{PType: "p", V0: "admin", V1: "/articles/add", V2: "POST"},              //POST /articles/add
	{PType: "p", V0: "admin", V1: "/tags/export", V2: "POST"},               //POST /tags/export
	{PType: "p", V0: "admin", V1: "/tags/import", V2: "POST"},               //POST /tags/import
	{PType: "p", V0: "admin", V1: "/tags/add", V2: "POST"},                  //POST /Tags/add
	{PType: "p", V0: "admin", V1: "/tags/all", V2: "GET"},                   //GET    /tags/all
	{PType: "p", V0: "admin", V1: " /articles/poster/generate", V2: "POST"}, //POST  /articles/poster/generate

}

var casbinrule_custom = []gormadapter.CasbinRule{
	{PType: "p", V0: "custom", V1: "/swagger/*any", V2: "GET"},    //GET    /swagger/*any
	{PType: "p", V0: "custom", V1: "/ws", V2: "GET"},              //GET    /ws
	{PType: "p", V0: "custom", V1: "/alive", V2: "GET"},           //GET    /alive
	{PType: "p", V0: "custom", V1: "/user/test", V2: "GET"},       //GET    /user/test
	{PType: "p", V0: "custom", V1: "/user/:name", V2: "GET"},      //GET    /user/:name
	{PType: "p", V0: "custom", V1: "/user", V2: "GET"},            //GET    /user
	{PType: "p", V0: "custom", V1: "/user/get", V2: "GET"},        //GET    /user/get
	{PType: "p", V0: "custom", V1: "/policy", V2: "GET"},          //GET    /policy
	{PType: "p", V0: "custom", V1: "/product", V2: "GET"},         //GET    /product
	{PType: "p", V0: "custom", V1: "/tags/get", V2: "GET"},        //GET    /tags/get
	{PType: "p", V0: "custom", V1: "/articles/getOne", V2: "GET"}, //GET    /articles/getOne
	{PType: "p", V0: "custom", V1: "/articles/getAll", V2: "GET"}, //GET    /articles/getAll
}

//casbin_rule 表数据初始化
func Casbin_Init() error {
	err := db.AutoMigrate(gormadapter.CasbinRule{}).Error
	if err != nil {
		logrus.Info("err = ", err)
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		defer func() { //处理崩溃
			err := recover()
			logrus.Info("recover ", err)
		}()

		if tx.Find(&[]gormadapter.CasbinRule{}).RowsAffected == 154 {
			logrus.Error("casbin_rule 表的初始数据已存在!")
			return nil
		}

		for _, v := range casbinrule {
			if err := tx.Create(&v).Error; err != nil { //遇到错误时回滚事务
				logrus.Error("create err", err)
				return err
			}
		}

		for _, v := range casbinrule_custom {
			if err := tx.Create(&v).Error; err != nil {
				logrus.Error("create casbinrule_custom err", err)
				return err
			}
		}

		logrus.Info("casbin_rule 表初始数据成功!")
		return nil
	})
}
