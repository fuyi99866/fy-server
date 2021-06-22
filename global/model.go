package global

import (
	"time"
)

type GVA_MODEL struct {
	ID uint `gorm:"primarykey"`
	CreateAt time.Time
	UpdateAt time.Time
	//DeleteAt gorm.DeleteAt `gorm:"index" json:"-"`
	DeleteAt time.Time
}
