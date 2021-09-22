package routers

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go_server/middleware/jwt"
	"go_server/pkg/setting"
	"go_server/routers/api"
	v1 "go_server/routers/api/v1"
	"go_server/routers/casbin/enforcer"
	"go_server/routers/websocket"
	"net/http"
)

//swag init --generalInfo .\routers\router.go

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery()) //崩溃恢复，直接返回500
	gin.SetMode(setting.RunMode)

	//2、绑定路由规则，执行的函数
	//gin.Context,封装了request和response

	//将访问路由到swagger的HTML页面
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // API 注释
	r.POST("/auth", v1.Auth)                                             // token鉴权
	r.POST("/upload", api.UploadImage)                                   //上传图片

	//websocket实现聊天室
	r.GET("/ws", websocket.NotifySocket)
	http.Handle("/",http.FileServer(http.Dir("dist")))
	//group1 := r.Group("api/v1")
	group1 := r.Group("")
	group1.Use(jwt.JWT())                                     //token 验证
	group1.Use(enforcer.Interceptor(enforcer.EnforcerTool())) //拦截器进行访问控制
	{
		group1.GET("/alive", v1.TokenAlive)
		user := group1.Group("user")
		{
			user.GET("/test", Response_test) //测试回复
			user.GET("/:name", v1.GetApiParam)
			user.GET("", v1.GetUsers)
			user.POST("", v1.AddUser)
			user.DELETE("delete", v1.DeleteUser)
			user.POST("update", v1.UpdateUser)
			user.GET("get", v1.GetOneUser)
		}

		policy := group1.Group("policy")
		{
			policy.POST("", v1.AddPolicy)
			policy.DELETE("", v1.DeletePolicy)
			policy.GET("", v1.GetPolicy)
		}

		product := group1.Group("product")
		{
			product.GET("", v1.GetProducts)
		}

		authority := group1.Group("authority")
		{
			authority.POST("add", v1.CreateAuthority)
			authority.POST("update", v1.UpdateAuthority)
			authority.POST("set", v1.SetAuthority)
			authority.DELETE("delete", v1.DeleteAuthority)
		}

		tags := group1.Group("tags")
		{
			tags.GET("get", v1.GetTags)
			tags.POST("get", v1.AddTag)
			tags.PUT("get", v1.EditTag)
			tags.DELETE("delete", v1.DeleteTag)
		}

		articles := group1.Group("articles")
		{
			articles.GET("getOne", v1.GetArticle)
			articles.GET("getAll", v1.GetArticles)
			articles.POST("add", v1.AddArticle)
			articles.PUT("put", v1.EditArticle)
			articles.DELETE("delete", v1.DeleteArticle)
		}
	}

	return r
}

//简单的回复成功
func Response_test(context *gin.Context) {
	message := "成功"
	code := 200
	data := "data"
	context.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
		"result":  "true",
	})
}
