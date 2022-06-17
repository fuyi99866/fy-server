package tag_service

import (
	"encoding/json"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx"
	"go_server/models"
	"go_server/pkg/export"
	"go_server/pkg/file"
	"go_server/pkg/gredis"
	"go_server/service/cache_service"
	"io"
	"strconv"
	"time"
)

type Tag struct {
	ID        int
	Name      string
	CreatedBy string
	State     int
	PageNum   int
	PageSize  int
}

func (t *Tag) ExistByID() (bool, error) {
	return models.ExistTagByID(t.ID)
}

//导出
func (t *Tag) Export() (string, error) {
	tags, err := t.GetAll()
	logrus.Info("get all tags: ", tags)
	if err != nil {
		return "", err
	}
	xlsFile := xlsx.NewFile()
	sheet, err := xlsFile.AddSheet("标签信息")
	if err != nil {
		return "", err
	}
	titles := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	row := sheet.AddRow()
	var cell *xlsx.Cell
	for _, title := range titles {
		cell = row.AddCell()
		cell.Value = title
	}
	for _, v := range tags {
		values := []string{
			strconv.Itoa(int(v.ID)),
			v.Name,
			v.CreatedBy,
			v.CreatedAt.String(),
			v.ModifiedBy,
			v.UpdatedAt.String(),
		}
		row = sheet.AddRow()
		for _, value := range values {
			cell = row.AddCell()
			cell.Value = value
		}
	}
	time := strconv.Itoa(int(time.Now().Unix()))
	filename := "tags-" + time + export.EXT
	dirFullPath := export.GetExcelFullPath()
	err = file.IsNotExistMkDir(dirFullPath)
	if err != nil {
		return "", err
	}

	err = xlsFile.Save(dirFullPath + filename)
	if err != nil {
		return "", err
	}

	return filename, nil
}

//导入
func (t *Tag) Import(r io.Reader) error {
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}
	rows := xlsx.GetRows("标签信息")
	for irow, row := range rows {
		if irow > 0 {
			var data []string
			for _, cell := range row {
				data = append(data, cell)
			}
			models.AddTag(data[1], 1, data[2])
		}
	}
	return nil
}

func (t *Tag) GetAll() ([]models.Tag, error) {
	var (
		tags, cacheTags []models.Tag
	)
	cache := cache_service.Tag{
		State: t.State,
		Name:  t.Name,
	}

	key := cache.GetTagsKey()
	logrus.Info("key = ", key)
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logrus.Error(err)
		} else {
			json.Unmarshal(data, &cacheTags)
			return cacheTags, nil
		}
	}
	tags, err := models.GetTags(t.getMap())
	if err != nil {
		return nil, err
	}
	gredis.Set(key, tags, 3600)
	return tags, nil
}

func (t *Tag) getMap() interface{} {
	maps := make(map[string]interface{})
	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}
	return maps
}

func (t *Tag) Count() (int, error) {
	return models.GetTagTotal(t.getMap())
}

func (t *Tag) ExistByName() (bool, error) {
	return models.ExistTagByName(t.Name)
}

func (t *Tag) Add() error {
	return models.AddTag(t.Name, t.State, t.CreatedBy)
}
