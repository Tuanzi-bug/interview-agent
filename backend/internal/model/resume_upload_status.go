package model

import (
	"time"
)

// ResumeUploadStatus 简历上传状态跟踪
type ResumeUploadStatus struct {
	ID uint64 `gorm:"primarykey"`

	// 关联信息
	UploadID string `gorm:"uniqueIndex;type:varchar(36);not null;comment:上传任务ID(UUID)"`
	UserID   uint   `gorm:"index;not null;comment:用户ID"`
	FilePath string `gorm:"type:varchar(500);not null;comment:文件路径"`
	ResumeID uint64 `gorm:"index;comment:解析完成后的简历ID"`

	// 状态字段
	Status   string `gorm:"type:varchar(20);not null;default:'pending';comment:状态:pending,extracting,analyzing,completed,failed"`
	Progress int    `gorm:"default:0;comment:进度百分比(0-100)"`
	Stage    string `gorm:"type:varchar(100);comment:当前阶段描述"`
	ErrorMsg string `gorm:"type:text;comment:错误信息"`

	// 性能指标
	ExtractDuration int64 `gorm:"comment:PDF提取耗时(ms)"`
	AnalyzeDuration int64 `gorm:"comment:AI分析耗时(ms)"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ResumeUploadStatus) TableName() string {
	return "resume_upload_status"
}

// ResumeUploadStatusDao 简历上传状态DAO
var ResumeUploadStatusDao = &resumeUploadStatusDao{}

type resumeUploadStatusDao struct{}

// CreateUploadStatus 创建上传状态记录
func (d *resumeUploadStatusDao) CreateUploadStatus(status *ResumeUploadStatus) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	return getDB().Create(status).Error
}

// GetByUploadID 根据UploadID获取状态
func (d *resumeUploadStatusDao) GetByUploadID(uploadID string) (*ResumeUploadStatus, error) {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	var status ResumeUploadStatus
	err := getDB().Where("upload_id = ?", uploadID).First(&status).Error
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// UpdateStatus 更新状态
func (d *resumeUploadStatusDao) UpdateStatus(uploadID string, updates map[string]interface{}) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	return getDB().Model(&ResumeUploadStatus{}).
		Where("upload_id = ?", uploadID).
		Updates(updates).Error
}
