package model

import (
	"time"
)

//var InterviewTitleDao _InterviewTitle

// InterviewTitle 面试题
type (
	//_InterviewTitle struct {
	//}
	InterviewTitle struct {
		ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
		Title     string    `json:"tile" gorm:"type:text;not null;comment:面试标题"`
		Type      string    `json:"type" gorm:"size:255;not null;comment:面试类型(基础概念、算法)"`
		Domain    string    `json:"domain" gorm:"size:255;not null;comment:面试领域(java、golang)"`
		Level     string    `json:"level" gorm:"size:128;comment:难度级别（简单、中等、困难）"`
		CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime:milli"`
		UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime:milli"`
		Deleted   int       `json:"deleted" gorm:"not null;default:0;comment:删除状态（0=未删除, 1=已删除）"`
	}
)

// TableName 指定表名
func (i *InterviewTitle) TableName() string {
	return "interview_title"
}
