package impl

import (
	interviewsapi "ai-eino-interview-agent/api/model/interviews"
	"ai-eino-interview-agent/internal/model"
	"context"
	"log"
	"strconv"
	"strings"
)

// InterviewServiceImpl 面试服务实现
type InterviewServiceImpl struct{}

// NewInterviewServiceImpl 创建面试服务实例
func NewInterviewServiceImpl() *InterviewServiceImpl {
	return &InterviewServiceImpl{}
}

// CreateInterviewRecord 创建面试记录，返回记录ID
func (s *InterviewServiceImpl) CreateInterviewRecord(ctx context.Context, dto *interviewsapi.InterviewRecordDTO) (uint64, error) {
	// 处理指针类型字段，提取值或使用默认值
	companyName := ""
	if dto.CompanyName != nil {
		companyName = *dto.CompanyName
	}

	positionName := ""
	if dto.PositionName != nil {
		positionName = *dto.PositionName
	}

	interviewDuration := ""
	if dto.InterviewDuration != nil {
		interviewDuration = *dto.InterviewDuration
	}

	// 初始化状态为pending（如果未提供）
	status := dto.Status
	if status == "" {
		status = "pending"
	}

	var duration int64 = 0
	if dto.Duration != nil {
		duration = *dto.Duration
	}

	record := &model.InterviewRecord{
		UserID:            uint(dto.UserID),
		Title:             dto.Title,
		Type:              dto.Type,
		Difficulty:        dto.Difficulty,
		CompanyName:       companyName,
		PositionName:      positionName,
		InterviewDuration: interviewDuration,
		Status:            status,
		Duration:          duration,
	}

	recordID, err := model.InterviewRecordDao.CreateInterviewRecord(record)
	if err != nil {
		log.Printf("[CreateInterviewRecord] 创建面试记录失败: %v", err)
		return 0, err
	}

	log.Printf("[CreateInterviewRecord] 面试记录创建成功，ID: %d，用户ID: %d，标题: %s", recordID, dto.UserID, dto.Title)
	return recordID, nil
}

// UpdateInterviewRecord 更新面试记录
func (s *InterviewServiceImpl) UpdateInterviewRecord(ctx context.Context, dto *interviewsapi.InterviewRecordDTO) error {
	// 处理指针类型字段，提取值或使用默认值
	companyName := ""
	if dto.CompanyName != nil {
		companyName = *dto.CompanyName
	}

	positionName := ""
	if dto.PositionName != nil {
		positionName = *dto.PositionName
	}

	interviewDuration := ""
	if dto.InterviewDuration != nil {
		interviewDuration = *dto.InterviewDuration
	}

	var duration int64 = 0
	if dto.Duration != nil {
		duration = *dto.Duration
	}

	record := &model.InterviewRecord{
		ID:                uint64(dto.ID),
		UserID:            uint(dto.UserID),
		Title:             dto.Title,
		Type:              dto.Type,
		Difficulty:        dto.Difficulty,
		CompanyName:       companyName,
		PositionName:      positionName,
		InterviewDuration: interviewDuration,
		Status:            dto.Status,
		Duration:          duration,
	}

	err := model.InterviewRecordDao.UpdateInterviewRecord(record)
	if err != nil {
		log.Printf("[UpdateInterviewRecord] 更新面试记录失败: %v", err)
		return err
	}

	log.Printf("[UpdateInterviewRecord] 面试记录更新成功，ID: %d，用户ID: %d，标题: %s", dto.ID, dto.UserID, dto.Title)
	return nil
}

// ListInterviewRecords 获取面试记录列表
func (s *InterviewServiceImpl) ListInterviewRecords(ctx context.Context, userID uint, page, pageSize *int32) ([]*interviewsapi.InterviewRecordDTO, int64, error) {
	// 设置默认分页参数
	pageNum := int32(1)
	pageSz := int32(10)

	if page != nil && *page > 0 {
		pageNum = *page
	}
	if pageSize != nil && *pageSize > 0 {
		pageSz = *pageSize
	}

	// 调用 DAO 层获取数据
	records, total, err := model.InterviewRecordDao.ListInterviewRecords(userID, &pageNum, &pageSz)
	if err != nil {
		log.Printf("[ListInterviewRecords] 查询面试记录失败: %v", err)
		return nil, 0, err
	}

	// 转换为 DTO
	dtoList := make([]*interviewsapi.InterviewRecordDTO, 0, len(records))
	for _, record := range records {
		dto := convertToInterviewRecordDTO(record)
		dtoList = append(dtoList, dto)
	}

	log.Printf("[ListInterviewRecords] 查询成功，用户ID: %d，总数: %d，页码: %d，每页: %d", userID, total, pageNum, pageSz)
	return dtoList, total, nil
}

// convertToInterviewRecordDTO 将 model.InterviewRecord 转换为 interviewsapi.InterviewRecordDTO
func convertToInterviewRecordDTO(record *model.InterviewRecord) *interviewsapi.InterviewRecordDTO {
	dto := interviewsapi.NewInterviewRecordDTO()
	dto.ID = int64(record.ID)
	dto.UserID = int32(record.UserID)
	dto.Title = record.Title
	dto.Type = record.Type
	dto.Difficulty = record.Difficulty
	dto.Status = record.Status

	if record.CompanyName != "" {
		v := record.CompanyName
		dto.CompanyName = &v
	}
	if record.PositionName != "" {
		v := record.PositionName
		dto.PositionName = &v
	}
	if record.InterviewDuration != "" {
		v := record.InterviewDuration
		dto.InterviewDuration = &v
	}
	if record.Duration != 0 {
		v := record.Duration
		dto.Duration = &v
	}
	if !record.CreatedAt.IsZero() {
		ms := record.CreatedAt.UnixNano() / int64(1000000)
		dto.CreatedAt = &ms
	}
	if !record.UpdatedAt.IsZero() {
		ms := record.UpdatedAt.UnixNano() / int64(1000000)
		dto.UpdatedAt = &ms
	}

	return dto
}

// SaveInterviewDialogues 保存面试对话和问题主题
func (s *InterviewServiceImpl) SaveInterviewDialogues(ctx context.Context, userID uint, recordID uint64, questions []interface{}, dialogues []interface{}) error {
	log.Printf("[SaveInterviewDialogues] 开始保存，问题数: %d, 对话数: %d", len(questions), len(dialogues))

	// 第一步：按 eval_dimension 去重，只创建 6 个主题
	// 同时建立 questionIndex 到 eval_dimension 的映射
	dimensionTopicMap := make(map[string]*model.InterviewQuestionTopic)
	questionIndexToDimension := make(map[uint32]string)             // questionIndex -> eval_dimension
	dimensionQuestionMap := make(map[string]map[string]interface{}) // 存储每个维度的第一个问题
	dimensionOrder := make(map[string]uint32)
	orderCounter := uint32(1)

	// 第一遍：找到每个维度的第一个问题，并建立 questionIndex 映射
	for _, q := range questions {
		qData, ok := q.(map[string]interface{})
		if !ok {
			continue
		}

		// 获取 order 字段作为 questionIndex
		order := toUint32(qData["order"])
		if order == 0 {
			continue
		}

		// 获取 eval_dimension，如果包含多个维度（用|分隔），只取第一个
		evalDim := toString(qData["eval_dimension"])
		if idx := strings.Index(evalDim, "|"); idx >= 0 {
			evalDim = evalDim[:idx]
		}

		// 记录 questionIndex 到 eval_dimension 的映射
		if _, exists := questionIndexToDimension[order]; !exists {
			questionIndexToDimension[order] = evalDim
			log.Printf("[SaveInterviewDialogues] 映射 questionIndex=%d -> eval_dimension=%s", order, evalDim)
		}

		// 如果这个维度还没有记录过，就记录这个问题
		if _, exists := dimensionQuestionMap[evalDim]; !exists {
			dimensionQuestionMap[evalDim] = qData
			log.Printf("[SaveInterviewDialogues] 找到维度 %s 的第一个问题: '%s...'", evalDim, truncateString(toString(qData["question_text"]), 30))
		}
	}

	// 第二遍：为每个维度创建主题
	for evalDim, qData := range dimensionQuestionMap {
		topic := &model.InterviewQuestionTopic{
			UserID:        userID,
			ReportID:      recordID,
			QuestionText:  toString(qData["question_text"]),
			DisplayOrder:  orderCounter,
			EvalDimension: evalDim,
		}

		if err := model.InterviewQuestionTopicDao.Create(topic); err != nil {
			log.Printf("[SaveInterviewDialogues] 创建主题失败: %v", err)
			return err
		}

		dimensionTopicMap[evalDim] = topic
		dimensionOrder[evalDim] = orderCounter
		log.Printf("[SaveInterviewDialogues] 创建维度主题 %d: %s (ID: %d)", orderCounter, evalDim, topic.ID)
		orderCounter++
	}

	// 第二步：保存对话记录 - 一次性处理所有对话，避免重复
	dialogueMap := make(map[uint32]*model.InterviewDialogue)

	for i, d := range dialogues {
		dData, ok := d.(map[string]interface{})
		if !ok {
			log.Printf("[SaveInterviewDialogues] 对话 %d: 类型转换失败，类型=%T", i, d)
			continue
		}

		displayOrder := toUint32(dData["display_order"])
		speakerType := toString(dData["speaker_type"])
		content := toString(dData["content"])

		log.Printf("[SaveInterviewDialogues] 对话 %d: displayOrder=%d, speakerType=%s, content='%s...'", i, displayOrder, speakerType, truncateString(content, 20))

		if displayOrder == 0 || content == "" {
			log.Printf("[SaveInterviewDialogues] 对话 %d 被过滤: displayOrder=%d, content长度=%d", i, displayOrder, len(content))
			continue
		}

		// 根据 displayOrder 确定属于哪个问题（displayOrder 的百位数字是问题索引）
		questionIndex := displayOrder / 100
		if questionIndex == 0 {
			log.Printf("[SaveInterviewDialogues] 对话被过滤: displayOrder=%d (questionIndex=0)", displayOrder)
			continue
		}

		// 从映射中获取该 questionIndex 对应的 eval_dimension
		evalDim, exists := questionIndexToDimension[questionIndex]
		if !exists {
			log.Printf("[SaveInterviewDialogues] 对话被过滤: displayOrder=%d (questionIndex=%d 没有对应的维度映射)", displayOrder, questionIndex)
			continue
		}

		log.Printf("[SaveInterviewDialogues] 对话 %d: displayOrder=%d, questionIndex=%d, evalDim=%s", i, displayOrder, questionIndex, evalDim)

		// 获取该维度的主题
		topic, ok := dimensionTopicMap[evalDim]
		if !ok {
			log.Printf("[SaveInterviewDialogues] 找不到维度主题: %s", evalDim)
			continue
		}
		log.Printf("[SaveInterviewDialogues] 分配 topic_id=%d 给 displayOrder=%d", topic.ID, displayOrder)

		// 如果该 displayOrder 的记录不存在，创建新记录
		if _, exists := dialogueMap[displayOrder]; !exists {
			dialogueMap[displayOrder] = &model.InterviewDialogue{
				UserID:       userID,
				TopicID:      topic.ID,
				ReportID:     recordID,
				Question:     "",
				Answer:       "",
				DisplayOrder: displayOrder,
			}
		}

		// 根据发言人类型填充对应字段
		if speakerType == "interviewer" {
			// 只保存第一个提问，避免重复
			if dialogueMap[displayOrder].Question == "" {
				dialogueMap[displayOrder].Question = content
				log.Printf("[SaveInterviewDialogues] 保存提问 displayOrder=%d, 内容='%s...'", displayOrder, truncateString(content, 30))
			}
		} else if speakerType == "candidate" {
			// 只保存第一个回答，避免重复
			if dialogueMap[displayOrder].Answer == "" {
				dialogueMap[displayOrder].Answer = content
				log.Printf("[SaveInterviewDialogues] 保存回答 displayOrder=%d, 内容='%s...'", displayOrder, truncateString(content, 30))
			}
		}
	}

	// 第三步：保存所有对话记录到数据库
	log.Printf("[SaveInterviewDialogues] 准备保存 %d 条对话记录", len(dialogueMap))
	for _, dialogue := range dialogueMap {
		// 只保存有内容的记录（至少有提问或回答）
		if dialogue.Question != "" || dialogue.Answer != "" {
			if err := model.InterviewDialogueDao.Create(dialogue); err != nil {
				log.Printf("[SaveInterviewDialogues] 保存失败: %v", err)
				return err
			}
		}
	}

	log.Printf("[SaveInterviewDialogues] 保存完成，共创建 %d 个维度主题，保存 %d 条对话记录", len(dimensionTopicMap), len(dialogueMap))
	return nil
}

// 辅助函数：截断字符串用于日志输出
func truncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen] + "..."
	}
	return s
}

// 辅助函数：将 interface{} 转换为 string
func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// 辅助函数：将 interface{} 转换为 uint32
func toUint32(v interface{}) uint32 {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return uint32(val)
	case float32:
		return uint32(val)
	case int:
		return uint32(val)
	case int32:
		return uint32(val)
	case int64:
		return uint32(val)
	case uint:
		return uint32(val)
	case uint32:
		return val
	case uint64:
		return uint32(val)
	case string:
		// 尝试解析字符串
		if i, err := strconv.ParseUint(val, 10, 32); err == nil {
			return uint32(i)
		}
	}
	return 0
}

// GetInterviewEvaluation 根据用户ID和报告ID获取面试评估报告
func (s *InterviewServiceImpl) GetInterviewEvaluation(ctx context.Context, userID uint, reportID uint64) (interface{}, error) {
	evaluation, err := model.InterviewEvaluationDao.GetEvaluationByUserIDAndReportID(userID, reportID)
	if err != nil {
		log.Printf("Failed to get interview evaluation: %v", err)
		return nil, err
	}

	// 返回评估数据
	return map[string]interface{}{
		"id":         evaluation.ID,
		"user_id":    evaluation.UserID,
		"report_id":  evaluation.ReportID,
		"comment":    evaluation.Comment,
		"score":      evaluation.Score,
		"dimensions": evaluation.Dimensions,
		"created_at": evaluation.CreatedAt,
		"updated_at": evaluation.UpdatedAt,
	}, nil
}

// GetAnswerReport 根据用户ID和报告ID获取答题报告
func (s *InterviewServiceImpl) GetAnswerReport(ctx context.Context, userID uint, reportID uint64) (interface{}, error) {
	report, err := model.AnswerReportDao.GetAnswerReportByUserIDAndReportID(userID, reportID)
	if err != nil {
		log.Printf("[GetAnswerReport] 获取答题报告失败: %v", err)
		return nil, err
	}

	// 返回答题报告数据
	return map[string]interface{}{
		"id":         report.ID,
		"user_id":    report.UserID,
		"report_id":  report.ReportID,
		"records":    report.Records,
		"deleted":    report.Deleted,
		"created_at": report.CreatedAt,
		"updated_at": report.UpdatedAt,
	}, nil
}
