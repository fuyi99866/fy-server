package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go_server/models"
	"go_server/pkg/app"
	"go_server/pkg/e"
	"go_server/pkg/util"
	"net/http"
	"strconv"
)

/**
对用户数据进行增删查改操作
*/

// @Summary   获取所有用户
// @Tags   用户
// @Accept json
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /user  [GET]
// @Security ApiKeyAuth
func GetUsers(c *gin.Context) {
	//分页显示所有用户
	appG := app.Gin{C: c}
	user, err := models.GetAllUser()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	data := make(map[string]interface{})
	data["lists"] = user
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// @Summary   增加用户
// @Tags   用户
// @Accept json
// @Produce  json
// @Param   body  body   models.UserRegister   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /user  [POST]
// @Security ApiKeyAuth
func AddUser(c *gin.Context) {
	logrus.Info("AddUser")
	appG := app.Gin{C: c}
	var reqInfo models.UserRegister
	err := c.ShouldBindJSON(&reqInfo)
	if err != nil {
		logrus.Info("AddUser param error")
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, err.Error())
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
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		logrus.Error("AddUser error")
		return
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
	}
}

// @Summary   更新用户信息
// @Tags   用户
// @Accept json
// @Produce  json
// @Param   body  body   models.UserRegister   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /user/update  [POST]
// @Security ApiKeyAuth
func UpdateUser(c *gin.Context) {
	logrus.Info("UpdateUser")
	appG := app.Gin{C: c}
	var reqInfo models.UserRegister
	err := c.ShouldBindJSON(&reqInfo)
	if err != nil {
		logrus.Info("AddUser param error")
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, err.Error())
		return
	}

	_err := models.UpdateUser(reqInfo)
	if _err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		logrus.Info("UpdateUser error")
		return
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
	}
}

// @Summary   获取单一用户
// @Tags   用户
// @Accept json
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /user/get  [GET]
// @Security ApiKeyAuthc
func GetOneUser(c *gin.Context) {
	appG := app.Gin{C: c}
	Authorization := c.GetHeader("Authorization")
	claims, err := util.ParseToken(Authorization)
	if err != nil {
		appG.Response(http.StatusForbidden, e.ACCESS_DENIED, nil)
		return
	}

	user, err := models.GetOneUser(claims.Username)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, user)
}

// @Summary   刪除用户
// @Tags   用户
// @Accept json
// @Produce  json
// @Param   body  body   models.UserRegister   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /user/delete  [delete]
// @Security ApiKeyAuth
func DeleteUser(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo models.UserRegister
	err := c.ShouldBindJSON(&reqInfo)
	if err != nil {
		logrus.Info("AddUser param error")
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, err.Error())
		return
	}

	_err := models.DeleteUser(reqInfo)
	if _err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		logrus.Error("DeleteUser error")
		return
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
	}
}

// @Summary 设置用户权限
// @Tags 用户
// @Accept json
// @Produce  json
// @Param   body  body   models.SetUserAuth   true "id:用户ID, auth_id:角色ID"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /user/set_auth [post]
// @Security ApiKeyAuth
func SetUserAuthority(c *gin.Context) {
	appG := app.Gin{C: c}
	var sua models.SetUserAuth
	_ = c.ShouldBindJSON(&sua)
	if UserVerifyErr := util.Verify(sua, util.SetUserAuthorityVerify); UserVerifyErr != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, UserVerifyErr.Error())
		return
	}
	id, _ := strconv.Atoi(sua.ID)
	if err := models.SetUserAuthority(id, sua.AuthorityId); err != nil {

		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
	}
}
