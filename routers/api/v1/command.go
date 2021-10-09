package v1

import (
	"github.com/gin-gonic/gin"
	"go_server/models"
	"go_server/pkg/app"
	"go_server/pkg/e"
	"go_server/pkg/logger"
	"net/http"
)

// @Summary   控制指令
// @Tags   控制
// @Accept json
// @Produce  json
// @Param   body  body   models.Command   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /cmd/set [POST]
// @Security ApiKeyAuth
func SetRobot(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo models.Command
	err := c.ShouldBindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	//TODO 给机器人发送指令 "start" ,"stop","pause","resume"
	logger.Info("给机器人发送指令 'start' ,'stop','pause','resume'")
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		logger.Info("send cmd failed: ", err)
		return
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
	}
}
