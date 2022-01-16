package routers

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go_server/conf"
	"go_server/middleware/jwt"
	"go_server/pkg/app"
	"go_server/pkg/e"
	"go_server/pkg/export"
	"go_server/pkg/logger"
	"go_server/pkg/setting"
	"go_server/pkg/upload"
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

	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	//将访问路由到swagger的HTML页面
	r.GET("/version", func(context *gin.Context) {
		context.JSON(http.StatusOK,conf.AppVersion)
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // API 注释
	r.POST("/auth", api.Auth)                                            // token鉴权
	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))           //下载导出的excel
	r.POST("/upload_img", api.UploadImage)                               // 上传图片
	r.POST("/upload_file", api.UploadFile)

	r.POST("/ifaas-hotel-robot-platform/api/hotel/tech/robot/disinfect/request/create/1.0", func(context *gin.Context) {
		buf := make([]byte, 1024)
		n, _ := context.Request.Body.Read(buf)
		logger.Info(string(buf[0:n]))

		appG := app.Gin{C: context}
		data := make(map[string]interface{})
		data["taskKey"] = "123"
		data["pointId"] = "2"
		data["vendor"] = "ubtech"
		appG.ResponseTest(http.StatusOK, data, e.SUCCESS)

	})

	r.POST("/ifaas-hotel-robot-platform/api/hotel/tech/robot/disinfect/depository/notify/1.0", func(context *gin.Context) {
		buf := make([]byte, 1024)
		n, _ := context.Request.Body.Read(buf)
		logger.Info(string(buf[0:n]))
		appG := app.Gin{C: context}
		data := make(map[string]interface{})
		data["pointId"] = "A2-3F-1K"
		appG.ResponseTest(http.StatusOK, data, e.SUCCESS)

	})

	r.POST("/ifaas-hotel-robot-platform/api/hotel/tech/robot/disinfect/robot/notify/1.0", func(context *gin.Context) {
		buf := make([]byte, 1024)
		n, _ := context.Request.Body.Read(buf)
		logger.Info(string(buf[0:n]))
		appG := app.Gin{C: context}
		appG.ResponseTest(http.StatusOK, nil, e.SUCCESS)

	})

	r.POST("/ifaas-hotel-robot-platform/api/hotel/tech/robot/status/report/1.0", func(context *gin.Context) {
		buf := make([]byte, 1024)
		n, _ := context.Request.Body.Read(buf)
		logger.Info(string(buf[0:n]))
		appG := app.Gin{C: context}
		appG.ResponseTest(http.StatusOK, nil, e.SUCCESS)

	})

	r.POST("/ifaas-authority/oauth/token", func(context *gin.Context) {
		buf := make([]byte, 1024)
		n, _ := context.Request.Body.Read(buf)
		logger.Info(string(buf[0:n]))
		//appG := app.Gin{C: context}
		data := make(map[string]interface{})
		data["access_token"] = "GJbqRK_cWlLTkhzWfKM/dXnwX58BQ=="
		data["token_type"] = "bearer"
		data["expires_in"] = 8639
		data["scope"] = "read write"

		context.JSON(http.StatusOK, gin.H{
			"respCode":     10000000,
			"respMessage":  "操作成功！",
			"access_token": "GJbqRK_cWlLTkhzWfKM/dXnwX58BQ==",
			"token_type":   "bearer",
			"expires_in":   8639,
			"scope":        "read write",
		})

	})

	//访问静态前端文件
	r.Static("static", "dist/static")
	r.Static("img", "dist/static/img")
	r.StaticFile("/", "dist/index.html")

	r.GET("/profile", api.GetConnectProfile)
	r.GET("/channel", websocket.NotifySocket)
	//group1 := r.Group("api/v1")
	group1 := r.Group("")
	group1.Use(jwt.JWT())                                     //token 验证
	group1.Use(enforcer.Interceptor(enforcer.EnforcerTool())) //拦截器进行访问控制
	{
		group1.GET("/alive", api.TokenAlive)
		user := group1.Group("user")
		{
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
			tags.GET("all", v1.GetTags)
			tags.POST("add", v1.AddTag)
			tags.PUT("get", v1.EditTag)
			tags.DELETE("delete", v1.DeleteTag)
			tags.POST("export", v1.ExportTag)
			tags.POST("import", v1.ImportTag)

		}

		articles := group1.Group("articles")
		{
			articles.GET("/:id", v1.GetArticle)
			articles.GET("all", v1.GetArticles)
			articles.POST("add", v1.AddArticle)
			articles.PUT("put", v1.EditArticle)
			articles.DELETE("delete", v1.DeleteArticle)
			articles.POST("poster/generate", v1.GenerateArticlePoster)
		}
		cmd := group1.Group("cmd")
		{
			cmd.POST("set", v1.SetRobot)
		}
	}

	return r
}
