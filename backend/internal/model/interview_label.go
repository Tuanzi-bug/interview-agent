package model

import (
	"time"
)

var InterviewLabelDao _InterviewLabel

// InterviewLabel 面试题标签
type (
	_InterviewLabel struct {
	}
	InterviewLabel struct {
		ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
		Label     string    `json:"label" gorm:"size:255;not null;comment:面试标签"`
		TitleId   uint64    `json:"title_id" gorm:"index;not null;comment:面试题目Id"`
		CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime:milli"`
		UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime:milli"`
		Deleted   int       `json:"deleted" gorm:"not null;default:0;comment:删除状态（0=未删除, 1=已删除）"`
	}
)

// TableName 指定表名
func (l *InterviewLabel) TableName() string {
	return "interview_label"
}

// Create 创建面试题标签
func (l *_InterviewLabel) CreateInterviewLabel(data *InterviewLabel) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	err := getDB().Model(&InterviewLabel{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
