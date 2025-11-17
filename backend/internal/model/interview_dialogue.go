package model

import (
	"time"
)

var InterviewDialogueDao _InterviewDialogue

// SpeakerType 发言人类型
type SpeakerType string

const (
	SpeakerTypeInterviewer SpeakerType = "interviewer" // 面试官
	SpeakerTypeCandidate   SpeakerType = "candidate"   // 候选人
)

// InterviewDialogue 面试对话记录
// 对应 "回答"、"追问" 这些交错的对话条目
type (
	_InterviewDialogue struct {
	}
	InterviewDialogue struct {
		ID           uint64      `json:"id" gorm:"primaryKey;autoIncrement"`
		TopicID      uint64      `json:"topic_id" gorm:"not null;index:idx_topic_id_order;comment:关联的问题主题ID"`
		SpeakerType  SpeakerType `json:"speaker_type" gorm:"type:enum('interviewer','candidate');not null;comment:发言人类型"`
		Content      string      `json:"content" gorm:"type:text;not null;comment:对话内容(回答或追问的文本)"`
		DisplayOrder uint32      `json:"display_order" gorm:"not null;default:0;index:idx_topic_id_order;comment:在主题内的显示顺序"`
		CreatedAt    time.Time   `json:"created_at" gorm:"autoCreateTime:milli"`
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
