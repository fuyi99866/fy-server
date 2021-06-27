package v1

import (
	"github.com/gin-gonic/gin"
	"go_server/models"
	"go_server/utils"
	"net/http"
)

// @Tags 系統
// @Summary 发送测试邮件
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {string} string "{"code":200,"data":{},"msg":"ok"}"
// @Failure 400 {object} utils.Response
// @Router /email [post]
func EmailTest(c *gin.Context) {
	appG := utils.Gin{C: c}
	if err := models.EmailTest(); err != nil {
		appG.Response(http.StatusInternalServerError, utils.ERROR, nil)
	} else {
		appG.Response(http.StatusOK, utils.SUCCESS, nil)
	}
}


