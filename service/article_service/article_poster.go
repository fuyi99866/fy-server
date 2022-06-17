package article_service

import (
	"github.com/golang/freetype"
	"github.com/sirupsen/logrus"
	"go_server/pkg/file"
	"go_server/pkg/qrcode"
	"go_server/pkg/setting"
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
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

type DrawText struct {
	JPG    draw.Image
	Merged *os.File

	Title string
	X0    int
	Y0    int
	Size0 float64

	SubTitle string
	X1       int
	Y1       int
	Size1    float64
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
	logrus.Info("1111111111111  ",a.CheckMergeImage(path))
	if !a.CheckMergeImage(path) {
		mergedF, err := a.OpenMergeImage(path)
		if err != nil {
			logrus.Info("mergedF  ",err)
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
			logrus.Info("qrF  ",err)
			return "", "", err
		}
		defer qrF.Close()
		logrus.Info("bgF  ",bgF)
		bgImage, err := jpeg.Decode(bgF)
		if err != nil {
			logrus.Info("bgImage  ",err)
			return "", "", err
		}
		logrus.Info("bgImage = ", bgImage)

		qrImage, err := jpeg.Decode(qrF)
		if err != nil {
			logrus.Info("qrImage  ",err)
			return "", "", err
		}
		logrus.Info("qrImage = ", qrImage)

		jpg := image.NewRGBA(image.Rect(
			a.Rect.X0,
			a.Rect.Y0,
			a.Rect.X1,
			a.Rect.Y1,
		))
		logrus.Info("jpg = ", jpg)

		draw.Draw(jpg,jpg.Bounds(),bgImage,bgImage.Bounds().Min,draw.Over)
		draw.Draw(jpg,jpg.Bounds(),qrImage,qrImage.Bounds().Min.Sub(image.Pt(a.Pt.X,a.Pt.Y)),draw.Over)
		jpeg.Encode(mergedF, jpg, nil)

		err = a.DrawPoster(&DrawText{
			JPG:      jpg,
			Merged:   mergedF,
			Title:    "golang Gin 系列文章",
			X0:       80,
			Y0:       160,
			Size0:    42,
			SubTitle: "---小厮",
			X1:       320,
			Y1:       220,
			Size1:    36,
		},"msyhbd.ttc")
		if err!=nil{
			return "","",err
		}
	}
	return fileName, path, nil
}

func (a *ArticlePosterBg) DrawPoster(d *DrawText, fontName string) error {
	fontSource := setting.AppSetting.RuntimeRootPath + setting.AppSetting.FontSavePath + fontName
	fontSourceBytes, err := ioutil.ReadFile(fontSource)
	if err != nil {
		return err
	}
	trueTypeFont, err := freetype.ParseFont(fontSourceBytes)
	if err != nil {
		return err
	}
	fc := freetype.NewContext()
	fc.SetDPI(72)
	fc.SetFont(trueTypeFont)
	fc.SetClip(d.JPG.Bounds())
	fc.SetDst(d.JPG)
	fc.SetSrc(image.Black)

	pt := freetype.Pt(d.X0, d.Y0)
	_, err = fc.DrawString(d.Title, pt)
	if err != nil {
		return err
	}
	fc.SetFontSize(d.Size1)
	_, err = fc.DrawString(d.SubTitle, freetype.Pt(d.X1, d.Y1))
	if err != nil {
		return err
	}
	err = jpeg.Encode(d.Merged, d.JPG, nil)
	if err != nil {
		return err
	}
	return err
}
