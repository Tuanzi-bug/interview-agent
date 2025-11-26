package impl

import (
	"ai-eino-interview-agent/internal/model"
	"context"
	"log"
)

// ResumeServer 简历管理服务实现
type ResumeServer struct{}

// NewResumeServer 创建简历管理服务实例
func NewResumeServer() *ResumeServer {
	return &ResumeServer{}
}

// UploadResume 上传简历，返回简历ID
func (s *ResumeServer) UploadResume(
	ctx context.Context,
	userID uint,
	fileName string,
	fileType string,
	fileSize int64,
	content string,
) (uint64, error) {
	resume := &model.Resume{
		UserID:    userID,
		Content:   content,
		FileName:  fileName,
		FileSize:  fileSize,
		FileType:  fileType,
		IsDefault: 0,
		Deleted:   0,
	}

	resumeID, err := model.ResumeDao.CreateResume(resume)
	if err != nil {
		log.Printf("[UploadResume] 创建简历失败: %v", err)
		return 0, err
	}

	log.Printf("[UploadResume] 简历上传成功: userID=%d, resumeID=%d, fileName=%s", userID, resumeID, fileName)
	return resumeID, nil
}

// GetResumeByID 根据简历ID获取简历详情
func (s *ResumeServer) GetResumeByID(
	ctx context.Context,
	resumeID uint64,
) (interface{}, error) {
	resume, err := model.ResumeDao.GetResumeByID(resumeID)
	if err != nil {
		log.Printf("[GetResumeByID] 获取简历失败: %v", err)
		return nil, err
	}

	return map[string]interface{}{
		"id":         resume.ID,
		"user_id":    resume.UserID,
		"content":    resume.Content,
		"file_name":  resume.FileName,
		"file_size":  resume.FileSize,
		"file_type":  resume.FileType,
		"is_default": resume.IsDefault,
		"created_at": resume.CreatedAt,
		"updated_at": resume.UpdatedAt,
	}, nil
}

// GetUserResumes 获取用户的所有简历列表
func (s *ResumeServer) GetUserResumes(
	ctx context.Context,
	userID uint,
) (interface{}, error) {
	resumes, err := model.ResumeDao.GetResumeByUserID(userID)
	if err != nil {
		log.Printf("[GetUserResumes] 获取用户简历列表失败: %v", err)
		return nil, err
	}

	var resumeList []map[string]interface{}
	for _, resume := range resumes {
		resumeList = append(resumeList, map[string]interface{}{
			"id":         resume.ID,
			"user_id":    resume.UserID,
			"file_name":  resume.FileName,
			"file_size":  resume.FileSize,
			"file_type":  resume.FileType,
			"is_default": resume.IsDefault,
			"created_at": resume.CreatedAt,
			"updated_at": resume.UpdatedAt,
		})
	}

	return map[string]interface{}{
		"resumes": resumeList,
		"count":   len(resumeList),
	}, nil
}

// GetDefaultResume 获取用户的默认简历
func (s *ResumeServer) GetDefaultResume(
	ctx context.Context,
	userID uint,
) (interface{}, error) {
	resume, err := model.ResumeDao.GetDefaultResumeByUserID(userID)
	if err != nil {
		log.Printf("[GetDefaultResume] 获取默认简历失败: %v", err)
		return nil, err
	}

	return map[string]interface{}{
		"id":         resume.ID,
		"user_id":    resume.UserID,
		"content":    resume.Content,
		"file_name":  resume.FileName,
		"file_size":  resume.FileSize,
		"file_type":  resume.FileType,
		"is_default": resume.IsDefault,
		"created_at": resume.CreatedAt,
		"updated_at": resume.UpdatedAt,
	}, nil
}

// SetDefaultResume 设置用户的默认简历
func (s *ResumeServer) SetDefaultResume(
	ctx context.Context,
	userID uint,
	resumeID uint64,
) error {
	err := model.ResumeDao.SetDefaultResume(userID, resumeID)
	if err != nil {
		log.Printf("[SetDefaultResume] 设置默认简历失败: %v", err)
		return err
	}

	log.Printf("[SetDefaultResume] 默认简历设置成功: userID=%d, resumeID=%d", userID, resumeID)
	return nil
}

// UpdateResume 更新简历信息
func (s *ResumeServer) UpdateResume(
	ctx context.Context,
	resumeID uint64,
	fileName string,
	content string,
) error {
	updates := make(map[string]interface{})
	if fileName != "" {
		updates["file_name"] = fileName
	}
	if content != "" {
		updates["content"] = content
	}

	if len(updates) == 0 {
		return nil
	}

	err := model.ResumeDao.UpdateResume(resumeID, updates)
	if err != nil {
		log.Printf("[UpdateResume] 更新简历失败: %v", err)
		return err
	}

	log.Printf("[UpdateResume] 简历更新成功: resumeID=%d", resumeID)
	return nil
}

// DeleteResume 删除简历
func (s *ResumeServer) DeleteResume(
	ctx context.Context,
	userID uint,
	resumeID uint64,
) error {
	// 验证简历属于该用户
	resume, err := model.ResumeDao.GetResumeByID(resumeID)
	if err != nil {
		log.Printf("[DeleteResume] 获取简历失败: %v", err)
		return err
	}

	if resume.UserID != userID {
		log.Printf("[DeleteResume] 用户无权删除该简历: userID=%d, resumeID=%d", userID, resumeID)
		return err
	}

	err = model.ResumeDao.DeleteResume(resumeID)
	if err != nil {
		log.Printf("[DeleteResume] 删除简历失败: %v", err)
		return err
	}

	log.Printf("[DeleteResume] 简历删除成功: userID=%d, resumeID=%d", userID, resumeID)
	return nil
}

// ListResumesByUserID 分页获取用户的简历列表
func (s *ResumeServer) ListResumesByUserID(
	ctx context.Context,
	userID uint,
	page, pageSize int32,
) (interface{}, int64, error) {
	resumes, total, err := model.ResumeDao.ListResumesByUserID(userID, page, pageSize)
	if err != nil {
		log.Printf("[ListResumesByUserID] 分页获取简历列表失败: %v", err)
		return nil, 0, err
	}

	var resumeList []map[string]interface{}
	for _, resume := range resumes {
		resumeList = append(resumeList, map[string]interface{}{
			"id":         resume.ID,
			"user_id":    resume.UserID,
			"file_name":  resume.FileName,
			"file_size":  resume.FileSize,
			"file_type":  resume.FileType,
			"is_default": resume.IsDefault,
			"created_at": resume.CreatedAt,
			"updated_at": resume.UpdatedAt,
		})
	}

	return map[string]interface{}{
		"resumes":   resumeList,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	}, total, nil
}
