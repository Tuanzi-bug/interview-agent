package model

import (
	"time"
)

var InterviewDialogueDao _InterviewDialogue

// InterviewDialogue 面试对话记录
// 对应 "回答"、"追问" 这些交错的对话条目
type (
	_InterviewDialogue struct {
	}
	InterviewDialogue struct {
		ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
		UserID    uint      `json:"user_id" gorm:"not null;index:idx_user_id;comment:用户ID"`
		ParentID  uint64    `json:"parent_id" gorm:"not null;index:idx_parent_id;comment:父对话记录ID"`
		ReportID  uint64    `json:"report_id" gorm:"not null;index:idx_report_id;comment:关联的面试报告ID"`
		Question  string    `json:"question" gorm:"type:text;comment:智能体的提问内容"`
		Answer    string    `json:"answer" gorm:"type:text;comment:用户的回答内容"`
		CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime:milli"`
	}
)

// TableName 指定表名
func (InterviewDialogue) TableName() string {
	return "interview_dialogues"
}

// Create 创建对话记录
func (dao *_InterviewDialogue) Create(dialogue *InterviewDialogue) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	return getDB().Create(dialogue).Error
}

// GetByReportID 根据报告ID获取所有对话记录
func (dao *_InterviewDialogue) GetByReportID(reportID uint64) ([]*InterviewDialogue, error) {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	var dialogues []*InterviewDialogue
	err := getDB().Where("report_id = ?", reportID).Order("created_at ASC").Find(&dialogues).Error
	return dialogues, err
}

// GetByParentID 根据父ID获取追问记录
func (dao *_InterviewDialogue) GetByParentID(parentID uint64) ([]*InterviewDialogue, error) {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	var dialogues []*InterviewDialogue
	err := getDB().Where("parent_id = ?", parentID).Order("created_at ASC").Find(&dialogues).Error
	return dialogues, err
}

// GetMainQuestions 获取所有主问题（ParentID = 0）
func (dao *_InterviewDialogue) GetMainQuestions(reportID uint64) ([]*InterviewDialogue, error) {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	var dialogues []*InterviewDialogue
	err := getDB().Where("report_id = ? AND parent_id = 0", reportID).Order("created_at ASC").Find(&dialogues).Error
	return dialogues, err
}

// GetFollowUpQuestions 获取所有追问（ParentID != 0）
func (dao *_InterviewDialogue) GetFollowUpQuestions(reportID uint64) ([]*InterviewDialogue, error) {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	var dialogues []*InterviewDialogue
	err := getDB().Where("report_id = ? AND parent_id != 0", reportID).Order("created_at ASC").Find(&dialogues).Error
	return dialogues, err
}
