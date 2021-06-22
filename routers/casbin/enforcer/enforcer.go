package enforcer
import (

	"fmt"
	"github.com/casbin/casbin"
	"github.com/casbin/gorm-adapter"
	"github.com/gin-gonic/gin"
	"go_server/routers/casbin/DB"
)

//拦截器函数
func Interceptor(e *casbin.Enforcer) gin.HandlerFunc {
	return func(context *gin.Context) {
		//获取资源
		obj := context.Request.URL.RequestURI()
		//获取方法
		act := context.Request.Method
		//获取实体
		sub := "admin"

		//判断策略中是否存在
		enforce := e.AddPolicy(sub, obj, act)
		if enforce {
			fmt.Println("通过权限")
			context.Next()
		} else {
			fmt.Println("没有通过权限")
			context.Abort()
		}
	}
}

/**
决策器
*/

func EnforcerTool() *casbin.Enforcer {
	//适配器,自动生成casbin_rule表，存放访问的路径和方法等信息
	adapter:= gormadapter.NewAdapterByDB(DB.MysqlTool())
	//创建一个casbin决策器需要一个模板文件和策略文件为参数
	Enforcer := casbin.NewEnforcer("conf/keymatch.conf",adapter)
	//从数据库加载策略
	Enforcer.LoadPolicy()
	return Enforcer
}



