package v1

import (
	"github.com/astaxie/beego/validation"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"go_server/models"
	"go_server/pkg/app"
	"go_server/pkg/e"
	"go_server/pkg/logger"
	"go_server/pkg/qrcode"
	"go_server/pkg/setting"
	"go_server/pkg/util"
	"go_server/service/article_service"
	"go_server/service/tag_service"
	"net/http"
)

// @Summary   查询文章
// @Tags   文章
// @Accept json
// @Produce  json
// @Param  id  path  int  true "文章ID"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /articles/{id}  [GET]
// @Security ApiKeyAuth
func GetArticle(c *gin.Context) {
	appG := app.Gin{C: c}

	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusInternalServerError, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	article, err := articleService.Get()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, article)
}

// @Summary   查询多篇文章
// @Tags   文章
// @Accept json
// @Produce  json
// @Param  tag_id  body  int  false "TagID"
// @Param  state  body  int  false "state"
// @Param  created_by  body  int  false "CreatedBy"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /articles/all  [GET]
// @Security ApiKeyAuth
func GetArticles(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state") //状态只允许0和1
	}

	var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		valid.Min(tagId, 1, "tag_id") //标签ID必须大于0
	}

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{
		TagID:    tagId,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}
	total, err := articleService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_AERTICLR_FAIL, nil)
		return
	}

	articles, err := articleService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}

	data := make(map[string]interface{})
	data["lists"] = articles
	data["total"] = total
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// @Summary   添加文章
// @Tags   文章
// @Accept json
// @Produce  json
// @Param  body  body  models.AddArticleForm  true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /articles/add  [POST]
// @Security ApiKeyAuth
func AddArticle(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo models.AddArticleForm

	err := c.ShouldBindJSON(&reqInfo)
	if err != nil {

		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	valid := validation.Validation{}
	valid.Min(reqInfo.TagID, 1, "tag_id").Message("标签ID必须大于0")
	valid.Required(reqInfo.Title, "title").Message("标题不能为空")
	valid.Required(reqInfo.Desc, "desc").Message("简述不能为空")
	valid.Required(reqInfo.Content, "content").Message("内容不能为空")
	valid.Required(reqInfo.CreatedBy, "created_by").Message("创建人不能为空")
	valid.Range(reqInfo.State, 0, 1, "state").Message("状态只允许0或1")

	if valid.HasErrors() {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	tag := tag_service.Tag{ID: reqInfo.TagID}
	exist, err := tag.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_TAG_FAIL, nil)
		return
	}

	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	articleService := article_service.Article{
		TagID:         reqInfo.TagID,
		Title:         reqInfo.Title,
		Desc:          reqInfo.Desc,
		Content:       reqInfo.Content,
		CoverImageUrl: reqInfo.CoverImageUrl,
		State:         reqInfo.State,
		CreatedBy:     reqInfo.CreatedBy,
	}

	if err := articleService.Add(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//修改文章
func EditArticle(c *gin.Context) {
	appG := app.Gin{C: c}

	valid := validation.Validation{}

	id := com.StrTo(c.Param("id")).MustInt()
	tagId := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	modifiedBy := c.Query("modified_by")

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Min(id, 1, "id").Message("ID必须大于0")
	valid.MaxSize(title, 100, "title").Message("标题最长为100字符")
	valid.MaxSize(desc, 255, "desc").Message("简述最长为255字符")
	valid.MaxSize(content, 65535, "content").Message("内容最长为65535字符")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")

	if ! valid.HasErrors() {
		exists, err := models.ExistArticleByID(id)
		if err != nil {
			appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		}

		if exists {
			exist, err := models.ExistTagByID(tagId)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR, nil)
			}
			if exist {
				data := make(map[string]interface{})
				if tagId > 0 {
					data["tag_id"] = tagId
				}
				if title != "" {
					data["title"] = title
				}
				if desc != "" {
					data["desc"] = desc
				}
				if content != "" {
					data["content"] = content
				}

				data["modified_by"] = modifiedBy

				models.EditArticle(id, data)
				appG.Response(http.StatusOK, e.SUCCESS, nil)
			} else {
				appG.Response(http.StatusInternalServerError, e.ERROR, nil)
			}
		} else {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		}
	} else {
		for _, err := range valid.Errors {
			logger.Info("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}

}

//删除文章
func DeleteArticle(c *gin.Context) {
	appG := app.Gin{C: c}

	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if ! valid.HasErrors() {
		exsits, err := models.ExistArticleByID(id)
		if err != nil {
			appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		}
		if exsits {
			models.DeleteArticle(id)
			appG.Response(http.StatusOK, e.SUCCESS, nil)
		} else {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		}
	} else {
		for _, err := range valid.Errors {
			logger.Info("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}
}

// @Summary   生成海报
// @Tags   文章
// @Accept json
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /articles/poster/generate  [POST]
// @Security ApiKeyAuth
func GenerateArticlePoster(c *gin.Context) {
	appG := app.Gin{C: c}
	article := &article_service.Article{}
	qr := qrcode.NewQrCode(qrcode.QRCODE_URL, 300, 300, qr.M, qr.Auto) // 目前写死 gin 系列路径，可自行增加业务逻辑
	posterName := article_service.GetPosterFlag() + "-" + qrcode.GetQrCodeFileName(qr.URL) + qr.GetQrCodeExt()
	articlePoster := article_service.NewArticlePoster(posterName, article, qr)
	articlePosterBgService := article_service.NewArticlePosterBg(
		"bg.jpg",
		articlePoster,
		&article_service.Rect{
			X0: 0,
			Y0: 0,
			X1: 550,
			Y1: 700,
		},
		&article_service.Pt{
			X: 125,
			Y: 298,
		},
	)
	_, filePath, err := articlePosterBgService.Generate()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GEN_ARTICLE_POSTER_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"poster_url":      qrcode.GetQrCodeFullUrl(posterName),
		"poster_save_url": filePath + posterName,
	})
}
