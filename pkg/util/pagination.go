package util

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"go_server/pkg/setting"
)

/**
分页页码的获取方法
*/
func GetPage(c *gin.Context) int {
	result := 0
	page, _ := com.StrTo(c.Query("page")).Int()
	if page > 0 {
		result = (page - 1) * setting.AppSetting.PageSize
	}
	return result
}
