package model

import (
	"time"
)

var InterviewRecordDao _InterviewRecord

// InterviewRecord 面试记录模型
type (
	_InterviewRecord struct {
	}
	InterviewRecord struct {
		ID                 uint64    `json:"id" gorm:"primaryKey;autoIncrement;comment:记录唯一ID"`
		UserID             uint64    `json:"user_id" gorm:"not null;index:idx_user_id;comment:关联的用户ID"`
		InterviewType      string    `json:"interview_type" gorm:"size:50;not null;comment:面试类型 (例如: 综合面试, 技术一面)"`
		ResumeName         string    `json:"resume_name" gorm:"size:255;comment:使用的简历文件名"`
		Difficulty         string    `json:"difficulty" gorm:"size:50;not null;comment:面试难度 (例如: 简单, 中等, 挑战)"`
		CompanyName        string    `json:"company_name" gorm:"size:100;not null;index:idx_company_name;comment:公司名称"`
		PositionName       string    `json:"position_name" gorm:"size:255;not null;comment:岗位名称"`
		InterviewTime      time.Time `json:"interview_time" gorm:"not null;index:idx_interview_time;comment:面试时间"`
		Score              uint8     `json:"score" gorm:"not null;index:idx_score;comment:本次面试评分 (0-100)"`
		OverallSummary     string    `json:"overall_summary" gorm:"type:text;comment:面试官综合点评-总结"`
		OverallStrengths   string    `json:"overall_strengths" gorm:"type:text;comment:面试官综合点评-优点列表 "`
		OverallWeaknesses  string    `json:"overall_weaknesses" gorm:"type:text;comment:面试官综合点评-缺点列表 "`
		OverallSuggestions string    `json:"overall_suggestions" gorm:"type:text;comment:面试官综合点评-建议列表"`
		CreatedAt          time.Time `json:"created_at" gorm:"autoCreateTime"`
		UpdatedAt          time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	}
)

// TableName 指定表名
func (i *InterviewRecord) TableName() string {
	return "interview_records"
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
func (i *_InterviewRecord) ListInterviewRecords(userID uint64, page, pageSize int) ([]*InterviewRecord, int64, error) {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	var records []*InterviewRecord
	var total int64

	query := getDB().Where("user_id = ?", userID)

	if err := query.Model(&InterviewRecord{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("interview_time DESC").
		Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// UpdateInterviewRecord 更新面试记录
func (i *_InterviewRecord) UpdateInterviewRecord(record *InterviewRecord) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	return getDB().Model(&InterviewRecord{}).
		Where("id = ?", record.ID).
		Updates(record).Error
}

// DeleteInterviewRecord 删除面试记录
func (i *_InterviewRecord) DeleteInterviewRecord(id uint64) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	return getDB().Where("id = ?", id).Delete(&InterviewRecord{}).Error
}
