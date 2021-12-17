package models

import "github.com/jinzhu/gorm"

type Robot struct {
	gorm.Model
	SN       string `json:"sn" gorm:"index;not null"` //是因为限定了SN是唯一的吗?
	NickName string `json:"nickname" gorm:"column:nickname;" `
	Company  string `json:"company" gorm:"column:company;"`
}

func GetAllRobot() ([]*Robot, error) {
	var robots []*Robot
	err := db.Find(&robots).Error
	if err != nil {
		return nil, err
	}

	return robots, nil
}
