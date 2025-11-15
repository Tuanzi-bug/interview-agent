package model

import (
	"time"
)

//var  UserTitleInteractDao _UserTitleInteract

// UserTitleInteract 用户面试题互动
type (
	//_UserTitleInteract struct {
	//}
	UserTitleInteract struct {
		ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
		UserId    string    `json:"user_id" gorm:"index;not null;comment:用户Id"`
		TitleId   uint64    `json:"title_id" gorm:"index;comment:面试题目Id"`
		Interact  int       `json:"interact" gorm:"not null;default:0;comment:互动类型（0=历史记录, 1=推荐, 2=收藏）"`
		CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime:milli"`
		UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime:milli"`
		Deleted   int       `json:"deleted" gorm:"not null;default:0;comment:删除状态（0=未删除, 1=已删除）"`
	}
)

// TableName 指定表名
func (i *UserTitleInteract) TableName() string {
	return "user_title_interact"
}
