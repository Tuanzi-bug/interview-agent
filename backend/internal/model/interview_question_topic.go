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
		UserID        uint      `json:"user_id" gorm:"not null;index:idx_user_id;comment:用户ID"`
		ReportID      uint64    `json:"report_id" gorm:"not null;index:idx_report_id;comment:关联的面试报告ID"`
		QuestionText  string    `json:"question_text" gorm:"type:text;not null;comment:面试官提出的问题内容"`
		DisplayOrder  uint32    `json:"display_order" gorm:"not null;default:0;comment:问题主题的显示顺序(例如1/5中的1)"`
		EvalDimension string    `json:"eval_dimension" gorm:"type:enum('professional_field','project_experience','technical_depth','technical_foundation','team_collaboration','system_architecture_design');not null;comment:问题维度"`
		CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime:milli"`
		UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime:milli"`
	}

	// QuestionTopicWithDialogues 问题主题与对话的组合结构体
	QuestionTopicWithDialogues struct {
		// 问题主题字段
		ID            uint64    `json:"id"`
		UserID        uint      `json:"user_id"`
		ReportID      uint64    `json:"report_id"`
		QuestionText  string    `json:"question_text"`
		DisplayOrder  uint32    `json:"display_order"`
		EvalDimension string    `json:"eval_dimension"`
		CreatedAt     time.Time `json:"created_at"`
		// 对话数组
		Dialogues []*DialogueWithoutTopic `json:"dialogues"`
	}

	// DialogueWithoutTopic 对话结构体（不包含主题信息）
	DialogueWithoutTopic struct {
		ID           uint64    `json:"id"`
		UserID       uint      `json:"user_id"`
		TopicID      uint64    `json:"topic_id"`
		Question     string    `json:"question"`
		Answer       string    `json:"answer"`
		DisplayOrder uint32    `json:"display_order"`
		CreatedAt    time.Time `json:"created_at"`
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

// GetWithDialoguesByUserIDAndReportID 通过 UserID 和 ReportID 查询问题主题及其关联的对话
// 返回 QuestionTopicWithDialogues 结构体数组，包含主题和对应的对话列表
func (dao *_InterviewQuestionTopic) GetWithDialoguesByUserIDAndReportID(userID uint, reportID uint64) ([]*QuestionTopicWithDialogues, error) {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}

	// 第一步：查询所有问题主题
	var topics []*InterviewQuestionTopic
	err := getDB().
		Where("user_id = ? AND report_id = ?", userID, reportID).
		Order("display_order ASC").
		Find(&topics).Error
	if err != nil {
		return nil, err
	}

	// 第二步：为每个主题查询对应的对话
	result := make([]*QuestionTopicWithDialogues, 0, len(topics))
	for _, topic := range topics {
		// 查询该主题下的所有对话
		var dialogues []*DialogueWithoutTopic
		err := getDB().
			Table("interview_dialogues").
			Where("topic_id = ?", topic.ID).
			Order("display_order ASC").
			Scan(&dialogues).Error
		if err != nil {
			return nil, err
		}

		// 构建组合结构体
		topicWithDialogues := &QuestionTopicWithDialogues{
			ID:            topic.ID,
			UserID:        topic.UserID,
			ReportID:      topic.ReportID,
			QuestionText:  topic.QuestionText,
			DisplayOrder:  topic.DisplayOrder,
			EvalDimension: topic.EvalDimension,
			CreatedAt:     topic.CreatedAt,
			Dialogues:     dialogues,
		}
		result = append(result, topicWithDialogues)
	}

	return result, nil
}
