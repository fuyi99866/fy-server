package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Tag struct {
	gorm.Model
	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

type AddTagForm struct {
	Name      string
	CreatedBy string
	State     int
}

type ExportTagForm struct {
	Name      string
	State     int
}

func GetTags( maps interface{}) ([]Tag, error) {
	var tags []Tag
	err := db.Where(maps).Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func GetTagTotal(maps interface{}) (int, error) {
	var count int
	err := db.Model(&Tag{}).Where(maps).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func ExistTagByName(name string) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("name = ?", name).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if tag.ID > 0 {
		return true, nil
	}
	return false, nil
}

func AddTag(name string, state int, createdBy string) error {
	err := db.Create(&Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	}).Error

	return err
}

func ExistTagByID(id int) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("id = ?", id).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}

	return false, nil
}

func DeleteTag(id int) error {
	err := db.Where("id = ?", id).Delete(&Tag{}).Error

	return err
}

func EditTag(id int, data interface{}) error {
	err := db.Model(&Tag{}).Where("id = ?", id).Updates(data).Error

	return err
}

//回调函数 Callbacks
/**
可以将回调方法定义为模型结构的指针，在创建、更新、查询、删除时将被调用，
如果任何回调返回错误，gorm将停止未来操作并回滚所有更改
*/
func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreateAt", time.Now().Unix())
	return nil
}

func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedAt", time.Now().Unix())
	return nil
}

//注意硬删除要使用 Unscoped()，这是 GORM 的约定
func CleanAllTag() error {
	err := db.Unscoped().Where("delete_on != ?", 0).Delete(&Tag{}).Error
	return err
}
