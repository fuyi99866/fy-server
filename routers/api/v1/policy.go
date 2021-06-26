package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	alog "go_server/log"
	"go_server/models"
	"go_server/routers/casbin/enforcer"
	"go_server/utils"
	"net/http"
)

// @Summary   增加访问权限
// @Tags   访问权限
// @Accept json
// @Produce  json
// @Param   body  body   models.UserPolicy   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} utils.Response
// @Router /policy  [POST]
// @Security ApiKeyAuth
func AddPolicy(c *gin.Context) {
	alog.MyLogger.Info("1增加Policy")
	appG := utils.Gin{C: c}
	var reqInfo models.UserPolicy
	err := c.ShouldBindJSON(&reqInfo)
	if err != nil {
		alog.MyLogger.Error("AddPolicy param error")
		appG.Response(http.StatusBadRequest, utils.INVALID_PARAMS, err.Error())
		return
	}

	e := enforcer.EnforcerTool()
	fmt.Println("增加Policy")
	if ok := e.AddPolicy(reqInfo.Username,reqInfo.URL,reqInfo.Type); !ok {
		fmt.Println("Policy已经存在")
		appG.Response(http.StatusInternalServerError, utils.ERROR, nil)
	} else {
		fmt.Println("增加成功")
		appG.Response(http.StatusOK, utils.SUCCESS, nil)
	}
}

// @Summary   删除访问权限
// @Tags   访问权限
// @Accept json
// @Produce  json
// @Param   body  body   models.UserPolicy   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} utils.Response
// @Router /policy  [DELETE]
// @Security ApiKeyAuth
func DeletePolicy(c *gin.Context) {
	logrus.Info("删除Policy")
	appG := utils.Gin{C: c}
	var reqInfo models.UserPolicy
	err := c.ShouldBindJSON(&reqInfo)
	if err != nil {
		alog.MyLogger.Error("AddPolicy param error")
		appG.Response(http.StatusBadRequest, utils.INVALID_PARAMS, err.Error())
		return
	}

	e := enforcer.EnforcerTool()
	fmt.Println("删除Policy")
	if ok := e.RemovePolicy(reqInfo.Username,reqInfo.URL,reqInfo.Type); !ok {
		fmt.Println("Policy不存在")
		appG.Response(http.StatusInternalServerError, utils.ERROR, nil)
	} else {
		fmt.Println("删除成功")
		appG.Response(http.StatusOK, utils.SUCCESS, nil)
	}
}

// @Summary   获取权限列表
// @Tags   访问权限
// @Accept json
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} utils.Response
// @Router /policy  [GET]
// @Security ApiKeyAuth
func GetPolicy(c *gin.Context) {
	logrus.Info("1查看Policy")
	appG := utils.Gin{C: c}
	e := enforcer.EnforcerTool()
	fmt.Println("2查看Policy")

	list := e.GetPolicy()
	for _, vlist := range list {
		for _, v := range vlist {
			fmt.Printf("value: %s, ", v)
		}
	}

	appG.Response(http.StatusOK, utils.SUCCESS, list)
}
