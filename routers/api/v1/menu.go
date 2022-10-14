package v1

import (
	"github.com/gin-gonic/gin"
	"go_server/models"
	"go_server/pkg/app"
	"go_server/pkg/e"
	"go_server/pkg/util"
	"net/http"
	"strings"
)

// @Summary   获取用户动态路由
// @Tags   菜单
// @Accept json
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /menu  [GET]
// @Security ApiKeyAuth
func GetMenu(c *gin.Context) {
	appG:=app.Gin{C: c}
	Authorization := c.GetHeader("Authorization")
	claims, err := util.ParseToken(Authorization)
	if err != nil {
		appG.Response(http.StatusForbidden, e.ACCESS_DENIED, nil)
		return
	}

	if err, menus := models.GetMenuTreeMap(claims.Id); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, menus)
	}
}


// @Summary   查询所有菜单
// @Tags   菜单
// @Accept json
// @Produce  json
// @Success 200 {object} app.Response{data=models.MenuResult}
// @Failure 400 {object} app.Response
// @Router /menu [GET]
// @Security ApiKeyAuth
func GetAllMenus(c *gin.Context) {
	appG := app.Gin{C: c}
	menus, err := models.GetAllMenus()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, menus)
}

// @Summary   添加菜单
// @Tags   菜单
// @Accept json
// @Produce  json
// @Param   body  body   models.MenuAdd   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /menu [POST]
// @Security ApiKeyAuth
func AddMenu(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo models.MenuAdd
	err := c.ShouldBindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, err.Error())
		return
	}
	menu := models.Menu{
		Name:      reqInfo.Name,
		Path:      reqInfo.Path,
		Component: reqInfo.Component,
		Url:       reqInfo.Url,
		Status:    reqInfo.Status,
	}
	_, err = models.AddMenu(menu)

	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary   删除菜单
// @Tags   菜单
// @Accept json
// @Produce  json
// @Param menuname query string true "菜单名"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /menu [DELETE]
// @Security ApiKeyAuth
func DeleteMenu(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.Query("menuname")
	name = strings.Replace(name, " ", "", -1)
	err := models.DeleteMenuByName(name)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}


