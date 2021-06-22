package models

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func GetAllUser() ([]*User, error) {
	var user []*User
	err := db.Find(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func AddUser(data map[string]interface{}) (id uint, err error) {
	user := User{
		Username:  data["username"].(string),
		Password:  data["password"].(string),
	}
	if err := db.Create(&user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

//通过用户和密码查询用户
func CheckUser(username, password string) (bool, error) {
	var user User
	err := db.Select("id").Where(&User{Username: username, Password: password}).First(&user).Error
	logrus.Debugln("CheckUser",  user.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if user.ID > 0 {
		return true, nil
	}
	return false, nil
}
