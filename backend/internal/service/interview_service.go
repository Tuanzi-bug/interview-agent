package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"ai-eino-interview-agent/internal/model"
	"ai-eino-interview-agent/internal/repository"
	"ai-eino-interview-agent/pkg/eino"

	"gorm.io/gorm"
)

// InterviewService 面试服务
type InterviewService struct{}

// NewInterviewService 创建面试服务实例
func NewInterviewService() *InterviewService {
	return &InterviewService{}
}

// CreateInterviewRequest 创建面试请求
type CreateInterviewRequest struct {
	ResumeID   uint   `json:"resume_id"`
	JobTitle   string `json:"job_title" binding:"required"`
	Difficulty string `json:"difficulty" binding:"required,oneof=easy medium hard"`
	Type       string `json:"type" binding:"required,oneof=comprehensive specialized resume-based"`
}

// SubmitAnswerRequest 提交回答请求
type SubmitAnswerRequest struct {
	QuestionID uint   `json:"question_id" binding:"required"`
	Answer     string `json:"answer" binding:"required"`
}

// CreateInterview 创建面试
func (s *InterviewService) CreateInterview(userID uint, req CreateInterviewRequest) (*model.Interview, error) {
	db := repository.GetDB()

	// 如果是基于简历的面试，验证简历存在且属于用户
	var resumeContent string
	if req.Type == "resume-based" && req.ResumeID > 0 {
		var resume model.Resume
		if err := db.Where("id = ? AND user_id = ?", req.ResumeID, userID).First(&resume).Error; err != nil {
			return nil, errors.New("简历不存在或无权限")
		}
		resumeContent = resume.Content
	}

	// 创建面试记录
	interview := model.Interview{
		UserID:     userID,
		ResumeID:   req.ResumeID,
		JobTitle:   req.JobTitle,
		Difficulty: req.Difficulty,
		Type:       req.Type,
		Status:     "pending",
	}

	if err := db.Create(&interview).Error; err != nil {
		return nil, err
	}

	// 预生成问题（可选，也可以在开始面试时生成）
	if req.Type == "resume-based" && resumeContent != "" {
		if err := s.generateQuestions(interview.ID, req.JobTitle, req.Difficulty, req.Type, resumeContent); err != nil {
			// 记录错误但不中断流程
			// 可以在开始面试时再次尝试生成
		}
	}

	return &interview, nil
}

// GetUserInterviews 获取用户的所有面试
func (s *InterviewService) GetUserInterviews(userID uint, status string) ([]model.Interview, error) {
	db := repository.GetDB()

	var interviews []model.Interview
	query := db.Where("user_id = ?", userID)
	
	// 如果指定了状态，添加状态过滤
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	result := query.Order("created_at DESC").Find(&interviews)
	if result.Error != nil {
		return nil, result.Error
	}

	return interviews, nil
}

// GetInterviewByID 根据ID获取面试详情
func (s *InterviewService) GetInterviewByID(interviewID, userID uint) (*model.Interview, error) {
	db := repository.GetDB()

	var interview model.Interview
	result := db.Preload("Questions").Where("id = ? AND user_id = ?", interviewID, userID).First(&interview)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("面试不存在")
		}
		return nil, result.Error
	}

	return &interview, nil
}

// StartInterview 开始面试
func (s *InterviewService) StartInterview(interviewID, userID uint) (*model.Interview, error) {
	db := repository.GetDB()

	// 获取面试记录
	interview, err := s.GetInterviewByID(interviewID, userID)
	if err != nil {
		return nil, err
	}

	// 检查状态
	if interview.Status != "pending" {
		return nil, errors.New("只能开始待开始状态的面试")
	}

	// 设置开始时间
	now := time.Now()
	interview.StartTime = &now
	interview.Status = "in_progress"

	// 如果还没有问题，生成问题
	var questions []model.Question
	db.Where("interview_id = ?", interviewID).Find(&questions)
	if len(questions) == 0 {
		// 获取简历内容（如果需要）
		var resumeContent string
		if interview.Type == "resume-based" && interview.ResumeID > 0 {
			var resume model.Resume
			if err := db.First(&resume, interview.ResumeID).Error; err == nil {
				resumeContent = resume.Content
			}
		}

		// 生成问题
		if err := s.generateQuestions(interview.ID, interview.JobTitle, interview.Difficulty, interview.Type, resumeContent); err != nil {
			return nil, err
		}
	}

	// 更新面试状态
	if err := db.Save(interview).Error; err != nil {
		return nil, err
	}

	return interview, nil
}

// EndInterview 结束面试
func (s *InterviewService) EndInterview(interviewID, userID uint) (*model.Interview, error) {
	db := repository.GetDB()

	// 获取面试记录
	interview, err := s.GetInterviewByID(interviewID, userID)
	if err != nil {
		return nil, err
	}

	// 检查状态
	if interview.Status != "in_progress" {
		return nil, errors.New("只能结束进行中的面试")
	}

	// 设置结束时间和计算时长
	now := time.Now()
	interview.EndTime = &now
	if interview.StartTime != nil {
		interview.Duration = int(now.Sub(*interview.StartTime).Seconds())
	}
	interview.Status = "completed"

	// 计算总分和生成评估报告
	if err := s.calculateScoreAndGenerateReport(interview); err != nil {
		// 记录错误但继续流程
	}

	// 更新面试状态
	if err := db.Save(interview).Error; err != nil {
		return nil, err
	}

	return interview, nil
}

// SubmitAnswer 提交回答
func (s *InterviewService) SubmitAnswer(interviewID, userID uint, req SubmitAnswerRequest) (*model.Question, error) {
	db := repository.GetDB()

	// 验证面试所有权和状态
	var interview model.Interview
	if err := db.Where("id = ? AND user_id = ? AND status = ?", interviewID, userID, "in_progress").First(&interview).Error; err != nil {
		return nil, errors.New("面试不存在、无权限或未开始")
	}

	// 验证问题归属
	var question model.Question
	if err := db.Where("id = ? AND interview_id = ?", req.QuestionID, interviewID).First(&question).Error; err != nil {
		return nil, errors.New("问题不存在")
	}

	// 更新回答
	question.Answer = req.Answer

	// 评估回答
	ctx := context.Background()
	feedback, score, err := eino.EvaluateAnswer(ctx, question.Content, req.Answer, interview.JobTitle)
	if err != nil {
		// 记录错误但继续保存回答
	} else {
		question.Feedback = feedback
		question.Score = score
	}

	if err := db.Save(&question).Error; err != nil {
		return nil, err
	}

	return &question, nil
}

// generateQuestions 生成面试问题
func (s *InterviewService) generateQuestions(interviewID uint, jobTitle, difficulty, interviewType, resumeContent string) error {
	db := repository.GetDB()

	// 使用Eino生成问题
	ctx := context.Background()
	questions, err := eino.GenerateInterviewQuestions(ctx, jobTitle, difficulty, interviewType, resumeContent)
	if err != nil {
		return err
	}

	// 保存问题到数据库
	for i, q := range questions {
		question := model.Question{
			InterviewID: interviewID,
			Content:     q,
			Type:        "technical", // 默认类型，可以根据实际问题类型调整
			Order:       i + 1,
		}

		if err := db.Create(&question).Error; err != nil {
			return err
		}
	}

	return nil
}

// calculateScoreAndGenerateReport 计算总分并生成评估报告
func (s *InterviewService) calculateScoreAndGenerateReport(interview *model.Interview) error {
	db := repository.GetDB()

	// 获取所有问题及回答
	var questions []model.Question
	db.Where("interview_id = ?", interview.ID).Find(&questions)

	// 计算总分
	totalScore := 0.0
	scoredCount := 0
	for _, q := range questions {
		if q.Score > 0 {
			totalScore += q.Score
			scoredCount++
		}
	}

	if scoredCount > 0 {
		interview.Score = totalScore / float64(scoredCount)
	}

	// 生成评估报告（这里简化实现，实际可以使用Eino生成更详细的报告）
	interview.Evaluation = "面试评估报告：\n"
	interview.Evaluation += "总分：" + strconv.FormatFloat(interview.Score, 'f', 2, 64) + "\n"
	interview.Evaluation += "面试时长：" + strconv.Itoa(interview.Duration) + "秒\n"
	interview.Evaluation += "问题数量：" + strconv.Itoa(len(questions)) + "个"

	return db.Save(interview).Error
}