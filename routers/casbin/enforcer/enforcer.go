package enforcer

import (
	"github.com/casbin/casbin"
	"github.com/casbin/gorm-adapter"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go_server/pkg/app"
	e2 "go_server/pkg/e"
	"go_server/routers/casbin/DB"
	"go_server/routers/jwt"
	"net/http"
)

//拦截器函数
func Interceptor(e *casbin.Enforcer) gin.HandlerFunc {
	return func(context *gin.Context) {
		appG := app.Gin{C: context}
		Authorization := context.GetHeader("Authorization")
		claims, err := jwt.ParseToken(Authorization)
		if err!=nil{
			appG.Response(http.StatusForbidden, e2.ACCESS_DENIED, nil)
			return
		}
        logrus.Info("claims===",claims)
		//获取请求的URI
		obj := context.Request.URL.RequestURI()
		//获取请求方法
		act := context.Request.Method
		//获取用户的角色
		sub := claims.Username

		//判断策略中是否存在
		if ok := e.Enforce(sub, obj, act); ok {
			logrus.Info("通过权限")
			context.Next()
		} else {
			logrus.Warn("没有通过权限")
			context.Abort()
			appG.Response(http.StatusForbidden, e2.ACCESS_DENIED, nil)
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
