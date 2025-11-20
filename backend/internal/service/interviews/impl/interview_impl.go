package impl

import (
	interviewsapi "ai-eino-interview-agent/api/model/interviews"
	"ai-eino-interview-agent/chatApp/agent/ext"
	"ai-eino-interview-agent/internal/model"
	"context"
	"log"
	"strconv"
	"strings"
	"time"
)

// InterviewServiceImpl 面试服务实现
type InterviewServiceImpl struct{}

// NewInterviewServiceImpl 创建面试服务实例
func NewInterviewServiceImpl() *InterviewServiceImpl {
	return &InterviewServiceImpl{}
}

//
//// StartInterviewStream 启动面试流程（流式）
//func (s *InterviewServiceImpl) StartInterviewStream(ctx context.Context, req *interviewsapi.StartInterviewRequest, maxQuestions int) (<-chan *interviewsapi.InterviewEvent, error) {
//	// 调用 chatApp 中的流式接口
//	eventChan, err := ext.StartInterviewStream(ctx, req.Query, maxQuestions)
//	if err != nil {
//		return nil, err
//	}
//
//	// 转换事件类型：从 ext.InterviewEvent 转换为 interviewsapi.InterviewEvent
//	apiEventChan := make(chan *interviewsapi.InterviewEvent, 10)
//	go func() {
//		defer close(apiEventChan)
//		for event := range eventChan {
//			apiEvent := convertToAPIEvent(event)
//			select {
//			case apiEventChan <- apiEvent:
//			case <-ctx.Done():
//				return
//			}
//		}
//	}()
//
//	return apiEventChan, nil
//}
//
//// ContinueInterview 继续面试流程（用于多轮对话）
//func (s *InterviewServiceImpl) ContinueInterview(ctx context.Context, req *interviewsapi.ContinueInterviewRequest, maxQuestions int) (<-chan *interviewsapi.InterviewEvent, error) {
//	// 调用 chatApp 中的继续面试接口
//	eventChan, err := ext.ContinueInterview(ctx, req.Query, maxQuestions)
//	if err != nil {
//		return nil, err
//	}
//
//	// 转换事件类型：从 ext.InterviewEvent 转换为 interviewsapi.InterviewEvent
//	apiEventChan := make(chan *interviewsapi.InterviewEvent, 10)
//	go func() {
//		defer close(apiEventChan)
//		for event := range eventChan {
//			apiEvent := convertToAPIEvent(event)
//			select {
//			case apiEventChan <- apiEvent:
//			case <-ctx.Done():
//		Query:          query,
//		Status:         "pending",
//		Messages:       "[]",
//		LastModifiedAt: now,
//	}
//	err := model.InterviewRecordDao.CreateInterviewRecord(record)
//	if err != nil {
//		return 0, err
//	}
//	return record.ID, nil
//}
//
//// UpdateInterviewRecord 更新面试记录（用于保存对话历史和状态）
//func (s *InterviewServiceImpl) UpdateInterviewRecord(ctx context.Context, recordID uint64, messages string, status, currentAgent string) error {
//	record := &model.InterviewRecord{
//		ID:           recordID,
//		Messages:     messages,
//		Status:       status,
//		CurrentAgent: currentAgent,
//	}
//	return model.InterviewRecordDao.UpdateInterviewRecord(record)
//}
//
//// CompleteInterviewRecord 完成面试记录（保存最终报告和评分）
//func (s *InterviewServiceImpl) CompleteInterviewRecord(ctx context.Context, recordID uint64, report string, duration int64, score *float64) error {
//	return model.InterviewRecordDao.CompleteInterviewRecord(recordID, report, duration, score)
//}

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

// 辅助函数：将 interface{} 转换为 float64
func toFloat64(v interface{}) float64 {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case uint:
		return float64(val)
	case uint32:
		return float64(val)
	case uint64:
		return float64(val)
	case string:
		// 尝试解析字符串
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return 0
}

func (s *InterviewServiceImpl) ListInterviewRecords(ctx context.Context, userID uint, page, pageSize *int32) ([]*interviewsapi.InterviewRecordDTO, int64, error) {
	records, total, err := model.InterviewRecordDao.ListInterviewRecords(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*interviewsapi.InterviewRecordDTO, 0, len(records))
	for _, r := range records {
		result = append(result, convertToInterviewRecordDTO(r))
	}
	return result, total, nil
}

func (s *InterviewServiceImpl) GetInterviewRecord(ctx context.Context, userID uint, recordID uint64) (*interviewsapi.InterviewRecordDTO, error) {
	record, err := model.InterviewRecordDao.GetInterviewRecordByID(recordID)
	if err != nil {
		return nil, err
	}
	if record.UserID != userID {
		return nil, nil
	}
	return convertToInterviewRecordDTO(record), nil
}

// convertToAPIEvent 将 ext.InterviewEvent 转换为 interviewsapi.InterviewEvent
func convertToAPIEvent(event *ext.InterviewEvent) *interviewsapi.InterviewEvent {
	apiEvent := interviewsapi.NewInterviewEvent()
	apiEvent.Type = event.Type

	if event.AgentName != "" {
		apiEvent.AgentName = &event.AgentName
	}
	if event.Message != "" {
		apiEvent.Message = &event.Message
	}
	if event.TransferTo != "" {
		apiEvent.TransferTo = &event.TransferTo
	}
	if event.Error != "" {
		apiEvent.Error = &event.Error
	}
	if event.Status != nil && *event.Status != "" {
		apiEvent.Status = event.Status
	}
	if event.Report != "" {
		apiEvent.Report = &event.Report
	}
	if event.Score != nil {
		apiEvent.Score = event.Score
	}
	if event.Duration > 0 {
		apiEvent.Duration = &event.Duration
	}
	if event.Feedback != "" {
		apiEvent.Feedback = &event.Feedback
	}
	if event.Messages != "" {
		apiEvent.Messages = &event.Messages
	}

	return apiEvent
}

func convertToInterviewRecordDTO(record *model.InterviewRecord) *interviewsapi.InterviewRecordDTO {
	dto := interviewsapi.NewInterviewRecordDTO()
	dto.ID = int64(record.ID)
	dto.UserID = int32(record.UserID)
	dto.Title = record.Title
	dto.Query = record.Query
	dto.Status = record.Status
	if record.Messages != "" {
		v := record.Messages
		dto.Messages = &v
	}
	if record.Report != "" {
		v := record.Report
		dto.Report = &v
	}
	if record.CurrentAgent != "" {
		v := record.CurrentAgent
		dto.CurrentAgent = &v
	}
	if record.Duration != 0 {
		v := record.Duration
		dto.Duration = &v
	}
	if record.Score != nil {
		v := *record.Score
		dto.Score = &v
	}
	if record.Feedback != "" {
		v := record.Feedback
		dto.Feedback = &v
	}
	if !record.CreatedAt.IsZero() {
		ms := record.CreatedAt.UnixNano() / int64(time.Millisecond)
		dto.CreatedAt = &ms
	}
	if !record.UpdatedAt.IsZero() {
		ms := record.UpdatedAt.UnixNano() / int64(time.Millisecond)
		dto.UpdatedAt = &ms
	}
	if record.CompletedAt != nil {
		ms := record.CompletedAt.UnixNano() / int64(time.Millisecond)
		dto.CompletedAt = &ms
	}
	return dto
}

//// SaveInterviewEvaluation 保存面试评估数据
//func (s *InterviewServiceImpl) SaveInterviewEvaluation(ctx context.Context, userID uint, reportID uint64, comment string, score float64, dimensions interface{}) error {
//	// 将 dimensions 转换为 []*model.EvaluationDimension
//	var dimensionList []*model.EvaluationDimension
//
//	// 如果 dimensions 是 []interface{}，则转换为 []*model.EvaluationDimension
//	if dimSlice, ok := dimensions.([]interface{}); ok {
//		for _, dim := range dimSlice {
//			if dimMap, ok := dim.(map[string]interface{}); ok {
//				evalDim := &model.EvaluationDimension{
//					DimensionName: toString(dimMap["dimension_name"]),
//					Evaluation:    toString(dimMap["evaluation"]),
//					Score:         toFloat64(dimMap["score"]),
//				}
//				dimensionList = append(dimensionList, evalDim)
//			}
//		}
//	}
//
//	// 创建评估记录
//	evaluation := &model.InterviewEvaluation{
//		UserID:     userID,
//		ReportID:   reportID,
//		Comment:    comment,
//		Score:      score,
//		Dimensions: dimensionList,
//		Deleted:    0,
//	}
//
//	// 保存到数据库
//	err := model.InterviewEvaluationDao.CreateEvaluation(evaluation)
//	if err != nil {
//		log.Printf("Failed to save interview evaluation: %v", err)
//		return err
//	}
//
//	return nil
//}

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
