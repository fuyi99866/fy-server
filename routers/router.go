package routers

import (
	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/examples/datamodel"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/plugins/admin"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/themes/adminlte"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go_server/middleware/jwt"
	"go_server/pkg/export"
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
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // API 注释
	r.POST("/auth", api.Auth)                                            // token鉴权
	r.StaticFS("/export",http.Dir(export.GetExcelFullPath()))            //下载导出的excel
	r.POST("/upload_img", api.UploadImage)                               // 上传图片
	r.POST("/upload_file", api.UploadFile)

	//websocket实现聊天室
	r.GET("/ws", websocket.NotifySocket)
	//访问静态前端文件
	r.Static("static", "dist/static")
	r.Static("img", "dist/img")
	r.StaticFile("/", "dist/index.html")

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

func Mcc(context *gin.Context)  {
	r := gin.Default()

	eng := engine.Default()

	// global config
	cfg := config.Config{
		Databases: config.DatabaseList{
			"default": {
				Host:         "127.0.0.1",
				Port:         "3306",
				User:         "root",
				Pwd:          "root",
				Name:         "godmin",
				MaxIdleCon: 50,
				MaxOpenCon: 150,
				Driver:       "mysql",
			},
		},
		UrlPrefix: "admin",
		// STORE 必须设置且保证有写权限，否则增加不了新的管理员用户
		Store: config.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
		Language: language.CN,
		// 开发模式
		Debug: true,
		// 日志文件位置，需为绝对路径
		InfoLogPath: "/var/logs/info.log",
		AccessLogPath: "/var/logs/access.log",
		ErrorLogPath: "/var/logs/error.log",
		ColorScheme: adminlte.ColorschemeSkinBlack,
	}

	// Generators： 详见 https://github.com/GoAdminGroup/go-admin/blob/master/examples/datamodel/tables.go
	adminPlugin := admin.NewAdmin(datamodel.Generators)

	// 增加 chartjs 组件
	template.AddComp(chartjs.NewChart())

	// 增加 generator, 第一个参数是对应的访问路由前缀
	// 例子:
	//
	// "user" => http://localhost:9033/admin/info/user
	//
	// adminPlugin.AddGenerator("user", datamodel.GetUserTable)

	// 自定义首页

	r.GET("/admin", func(ctx *gin.Context) {
		engine.Content(ctx, func(ctx interface{}) (types.Panel, error) {
			return datamodel.GetContent(nil)
		})
	})

	_ = eng.AddConfig(&cfg).AddPlugins(adminPlugin).Use(r)

	//_ = r.Run(":9033")
}
