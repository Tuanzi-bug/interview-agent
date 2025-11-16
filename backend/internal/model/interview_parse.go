package model

import (
	"time"
)

var InterviewParseDao _InterviewParse

// InterviewParse 面试题解析
type (
	_InterviewParse struct {
	}
	InterviewParse struct {
		ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
		TitleId   uint64    `json:"title_id" gorm:"index;not null;comment:面试题目Id"`
		Parse     string    `json:"parse" gorm:"type:text;not null;comment:面试题解析"`
		Url       string    `json:"url" gorm:"type:text;comment:面试题解析图片URL "`
		CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime:milli"`
		UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime:milli"`
		Deleted   int       `json:"deleted" gorm:"not null;default:0;comment:删除状态（0=未删除, 1=已删除）"`
	}
)

// TableName 指定表名
func (p *InterviewParse) TableName() string {
	return "interview_parse"
}

// Create 创建面试题解析
func (p *_InterviewParse) CreateInterviewParse(data *InterviewParse) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	err := getDB().Model(&InterviewParse{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
