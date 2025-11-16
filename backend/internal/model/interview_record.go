package model

import (
	"fmt"
	"time"
)

var InterviewRecordDao _InterviewRecord

// InterviewRecord 面试记录模型
type (
	_InterviewRecord struct {
	}
	InterviewRecord struct {
		ID             uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
		UserID         uint       `json:"user_id" gorm:"index;not null;comment:用户ID"`
		Title          string     `json:"title" gorm:"size:255;not null;comment:面试标题"`
		Query          string     `json:"query" gorm:"type:text;not null;comment:初始查询/问题"`
		Messages       string     `json:"messages" gorm:"type:json;comment:对话历史（JSON格式）"`
		Report         string     `json:"report" gorm:"type:text;comment:最终报告"`
		Status         string     `json:"status" gorm:"size:50;not null;default:'pending';comment:面试状态（pending/resume_analysis/question_generation/answer_evaluation/report_generation/completed）"`
		CurrentAgent   string     `json:"current_agent" gorm:"size:100;comment:当前活跃的Agent名称"`
		Duration       int64      `json:"duration" gorm:"comment:面试耗时（秒）"`
		Score          *float64   `json:"score" gorm:"comment:面试评分"`
		Feedback       string     `json:"feedback" gorm:"type:text;comment:反馈信息"`
		LastModifiedAt time.Time  `json:"last_modified_at" gorm:"autoCreateTime:milli;comment:最后修改时间戳，用于并发控制"`
		CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime:milli"`
		UpdatedAt      time.Time  `json:"updated_at" gorm:"autoUpdateTime:milli"`
		CompletedAt    *time.Time `json:"completed_at" gorm:"comment:完成时间"`
	}
)

// TableName 指定表名
func (i *InterviewRecord) TableName() string {
	return "interview_record"
}

// CreateInterviewRecord 创建面试记录
func (i *_InterviewRecord) CreateInterviewRecord(record *InterviewRecord) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	return getDB().Create(record).Error
}

// GetInterviewRecordByID 根据ID查询面试记录
func (i *_InterviewRecord) GetInterviewRecordByID(id uint64) (*InterviewRecord, error) {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	var record InterviewRecord
	err := getDB().Where("id = ?", id).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// ListInterviewRecords 查询用户的面试记录列表
func (i *_InterviewRecord) ListInterviewRecords(userID uint, page, pageSize int) ([]*InterviewRecord, int64, error) {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	var records []*InterviewRecord
	var total int64

	query := getDB().Where("user_id = ?", userID)

	// 获取总数
	if err := query.Model(&InterviewRecord{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := query.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// UpdateInterviewRecord 更新面试记录（带时间戳检查）
func (i *_InterviewRecord) UpdateInterviewRecord(record *InterviewRecord) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}

	// 先获取当前记录的最后修改时间
	var currentRecord InterviewRecord
	if err := getDB().Where("id = ?", record.ID).First(&currentRecord).Error; err != nil {
		return err
	}

	// 检查时间戳：如果当前记录的修改时间晚于请求中的时间戳，说明有更新的操作
	if !record.LastModifiedAt.IsZero() && currentRecord.LastModifiedAt.After(record.LastModifiedAt) {
		return fmt.Errorf("update rejected: record was modified after this update was initiated (current: %v, request: %v)",
			currentRecord.LastModifiedAt, record.LastModifiedAt)
	}

	// 更新记录和时间戳
	now := time.Now()
	result := getDB().Model(&InterviewRecord{}).
		Where("id = ?", record.ID).
		Updates(map[string]interface{}{
			"messages":         record.Messages,
			"status":           record.Status,
			"current_agent":    record.CurrentAgent,
			"last_modified_at": now,
		})

	return result.Error
}

// UpdateInterviewRecordStatus 更新面试记录状态
func (i *_InterviewRecord) UpdateInterviewRecordStatus(id uint64, status string, currentAgent string) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	updates := make(map[string]interface{})
	if status != "" {
		updates["status"] = status
	}
	if currentAgent != "" {
		updates["current_agent"] = currentAgent
	}
	if len(updates) == 0 {
		return nil // 没有需要更新的字段
	}
	return getDB().Model(&InterviewRecord{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// CompleteInterviewRecord 完成面试记录（带时间戳检查）
func (i *_InterviewRecord) CompleteInterviewRecord(id uint64, report string, duration int64, score *float64) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}

	// 先获取当前记录，检查是否已完成
	var record InterviewRecord
	if err := getDB().Where("id = ?", id).First(&record).Error; err != nil {
		return err
	}

	// 防止重复完成：如果已经是 completed 状态，返回错误
	if record.Status == "completed" {
		return fmt.Errorf("complete failed: record already completed (id=%d)", id)
	}

	now := time.Now()
	result := getDB().Model(&InterviewRecord{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":           "completed",
			"report":           report,
			"duration":         duration,
			"score":            score,
			"completed_at":     now,
			"last_modified_at": now,
		})

	return result.Error
}

// DeleteInterviewRecord 删除面试记录
func (i *_InterviewRecord) DeleteInterviewRecord(id uint64) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	return getDB().Where("id = ?", id).Delete(&InterviewRecord{}).Error
}
