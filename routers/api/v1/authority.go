package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go_server/models"
	"go_server/pkg/e"

	//"go_server/pgo kg/e"
	"go_server/pkg/app"
	"net/http"
)

// @Summary   创建角色
// @Tags   角色
// @Accept json
// @Produce  json
// @Param   body  body   models.Authority   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /authority/add  [POST]
// @Security ApiKeyAuth
func CreateAuthority(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo models.Authority
	err := c.ShouldBindJSON(&reqInfo)
	if err != nil {
		logrus.Info("CreateAuthority param error")
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, err.Error())
		return
	}
	/*
		menu := map[string]interface{}{
			"id":       reqInfo.AuthorityId,
			"name":     reqInfo.AuthorityName,
			"parentId": reqInfo.ParentId,
		}*/

	if _,err := models.CreateAuthority(reqInfo); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		logrus.Error("CreateAuthority error")
		return
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
	}
}

// @Tags 角色
// @Summary 删除角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body models.Authority true "删除角色"
// @Success 200 {string} string "{"code":200,"data":{},"msg":"ok"}"
// @Router /authority/delete [delete]
func DeleteAuthority(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo models.Authority
	err := c.ShouldBindJSON(&reqInfo)
	if err != nil {
		logrus.Info("CreateAuthority param error")
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, err.Error())
		return
	}
	if err := models.DeleteAuthority(reqInfo); err != nil { // 删除角色之前需要判断是否有用户正在使用此角色
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		logrus.Error("DeleteAuthority error")
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
	}
}

// @Tags 角色
// @Summary 更新角色信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body models.Authority true "权限id, 权限名, 父角色id"
// @Success 200 {string} string "{"code":200,"data":{},"msg":"ok"}"
// @Router /authority/update [post]
func UpdateAuthority(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo models.Authority
	err := c.ShouldBindJSON(&reqInfo)
	if err != nil {
		logrus.Info("CreateAuthority param error")
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, err.Error())
		return
	}

	if err, authority := models.UpdateAuthority(reqInfo); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		logrus.Error("UpdateAuthority error")
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, authority)
	}
}

// @Tags 角色
// @Summary 设置角色资源权限
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body models.Authority true "设置角色资源权限"
// @Success 200 {string} string "{"code":200,"data":{},"msg":"ok"}"
// @Router /authority/set [post]
func SetAuthority(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo models.Authority
	err := c.ShouldBindJSON(&reqInfo)
	if err != nil {
		logrus.Info("CreateAuthority param error")
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, err.Error())
	}

	if err := models.SetAuthority(reqInfo); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		logrus.Error("SetDataAuthority error")
		return
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
	}
}



