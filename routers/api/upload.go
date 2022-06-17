package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go_server/pkg/app"
	"go_server/pkg/e"
	"go_server/pkg/upload"
	"net/http"
	"path/filepath"
)

//TODO 图片上传接口，需要实现

// @Summary   上传图片
// @Tags   上传下载
// @Accept multipart/form-data
// @Param file formData file true "file"
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /upload_img  [POST]
// @Security ApiKeyAuth
func UploadImage(c *gin.Context) {
	appG := app.Gin{C: c}

	data := make(map[string]string)
	file, image, err := c.Request.FormFile("file") //获取上传的图片，返回提供表单键的第一个文件
	if err != nil {
		logrus.Error(err)
		appG.Response(http.StatusBadRequest, e.ERROR, err)
		return
	}

	if image == nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	} else {
		imageName := upload.GetImageName(image.Filename)
		fullPath := upload.GetImageFullPath()
		savePath := upload.GetImagePath()
		src := fullPath + imageName
		logrus.Info("file = ", file)
		logrus.Info("image = ", image)
		logrus.Info("imageName = ", imageName, upload.CheckImageExt(imageName))
		if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) { //检查图片的大小和后缀
			appG.Response(http.StatusBadRequest, e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT, nil)
			return
		} else {
			err := upload.CheckImage(fullPath) //检查上传图片所需（权限、文件夹）
			if err != nil {
				logrus.Error(err)
				appG.Response(http.StatusBadRequest, e.ERROR_UPLOAD_CHECK_IMAGE_FAIL, nil)
			} else if err := c.SaveUploadedFile(image, src); err != nil { //保存图片
				logrus.Error(err)
				appG.Response(http.StatusBadRequest, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
			} else {
				data["image_url"] = upload.GetImageFullUrl(imageName)
				data["image_save_url"] = savePath + imageName
				appG.Response(http.StatusOK, e.SUCCESS, data)
			}
		}
	}
}

// @Summary   上传文件
// @Tags   上传下载
// @Accept multipart/form-data
// @Param file formData file true "file"
// @Produce  json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /upload_file  [POST]
// @Security ApiKeyAuth
func UploadFile(c *gin.Context) {
	appG := app.Gin{C: c}
	data := make(map[string]string)
	// 单文件
	file, err := c.FormFile("file") //解析提交的表单
	if err != nil || file == nil {
		logrus.Info("UploadFill", err)
		c.String(http.StatusBadRequest, "")
		return
	}
	logrus.Info("file:", file)
	logrus.Info("UploadFile:", file.Filename)

	//设置文件存储的地址
	//fullPath := setting.AppSetting.RuntimeRootPath + "upload/images/"
	fullPath := "dist/"
	filename := fullPath + filepath.Base(file.Filename)
	// 上传文件到指定的路径
	c.SaveUploadedFile(file, filename)

	url := "http://127.0.0.1:8081"+"/upload/images/"+file.Filename
	logrus.Info("Url   ", url)
	data["url"] = url
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
