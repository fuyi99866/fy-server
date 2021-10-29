package v1

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"go_server/models"
	"go_server/pkg/app"
	"go_server/pkg/e"
	"go_server/pkg/export"
	"go_server/pkg/setting"
	"go_server/pkg/util"
	"go_server/service/tag_service"
	"net/http"
)

// @Summary   查询多个标签
// @Tags   标签
// @Accept json
// @Produce  json
// @Param  name  query  string  false "Name"
// @Param  state  query  int  false "state"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /tags/all  [GET]
// @Security ApiKeyAuth
func GetTags(c *gin.Context) {
	appG := app.Gin{C: c}
	//c.Query可用于获取?name=test&state=1这类URL参数
	name := c.Query("name")

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt() //String转换成int
	}

	tag := tag_service.Tag{
		Name:     name,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}
	tags, err := tag.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_TAG_FAIL, nil)
		return
	}

	count, err := tag.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"lists": tags,
		"total": count,
	})
}

// @Summary   添加标签
// @Tags   标签
// @Accept json
// @Produce  json
// @Param  body  body  models.AddTagForm  true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /tags/add  [POST]
// @Security ApiKeyAuth
func AddTag(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo models.AddTagForm
	err := c.ShouldBindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	valid := validation.Validation{}
	valid.Required(reqInfo.Name, "name").Message("名称不能为空")
	valid.MaxSize(reqInfo.Name, 100, "name").Message("名称最长为100字符")
	valid.Required(reqInfo.CreatedBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(reqInfo.CreatedBy, 100, "created_by").Message("创建人最长为100字符")
	valid.Range(reqInfo.State, 0, 1, "state").Message("状态只允许0或1")

	if valid.HasErrors() {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	tag := tag_service.Tag{
		Name:      reqInfo.Name,
		CreatedBy: reqInfo.CreatedBy,
		State:     reqInfo.State,
	}

	exsit, err := tag.ExistByName()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_TAG_FAIL, nil)
		return
	}
	if exsit {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG, nil)
		return
	}
	err = tag.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func EditTag(c *gin.Context) {
	appG := app.Gin{C: c}

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

	if !valid.HasErrors() {
		exist, err := models.ExistTagByID(id)
		if err != nil {
			appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
			return
		}
		if exist {
			data := make(map[string]interface{})
			data["modified_at"] = modifiedBy
			if name != "" {
				data["name"] = name
			}
			if state != -1 {
				data["state"] = state
			}
			models.EditTag(id, data)
			appG.Response(http.StatusOK, e.SUCCESS, data)
		} else {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		}
	}
}

func DeleteTag(c *gin.Context) {
	appG := app.Gin{C: c}

	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if ! valid.HasErrors() {
		exist, err := models.ExistTagByID(id)
		if err != nil {
			appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
			return
		}
		if exist {
			models.DeleteTag(id)
			appG.Response(http.StatusOK, e.SUCCESS, nil)
		} else {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		}
	}
}

// @Summary   导出标签
// @Tags   标签
// @Accept json
// @Produce  json
// @Param  body  body  models.ExportTagForm  false "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /tags/export  [POST]
// @Security ApiKeyAuth
func ExportTag(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo models.ExportTagForm
	err := c.ShouldBindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	tag := tag_service.Tag{
		Name:  reqInfo.Name,
		State: reqInfo.State,
	}

	filename, err := tag.Export()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXPORT_TAG_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, http.StatusOK, map[string]interface{}{
		"export_url":      export.GetExcelFullUrl(filename),
		"export_save_url": export.GetExcelPath() + filename,
	})
}

// @Summary   导入标签
// @Tags   标签
// @Accept multipart/form-data
// @Param file formData file true "file"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /tags/import  [POST]
// @Security ApiKeyAuth
func ImportTag(c *gin.Context) {
	appG := app.Gin{C: c}
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	tag := tag_service.Tag{}
	err = tag.Import(file)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_IMPORT_TAG_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
