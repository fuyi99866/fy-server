package v1

import (

	"github.com/gin-gonic/gin"
	"go_server/models"
)

// @Summary   获取所有产品
// @Tags   用户
// @Accept json
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /users  [GET]

func GetProducts(c *gin.Context) {
	models.GetAllRobot()
}
