package impl

import (
	interviewsapi "ai-eino-interview-agent/api/model/interviews"
	"ai-eino-interview-agent/chatApp/agent/ext"
	"ai-eino-interview-agent/internal/model"
	"context"
	"log"
	"time"
)

// InterviewServiceImpl 面试服务实现
type InterviewServiceImpl struct{}

// NewInterviewServiceImpl 创建面试服务实例
func NewInterviewServiceImpl() *InterviewServiceImpl {
	return &InterviewServiceImpl{}
}

// StartInterviewStream 启动面试流程（流式）
func (s *InterviewServiceImpl) StartInterviewStream(ctx context.Context, req *interviewsapi.StartInterviewRequest, maxQuestions int) (<-chan *interviewsapi.InterviewEvent, error) {
	// 调用 chatApp 中的流式接口
	eventChan, err := ext.StartInterviewStream(ctx, req.Query, maxQuestions)
	if err != nil {
		return nil, err
	}

	// 转换事件类型：从 ext.InterviewEvent 转换为 interviewsapi.InterviewEvent
	apiEventChan := make(chan *interviewsapi.InterviewEvent, 10)
	go func() {
		defer close(apiEventChan)
		for event := range eventChan {
			apiEvent := convertToAPIEvent(event)
			select {
			case apiEventChan <- apiEvent:
			case <-ctx.Done():
				return
			}
		}
	}()

	return apiEventChan, nil
}

// ContinueInterview 继续面试流程（用于多轮对话）
func (s *InterviewServiceImpl) ContinueInterview(ctx context.Context, req *interviewsapi.ContinueInterviewRequest, maxQuestions int) (<-chan *interviewsapi.InterviewEvent, error) {
	// 调用 chatApp 中的继续面试接口
	eventChan, err := ext.ContinueInterview(ctx, req.Query, maxQuestions)
	if err != nil {
		return nil, err
	}

	// 转换事件类型：从 ext.InterviewEvent 转换为 interviewsapi.InterviewEvent
	apiEventChan := make(chan *interviewsapi.InterviewEvent, 10)
	go func() {
		defer close(apiEventChan)
		for event := range eventChan {
			apiEvent := convertToAPIEvent(event)
			select {
			case apiEventChan <- apiEvent:
			case <-ctx.Done():
				return
			}
		}
	}()

	return apiEventChan, nil
}

// SaveInterviewRecord 保存面试记录
func (s *InterviewServiceImpl) SaveInterviewRecord(ctx context.Context, userID uint, title, query string) (uint64, error) {
	now := time.Now()
	record := &model.InterviewRecord{
		UserID:         userID,
		Title:          title,
		Query:          query,
		Status:         "pending",
		Messages:       "[]",
		LastModifiedAt: now,
	}
	err := model.InterviewRecordDao.CreateInterviewRecord(record)
	if err != nil {
		return 0, err
	}
	return record.ID, nil
}

// UpdateInterviewRecord 更新面试记录（用于保存对话历史和状态）
func (s *InterviewServiceImpl) UpdateInterviewRecord(ctx context.Context, recordID uint64, messages string, status, currentAgent string) error {
	record := &model.InterviewRecord{
		ID:           recordID,
		Messages:     messages,
		Status:       status,
		CurrentAgent: currentAgent,
	}
	return model.InterviewRecordDao.UpdateInterviewRecord(record)
}

// CompleteInterviewRecord 完成面试记录（保存最终报告和评分）
func (s *InterviewServiceImpl) CompleteInterviewRecord(ctx context.Context, recordID uint64, report string, duration int64, score *float64) error {
	return model.InterviewRecordDao.CompleteInterviewRecord(recordID, report, duration, score)
}

// SaveInterviewDialogues 保存面试对话和问题主题
func (s *InterviewServiceImpl) SaveInterviewDialogues(ctx context.Context, userID uint, recordID uint64, questions []interface{}, dialogues []interface{}) error {
	log.Printf("[SaveInterviewDialogues] 开始保存，问题数: %d, 对话数: %d", len(questions), len(dialogues))

	// 保存问题主题
	for i, q := range questions {
		qData, ok := q.(map[string]interface{})
		if !ok {
			continue
		}

		topic := &model.InterviewQuestionTopic{
			UserID:        userID,
			ReportID:      recordID,
			QuestionText:  toString(qData["question_text"]),
			DisplayOrder:  uint32(i + 1),
			EvalDimension: toString(qData["eval_dimension"]),
		}

		if err := model.InterviewQuestionTopicDao.Create(topic); err != nil {
			return err
		}

		// 保存对应的对话记录 - 将提问和回答配对保存在一条记录中
		dialogueMap := make(map[uint32]*model.InterviewDialogue)
		dialogueCounter := uint32(1) // 用于自动分配 display_order

		log.Printf("[SaveInterviewDialogues] 问题 %d: 开始收集对话", i+1)

		// 第一遍：遍历所有对话，收集属于当前问题的对话
		for _, d := range dialogues {
			dData, ok := d.(map[string]interface{})
			if !ok {
				continue
			}

			displayOrder := toUint32(dData["display_order"])
			speakerType := toString(dData["speaker_type"])
			content := toString(dData["content"])

			// 如果 display_order 为 0，自动分配
			if displayOrder == 0 {
				// 自动分配 display_order：基于问题索引和对话计数
				// 例如：第1个问题（i=0）的对话 display_order 为 1, 2, 3...
				//      第2个问题（i=1）的对话 display_order 为 101, 102, 103...
				// 注意：这里我们只在 interviewer 类型时增加计数，这样 interviewer 和 candidate 会共享同一个 displayOrder
				if speakerType == "interviewer" {
					displayOrder = uint32(i)*100 + dialogueCounter
					log.Printf("[SaveInterviewDialogues] 对话自动分配 displayOrder=%d, speaker_type=%s", displayOrder, speakerType)
				} else {
					// candidate 回答使用前一个 interviewer 的 displayOrder
					if dialogueCounter > 1 {
						displayOrder = uint32(i)*100 + (dialogueCounter - 1)
					} else {
						displayOrder = uint32(i)*100 + dialogueCounter
					}
					log.Printf("[SaveInterviewDialogues] 对话自动分配 displayOrder=%d, speaker_type=%s", displayOrder, speakerType)
				}
			} else {
				// 检查对话是否属于当前问题
				minDisplayOrder := uint32(i) * 100
				maxDisplayOrder := uint32(i)*100 + 199
				if displayOrder <= minDisplayOrder || displayOrder > maxDisplayOrder {
					log.Printf("[SaveInterviewDialogues] 对话被过滤: displayOrder=%d (超出范围), speaker_type=%s", displayOrder, speakerType)
					continue
				}
				log.Printf("[SaveInterviewDialogues] 对话被收集: displayOrder=%d, speaker_type=%s", displayOrder, speakerType)
			}

			// 如果该 displayOrder 的记录不存在，创建新记录
			if _, exists := dialogueMap[displayOrder]; !exists {
				dialogueMap[displayOrder] = &model.InterviewDialogue{
					UserID:       userID,
					TopicID:      topic.ID,
					Question:     "",
					Answer:       "",
					DisplayOrder: displayOrder,
				}
			}

			// 根据发言人类型填充对应字段
			if speakerType == "interviewer" {
				dialogueMap[displayOrder].Question = content
				// 只在 interviewer 时增加计数，这样下一个 candidate 会使用相同的 displayOrder
				if displayOrder == uint32(i)*100+dialogueCounter {
					dialogueCounter++
				}
			} else if speakerType == "candidate" {
				dialogueMap[displayOrder].Answer = content
			}
		}

		// 第二遍：保存所有对话记录到数据库
		log.Printf("[SaveInterviewDialogues] 问题 %d: 收集到 %d 条对话", i+1, len(dialogueMap))
		for displayOrder, dialogue := range dialogueMap {
			// 只保存有内容的记录（至少有提问或回答）
			if dialogue.Question != "" || dialogue.Answer != "" {
				log.Printf("[SaveInterviewDialogues] 保存对话 displayOrder=%d, 问题='%s...', 回答='%s...'",
					displayOrder,
					truncateString(dialogue.Question, 30),
					truncateString(dialogue.Answer, 30))
				if err := model.InterviewDialogueDao.Create(dialogue); err != nil {
					log.Printf("[SaveInterviewDialogues] 保存失败: %v", err)
					return err
				}
			}
		}
	}

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
	if f, ok := v.(float64); ok {
		return uint32(f)
	}
	if i, ok := v.(int); ok {
		return uint32(i)
	}
	return 0
}

func (s *InterviewServiceImpl) ListInterviewRecords(ctx context.Context, userID uint, page, pageSize int) ([]*interviewsapi.InterviewRecordDTO, int64, error) {
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
