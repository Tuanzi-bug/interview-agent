package service

import (
	"context"
	"errors"

	"ai-eino-interview-agent/internal/model"
	"ai-eino-interview-agent/internal/repository"
	"ai-eino-interview-agent/pkg/eino"

	"gorm.io/gorm"
)

// ResumeService 简历服务
type ResumeService struct{}

// NewResumeService 创建简历服务实例
func NewResumeService() *ResumeService {
	return &ResumeService{}
}

// CreateResumeRequest 创建简历请求
type CreateResumeRequest struct {
	Name    string `json:"name" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// CreateResume 创建简历
func (s *ResumeService) CreateResume(userID uint, req CreateResumeRequest) (*model.Resume, error) {
	db := repository.GetDB()

	// 创建简历记录
	resume := model.Resume{
		UserID:  userID,
		Name:    req.Name,
		Content: req.Content,
	}

	if err := db.Create(&resume).Error; err != nil {
		return nil, err
	}

	return &resume, nil
}

// GetUserResumes 获取用户的所有简历
func (s *ResumeService) GetUserResumes(userID uint) ([]model.Resume, error) {
	db := repository.GetDB()

	var resumes []model.Resume
	result := db.Where("user_id = ?", userID).Find(&resumes)
	if result.Error != nil {
		return nil, result.Error
	}

	return resumes, nil
}

// GetResumeByID 根据ID获取简历
func (s *ResumeService) GetResumeByID(resumeID, userID uint) (*model.Resume, error) {
	db := repository.GetDB()

	var resume model.Resume
	result := db.Where("id = ? AND user_id = ?", resumeID, userID).First(&resume)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("简历不存在")
		}
		return nil, result.Error
	}

	return &resume, nil
}

// UpdateResume 更新简历
func (s *ResumeService) UpdateResume(resumeID, userID uint, req CreateResumeRequest) (*model.Resume, error) {
	db := repository.GetDB()

	// 查找简历并验证所有权
	var resume model.Resume
	if err := db.Where("id = ? AND user_id = ?", resumeID, userID).First(&resume).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("简历不存在或无权限")
		}
		return nil, err
	}

	// 更新简历信息
	updates := map[string]interface{}{
		"name":    req.Name,
		"content": req.Content,
	}

	if err := db.Model(&resume).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 清除旧的分析结果
	if resume.Analysis != "" {
		if err := db.Model(&resume).Update("analysis", "").Error; err != nil {
			return nil, err
		}
	}

	// 重新加载简历信息
	if err := db.First(&resume, resumeID).Error; err != nil {
		return nil, err
	}

	return &resume, nil
}

// DeleteResume 删除简历
func (s *ResumeService) DeleteResume(resumeID, userID uint) error {
	db := repository.GetDB()

	// 验证简历所有权
	var resume model.Resume
	if err := db.Where("id = ? AND user_id = ?", resumeID, userID).First(&resume).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("简历不存在或无权限")
		}
		return err
	}

	// 删除简历
	if err := db.Delete(&resume).Error; err != nil {
		return err
	}

	return nil
}

// AnalyzeResume 分析简历
func (s *ResumeService) AnalyzeResume(resumeID, userID uint) (*model.Resume, error) {
	db := repository.GetDB()

	// 获取简历
	resume, err := s.GetResumeByID(resumeID, userID)
	if err != nil {
		return nil, err
	}

	// 使用Eino分析简历
	ctx := context.Background()
	analysis, err := eino.AnalyzeResume(ctx, resume.Content)
	if err != nil {
		return nil, err
	}

	// 更新分析结果
	if err := db.Model(resume).Update("analysis", analysis).Error; err != nil {
		return nil, err
	}

	// 重新加载简历信息
	resume.Analysis = analysis

	return resume, nil
}