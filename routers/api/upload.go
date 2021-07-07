package api

import (
	"github.com/gin-gonic/gin"
	"go_server/pkg/app"
	"go_server/pkg/e"
	"go_server/pkg/logger"
	"go_server/pkg/upload"
	"net/http"
)

func UploadImage(c *gin.Context) {
	appG := app.Gin{C: c}

	data := make(map[string]string)
	file, image, err := c.Request.FormFile("image")
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusBadRequest, e.ERROR, err)
	}

	if image == nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
	} else {
		imageName := upload.GetImageName(image.Filename)
		fullPath := upload.GetImageFullPath()
		savePath := upload.GetImagePath()

		src := fullPath + imageName
		if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
			appG.Response(http.StatusBadRequest, e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT, nil)
		}else {
			err:=upload.CheckImage(fullPath)
			if err != nil{
				logger.Error(err)
				appG.Response(http.StatusBadRequest, e.ERROR_UPLOAD_CHECK_IMAGE_FAIL, nil)
			}else if err:=c.SaveUploadedFile(image,src);err!=nil{
				logger.Error(err)
				appG.Response(http.StatusBadRequest, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
			}else {
				data["image_url"] = upload.GetImageFullUrl(imageName)
				data["image_save_url"] = savePath + imageName
				appG.Response(http.StatusOK,e.SUCCESS,data)
			}
		}
	}

}
