package api

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"go_server/models"
	"go_server/pkg/app"
	"go_server/pkg/e"
	"go_server/pkg/logger"
	"go_server/pkg/util"
	"net/http"
)

// @Summary  检查token是否过期
// @Tags 鉴权
// @Accept json
// @Produce json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /alive  [get]
// @Security ApiKeyAuth
func TokenAlive(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary   登录获取登录token 信息
// @Tags 鉴权
// @Accept json
// @Produce  json
// @Param   body  body   models.UserLogin   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /auth  [POST]
func Auth(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo models.UserLogin
	if err := c.ShouldBindJSON(&reqInfo); err != nil {
		body, _ := c.GetRawData()
		logger.Info("Auth request: ", string(body))
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	logger.Info("Auth: ", reqInfo.Username, reqInfo.Password)

	//有效性验证
	valid := validation.Validation{}
	valid.AlphaNumeric(reqInfo.Username, "username").Message("非法用户名")
	valid.Required(reqInfo.Username, "User").Message("用户名不能为空")
	valid.MinSize(reqInfo.Username, 1, "User").Message("用户名最短为1位")
	valid.MaxSize(reqInfo.Password, 100, "username").Message("最长为100字符")
	valid.MaxSize(reqInfo.Password, 100, "password").Message("最长为100字符")
	valid.Required(reqInfo.Password, "Pwd").Message("密码不能为空")
	valid.MinSize(reqInfo.Password, 6, "Pwd").Message("密码最短为6位")
	//检测用户的账号密码，看能否正常登录
	isOk, err := models.CheckUser(reqInfo.Username, reqInfo.Password)

	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	if !isOk {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	token, err := util.GenerateToken(reqInfo.Username, reqInfo.Password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
}

