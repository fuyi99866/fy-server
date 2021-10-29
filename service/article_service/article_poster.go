package article_service

import (
	"go_server/pkg/file"
	"go_server/pkg/logger"
	"go_server/pkg/qrcode"
	"image"
	"image/jpeg"
	"os"
)

type ArticlePoster struct {
	PosterName string
	*Article
	Qr *qrcode.QrcCode
}

func NewArticlePoster(posterName string, article *Article, qr *qrcode.QrcCode) *ArticlePoster {
	return &ArticlePoster{
		PosterName: posterName,
		Article:    article,
		Qr:         qr,
	}
}

func GetPosterFlag() string {
	return "poster"
}

func (a *ArticlePoster) CheckMergeImage(path string) bool {
	if file.CheckNotExist(path+a.PosterName) == true {
		return false
	}
	return true
}

func (a *ArticlePoster) OpenMergeImage(path string) (*os.File, error) {
	f, err := file.MustOpen(a.PosterName, path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

type ArticlePosterBg struct {
	Name string
	*ArticlePoster
	*Rect
	*Pt
}

type Rect struct {
	Name string
	X0   int
	Y0   int
	X1   int
	Y1   int
}

type Pt struct {
	X int
	Y int
}

func NewArticlePosterBg(name string, ap *ArticlePoster, rect *Rect, pt *Pt) *ArticlePosterBg {
	return &ArticlePosterBg{
		Name:          name,
		ArticlePoster: ap,
		Rect:          rect,
		Pt:            pt,
	}
}

func (a *ArticlePosterBg) Generate() (string, string, error) {
	fullPath := qrcode.GetQrCodeFullPath()
	fileName, path, err := a.Qr.Encode(fullPath)
	if err != nil {
		return "", "", err
	}
	if !a.CheckMergeImage(path) {
		mergedF, err := a.OpenMergeImage(path)
		if err != nil {
			return "", "", err
		}
		defer mergedF.Close()

		bgF, err := file.MustOpen(a.Name, path)
		if err != nil {
			return "", "", err
		}
		defer bgF.Close()

		qrF, err := file.MustOpen(fileName, path)
		if err != nil {
			return "", "", err
		}
		defer qrF.Close()
		bgImage, err := jpeg.Decode(bgF)
		if err != nil {
			return "", "", err
		}
		logger.Info("bgImage = ", bgImage)
		qrImage, err := jpeg.Decode(qrF)
		if err != nil {
			return "", "", err
		}
		logger.Info("qrImage = ", qrImage)
		jpg := image.NewRGBA(image.Rect(
			a.Rect.X0,
			a.Rect.Y0,
			a.Rect.X1,
			a.Rect.Y1,
		))
		logger.Info("jpg = ", jpg)

	}
	return fileName, path, nil
}
