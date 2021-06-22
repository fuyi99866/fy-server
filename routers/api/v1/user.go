package v1

import (

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go_server/models"
	"go_server/routers/jwt"
	"go_server/utils"
	"net/http"
	"strings"
)

/**
对用户数据进行增删查改操作
*/

// @Summary  检查token是否过期
// @Tags 鉴权
// @Accept json
// @Produce json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} utils.Response
// @Router /alive  [get]
// @Security ApiKeyAuth
func TokenAlive(c *gin.Context) {
	appG := utils.Gin{C: c}
	appG.Response(http.StatusOK, utils.SUCCESS, nil)
}

// @Summary   登录获取登录token 信息
// @Tags 鉴权
// @Accept json
// @Produce  json
// @Param   body  body   models.UserLogin   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} utils.Response
// @Router /auth  [POST]
func Auth(c *gin.Context) {
	appG := utils.Gin{C: c}
	var reqInfo models.UserLogin
	if err := c.ShouldBindJSON(&reqInfo); err != nil {
		body, _ := c.GetRawData()
		logrus.Println("Auth request: ", string(body))
		appG.Response(http.StatusBadRequest, utils.INVALID_PARAMS, nil)
		return
	}
	logrus.Println("Auth: ", reqInfo.Username, reqInfo.Password)

	//有效性验证
	valid := validation.Validation{}
	valid.AlphaNumeric(reqInfo.Username,"username").Message("非法用户名")
	valid.Required(reqInfo.Username,"User").Message("用户名不能为空")
	valid.MinSize(reqInfo.Username,1,"User").Message("用户名最短为1位")
	valid.MaxSize(reqInfo.Password,100,"username").Message("最长为100字符")
	valid.MaxSize(reqInfo.Password,100,"password").Message("最长为100字符")
	valid.Required(reqInfo.Password, "Pwd").Message("密码不能为空")
	valid.MinSize(reqInfo.Password, 6, "Pwd").Message("密码最短为6位")
	//检测用户的账号密码，看能否正常登录
	isOk, err := models.CheckUser(reqInfo.Username, reqInfo.Password)

	if err != nil {
		appG.Response(http.StatusInternalServerError, utils.ERROR, nil)
		return
	}

	if !isOk {
		appG.Response(http.StatusUnauthorized, utils.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	token, err := jwt.GenerateToken(reqInfo.Username, reqInfo.Password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, utils.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, utils.SUCCESS, map[string]string{
		"token": token,
	})
}

//// @Summary   获取当前登录用户信息
//// @Tags  用户
//// @Accept json
//// @Produce  json
//// @Success 200 {string} json "{ "code": 200, "data": {"lists":""}, "msg": "ok" }"
//// @Failure 400 {object} app.Response
//// @Router /users/current_user  [GET]
//// @Security ApiKeyAuth
//func GetCurrentUserInfo(c *gin.Context) {
//	appG := app.Gin{C: c}
//	Authorization := c.GetHeader("Authorization")
//	t, err := jwt.Parse(Authorization, func(*jwt.Token) (interface{}, error) {
//		return jwtGet.JwtSecret, nil
//	})
//
//	if err != nil {
//		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH, nil)
//		return
//	}
//	u := jwtGet.GetIdFromClaims("username", t.Claims)
//	userService := user_service.User{
//		Username: u,
//		PageNum:  util.GetPage(c),
//		PageSize: conf.AppSetting.PageSize,
//	}
//
//	user, err := userService.Get()//get user by username
//
//	if err != nil {
//		appG.Response(http.StatusInternalServerError, e.ERROR_GET_S_FAIL, nil)
//		return
//	}
//
//	user.Password = ""
//	//user.RobotList = ""
//	data := make(map[string]interface{})
//	data["lists"] = user
//
//	appG.Response(http.StatusOK, e.SUCCESS, data)
//}

// @Summary   获取所有用户
// @Tags   用户
// @Accept json
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} utils.Response
// @Router /user  [GET]
// @Security ApiKeyAuth
func GetUsers(c *gin.Context) {
	appG := utils.Gin{C: c}
	user, err := models.GetAllUser()
	if err != nil {
		appG.Response(http.StatusInternalServerError, utils.ERROR, nil)
		return
	}
	data := make(map[string]interface{})
	data["lists"] = user
	appG.Response(http.StatusOK, utils.SUCCESS, data)
}

// @Summary   增加用户
// @Tags   用户
// @Accept json
// @Produce  json
// @Param   body  body   models.UserRegister   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} utils.Response
// @Router /user  [POST]
// @Security ApiKeyAuth
func AddUser(c *gin.Context) {
	logrus.Info("AddUser")
	appG := utils.Gin{C: c}
	var reqInfo models.UserRegister
	err := c.ShouldBindJSON(&reqInfo)
	if err != nil {
		logrus.Infoln("AddUser param error")
		appG.Response(http.StatusBadRequest, utils.INVALID_PARAMS, err.Error())
		return
	}
	menu := map[string]interface{}{
		"username":   reqInfo.Username,
		"password":   reqInfo.Password,
		"nickname":   reqInfo.NickName,
		"company_id": reqInfo.CompanyID,
	}

	_, _err := models.AddUser(menu)
	if _err != nil {
		appG.Response(http.StatusInternalServerError, utils.ERROR, nil)
		return
	} else {
		appG.Response(http.StatusOK, utils.SUCCESS, nil)
	}
}

//Param方法来获取API参数
func GetApiParam(context *gin.Context) {
	name := context.Param("name")
	action := context.Param("action")
	//截取
	action = strings.Trim(action, "/")

	context.String(http.StatusOK, name+" is "+action)
}