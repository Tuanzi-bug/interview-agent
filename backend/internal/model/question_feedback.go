package model

import (
	"time"
)

var QuestionFeedbackDao _QuestionFeedback

type (
	_QuestionFeedback struct {
	}
	QuestionFeedback struct {
		ID                uint64    `json:"id" gorm:"primaryKey;autoIncrement;comment:ID"`
		TopicID           uint64    `json:"topic_id" gorm:"not null;comment:关联的问题主题ID。每一轮一个 topic"`
		InterviewRecordID uint64    `json:"interview_record_id" gorm:"not null;index:idx_interview_record_id;comment:关联的主面试记录ID"`
		Score             uint8     `json:"score" gorm:"not null;comment:本题得分 (0-100)"`
		KeyAssessment     string    `json:"key_assessment" gorm:"size:255;default:'';comment:重点考察 需要提前设置"`
		DifficultyLevel   string    `json:"difficulty_level" gorm:"size:50;default:'';comment:难度等级"`
		Strengths         string    `json:"strengths" gorm:"type:text;comment:优势点评"`
		Weaknesses        string    `json:"weaknesses" gorm:"type:text;comment:不足点评"`
		Suggestions       string    `json:"suggestions" gorm:"type:text;comment:建议内容"`
		AssessmentPoints  string    `json:"assessment_points" gorm:"type:text;comment:考点分析"`
		ReferenceAnswer   string    `json:"reference_answer" gorm:"type:text;comment:参考答案"`
		ThoughtProcess    string    `json:"thought_process" gorm:"size:255;comment:答题思路"`
		CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
		UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	}
)

// TableName �h
func (q *QuestionFeedback) TableName() string {
	return "question_feedback"
}

// CreateQuestionFeedback ��T͈
func (q *_QuestionFeedback) CreateQuestionFeedback(feedback *QuestionFeedback) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	return getDB().Create(feedback).Error
}

// GetQuestionFeedbackByID 9nID���T͈
func (q *_QuestionFeedback) GetQuestionFeedbackByID(id uint64) (*QuestionFeedback, error) {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	var feedback QuestionFeedback
	err := getDB().Where("id = ?", id).First(&feedback).Error
	if err != nil {
		return nil, err
	}
	return &feedback, nil
}

// ListQuestionFeedbackByInterviewRecordID 9nbհUID���T͈h
func (q *_QuestionFeedback) ListQuestionFeedbackByInterviewRecordID(interviewRecordID uint64) ([]*QuestionFeedback, error) {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	var feedbacks []*QuestionFeedback
	err := getDB().Where("interview_record_id = ?", interviewRecordID).
		Order("created_at ASC").
		Find(&feedbacks).Error
	if err != nil {
		return nil, err
	}
	return feedbacks, nil
}

// ListQuestionFeedbackByTopicID 9n;�ID���T͈h
func (q *_QuestionFeedback) ListQuestionFeedbackByTopicID(topicID uint64) ([]*QuestionFeedback, error) {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	var feedbacks []*QuestionFeedback
	err := getDB().Where("topic_id = ?", topicID).
		Order("created_at ASC").
		Find(&feedbacks).Error
	if err != nil {
		return nil, err
	}
	return feedbacks, nil
}

// UpdateQuestionFeedback ���T͈
func (q *_QuestionFeedback) UpdateQuestionFeedback(feedback *QuestionFeedback) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	return getDB().Model(&QuestionFeedback{}).
		Where("id = ?", feedback.ID).
		Updates(feedback).Error
}

// DeleteQuestionFeedback  d�T͈
func (q *_QuestionFeedback) DeleteQuestionFeedback(id uint64) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	return getDB().Where("id = ?", id).Delete(&QuestionFeedback{}).Error
}
