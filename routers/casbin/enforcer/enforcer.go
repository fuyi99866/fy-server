package enforcer

import (
	"fmt"
	"github.com/casbin/casbin"
	"github.com/casbin/gorm-adapter"
	"github.com/gin-gonic/gin"
	alog "go_server/log"
	"go_server/routers/casbin/DB"
	"go_server/routers/jwt"
	"go_server/utils"
	"net/http"
)

//拦截器函数
func Interceptor(e *casbin.Enforcer) gin.HandlerFunc {
	return func(context *gin.Context) {
		appG := utils.Gin{C: context}
		Authorization := context.GetHeader("Authorization")
		claims, err := jwt.ParseToken(Authorization)
		if err!=nil{
			appG.Response(http.StatusForbidden, utils.ACCESS_DENIED, nil)
			return
		}
		alog.MyLogger.Info("claims===",claims)
        fmt.Println("claims===",claims)
		//获取请求的URI
		obj := context.Request.URL.RequestURI()
		//获取请求方法
		act := context.Request.Method
		//获取用户的角色
		sub := claims.Username

		//判断策略中是否存在
		if ok := e.Enforce(sub, obj, act); ok {
			alog.MyLogger.Info("通过权限")
			context.Next()
		} else {
			alog.MyLogger.Warn("没有通过权限")
			context.Abort()
			appG.Response(http.StatusForbidden, utils.ACCESS_DENIED, nil)
		}
	}
}

/**
决策器
*/

func EnforcerTool() *casbin.Enforcer {
	//适配器,自动生成casbin_rule表，存放访问的路径和方法等信息
	adapter := gormadapter.NewAdapterByDB(DB.MysqlTool())
	//创建一个casbin决策器需要一个模板文件和策略文件为参数
	Enforcer := casbin.NewEnforcer("conf/keymatch.conf", adapter)
	//从数据库加载策略
	Enforcer.LoadPolicy()
	return Enforcer
}
