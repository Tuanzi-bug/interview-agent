package model

import (
	"time"
)

var InterviewQuestionTopicDao _InterviewQuestionTopic

// InterviewQuestionTopic 面试问题主题
// 对应截图中 "1/5 我看你简历上提到了..." 这一整个区块
type (
	_InterviewQuestionTopic struct {
	}
	InterviewQuestionTopic struct {
		ID            uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
		ReportID      uint64    `json:"report_id" gorm:"not null;index:idx_report_id;comment:关联的面试报告ID"`
		QuestionText  string    `json:"question_text" gorm:"type:text;not null;comment:面试官提出的初始问题内容"`
		DisplayOrder  uint32    `json:"display_order" gorm:"not null;default:0;comment:问题主题的显示顺序(例如1/5中的1)"`
		EvalDimension string    `json:"eval_dimension" gorm:"type:enum('professional_field','project_experience','technical_depth','technical_foundation','team_collaboration','system_architecture_design');not null;comment:问题维度"`
		CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime:milli"`
		UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime:milli"`
	}
)

// TableName 指定表名
func (InterviewQuestionTopic) TableName() string {
	return "interview_question_topics"
}

// Create 创建问题主题
func (dao *_InterviewQuestionTopic) Create(topic *InterviewQuestionTopic) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	return getDB().Create(topic).Error
}
