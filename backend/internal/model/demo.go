package model

import (
	"time"

	"gorm.io/gorm"
)

var DemoDao _Demo

// Demo 用户模型
type (
	_Demo struct {
	}
	Demo struct {
		ID        uint           `json:"id" gorm:"primaryKey"`
		Name      string         `json:"name" gorm:"size:50;not null"`
		CreatedAt time.Time      `json:"created_at"`
		UpdatedAt time.Time      `json:"updated_at"`
		DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	}
)

func (Demo) TableName() string {
	return "Demo"
}

// Create 创建用户记录
func (u *_Demo) Create(Demo *Demo) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	return getDB().Create(Demo).Error
}
