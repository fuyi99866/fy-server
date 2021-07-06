package v1

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"go_server/models"
	"go_server/pkg/app"
	"go_server/pkg/e"
	"go_server/pkg/setting"
	"go_server/pkg/util"
	"net/http"
)

//获取多个文章的标签
func GetTags(c *gin.Context) {
	appG:=app.Gin{C:c}
	//c.Query可用于获取?name=test&state=1这类URL参数
	name := c.Query("name")
	maps := make(map[string]interface{})
	data := make(map[string]interface{})
	if name != "" {
		maps["name"] = name
	}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt() //String转换成int
		maps["state"] = state
	}
	data["lists"]=models.GetTags(util.GetPage(c),setting.AppSetting.PAGE_SIZE,maps)
	data["total"]=models.GetTagTotal(maps)

	appG.Response(http.StatusOK, e.SUCCESS, data)
}

func AddTag(c *gin.Context) {
	appG:=app.Gin{C:c}

	name := c.Query("name")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createdBy := c.Query("created_by")
	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	if !valid.HasErrors(){
		models.AddTag(name,state,createdBy)
		appG.Response(http.StatusOK, e.SUCCESS, nil)
	}else {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
	}
}

func EditTag(c *gin.Context) {
	appG:=app.Gin{C:c}

	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")
	valid := validation.Validation{}
	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Required(id, "id").Message("ID不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")

	if !valid.HasErrors(){
		if models.ExistTagByID(id){
			data:=make(map[string]interface{})
			data["modified_at"] = modifiedBy
			if name != ""{
				data["name"] = name
			}
			if state != -1{
				data["state"] = state
			}
			models.EditTag(id,data)
			appG.Response(http.StatusOK, e.SUCCESS, data)
		}else {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		}
	}
}

func DeleteTag(c *gin.Context) {
	appG:=app.Gin{C:c}

	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if ! valid.HasErrors() {

		if models.ExistTagByID(id) {
			models.DeleteTag(id)
			appG.Response(http.StatusOK, e.SUCCESS, nil)
		} else {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		}
	}
}
