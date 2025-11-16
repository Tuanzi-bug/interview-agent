package model

import (
	"time"
)

var InterviewTitleDao _InterviewTitle

// InterviewTitle 面试题
type (
	_InterviewTitle struct {
	}
	InterviewTitle struct {
		ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
		Title     string    `json:"title" gorm:"type:text;not null;comment:面试标题"`
		Type      string    `json:"type" gorm:"size:255;not null;comment:面试类型(基础概念、算法)"`
		Domain    string    `json:"domain" gorm:"size:255;not null;comment:面试领域(java、golang)"`
		Level     string    `json:"level" gorm:"size:128;not null;comment:难度级别（简单、中等、困难）"`
		CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime:milli"`
		UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime:milli"`
		Deleted   int       `json:"deleted" gorm:"not null;default:0;comment:删除状态（0=未删除, 1=已删除）"`
	}
)

// TableName 指定表名
func (t *InterviewTitle) TableName() string {
	return "interview_title"
}

// Create 创建面试题
func (t *_InterviewTitle) CreateInterviewTitle(data *InterviewTitle) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	err := getDB().Model(&InterviewTitle{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

// FindByTitle
func (t *_InterviewTitle) FindInterviewTitleByTitle(title string) (*InterviewTitle, error) {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	var interviewTitle InterviewTitle
	err := getDB().
		Where("title = ?", title).
		First(&interviewTitle).Error
	if err != nil {
		return nil, err
	}
	return &interviewTitle, nil
}

// FindByID
func (t *_InterviewTitle) FindInterviewTitleById(id uint64) (*InterviewTitle, error) {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	var interviewTitle InterviewTitle
	err := getDB().
		Where("id = ?", id).
		First(&interviewTitle).Error
	if err != nil {
		return nil, err
	}
	return &interviewTitle, nil
}
