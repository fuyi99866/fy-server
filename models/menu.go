package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

type SysMenu struct {
	SysBaseMenu
	MenuId      string                 `json:"menuId" gorm:"comment:菜单ID"`
	AuthorityId string                 `json:"-" gorm:"comment:角色ID"`
	Children    []SysMenu              `json:"children" gorm:"-"`
	Parameters  []SysBaseMenuParameter `json:"parameters" gorm:"foreignKey:SysBaseMenuID;references:MenuId"`
}

func (s SysMenu) TableName() string {
	return "authority_menu"
}

func GetMenuTreeMap(authorityId string) (err error, treeMap map[string][]SysMenu) {
	var allMenus []SysMenu
	treeMap = make(map[string][]SysMenu)
	err = db.Where("authority_id = ?", authorityId).Order("sort").Preload("Parameters").Find(&allMenus).Error
	for _, v := range allMenus {
		treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
	}
	return err, treeMap
}


//////////////////////////////////////////////////////
type Menu struct {
	gorm.Model
	Name      string `json:"name" gorm:"column:name"`
	Path      string `json:"path" gorm:"column:path"`
	Component string `json:"component" gorm:"column:component"`
	Url       string `json:"url" gorm:"column:url"`
	Status    int    `json:"status" gorm:"column:status"`
}

type MenuResult struct {
	ID        string `json:"id"`        //菜单id
	Name      string `json:"name"`      //菜单名
	Path      string `json:"path"`      //菜单路径
	Component string `json:"component"` //菜单组件
	Url       string `json:"url"`       //菜单接口
	Status    int    `json:"status"`    //菜单状态
}

type MenuAdd struct {
	Name      string `json:"name"`      //菜单名
	Path      string `json:"path"`      //菜单路径
	Component string `json:"component"` //菜单组件
	Url       string `json:"url"`       //菜单接口
	Status    int    `json:"status"`    //菜单状态
}

//查询所有菜单
func GetAllMenus() ([]*MenuResult, error) {
	var menus []*MenuResult
	err := db.Model(&Menu{}).Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func ExistMenu(name string) (bool, error) {
	var menu Menu
	err := db.Model(&Menu{}).Where("name = ?", name).First(&menu).Error

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}

	if err != nil && err != gorm.ErrRecordNotFound{
		return false, err
	}

	return true, nil
}

//添加菜单
func AddMenu(menu Menu) (uint, error) {
	//查询菜单是否已存在
	exist, err := ExistMenu(menu.Name)
	if err != nil {
		return 0, err
	}

	if exist {
		return 0, errors.New("menu exist")
	}

	err = db.Model(&Menu{}).Create(&menu).Error
	if err != nil {
		return 0, err
	}
	return menu.ID, nil
}

//删除菜单
func DeleteMenuByName(name string) error {
	var menu *Menu
	err := db.Model(&Menu{}).Where("name = ?", name).Delete(&menu).Error
	if err != nil {
		return err
	}
	return nil
}
