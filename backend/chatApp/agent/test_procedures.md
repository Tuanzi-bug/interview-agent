# Go-Eino Interview Agent Test Procedures

## 1. Go Interview Agent Core Functionality

### 1.1 Interview Agent Initialization Testing

#### TC-001: GoInterviewAgent Initialization
**Prerequisites:**
- Go development environment set up
- Required dependencies installed

**Test Steps:**
1. Import the GoInterviewAgent package
2. Create a new instance with valid configuration:
   ```go
   config := &agent.Config{
       ModelName:    "gpt-4",
       MaxQuestions: 5,
       Difficulty:   "medium",
   }
   goAgent, err := NewGoInterviewAgent(config)
   ```
3. Verify no error is returned
4. Verify the agent is not nil
5. Verify configuration parameters are correctly set

#### TC-002: Interface Implementation
**Test Steps:**
1. Create a new GoInterviewAgent instance
2. Verify it implements the InterviewAgent interface using type assertion:
   ```go
   var agentInterface agent.InterviewAgent = goAgent
   ```
3. Verify no compile errors occur

### 1.2 First Question Generation Testing

#### TC-003 to TC-005: First Question Generation
**Prerequisites:**
- Initialized GoInterviewAgent
- Valid resume data for Go candidates

**Test Steps:**
1. Prepare candidate profile with Go experience
2. Generate first question:
   ```go
   question, err := goAgent.GenerateFirstQuestion(candidateProfile)
   ```
3. Verify no error is returned
4. Verify returned QuestionData contains all required fields:
   - QuestionID
   - Content
   - Type
   - Category
   - Difficulty
5. Verify question is relevant to Go programming
6. Repeat with different difficulty levels and verify question complexity adjusts accordingly

### 1.3 Answer Evaluation Testing

#### TC-007 to TC-009: Answer Evaluation
**Prerequisites:**
- Initialized GoInterviewAgent
- Generated question from previous step

**Test Steps:**
1. Prepare correct Go code answer:
   ```go
   correctAnswer := `package main
   
   import "fmt"
   
   func main() {
       fmt.Println("Hello, Go!")
   }
   `
   ```
2. Evaluate correct answer:
   ```go
   evaluation, err := goAgent.EvaluateAnswer(question.QuestionID, correctAnswer)
   ```
3. Verify no error is returned
4. Verify evaluation contains expected fields:
   - Score
   - Feedback
   - Metrics (completeness, correctness, efficiency)
5. Verify correct answers receive high scores and positive feedback
6. Repeat with incorrect or incomplete Go code and verify constructive feedback

## 2. API Integration Testing

### 2.1 Interview Session Management

#### Starting a New Interview
**Test Steps:**
1. Send POST request to `/api/interviews/start`:
   ```http
   POST /api/interviews/start
   Content-Type: application/json
   Authorization: Bearer {token}
   
   {
     "user_id": "test_user_123",
     "resume_id": "resume_456",
     "difficulty": "medium",
     "interview_type": "go_specialized"
   }
   ```
2. Verify HTTP 200 response
3. Extract session ID from response
4. Verify interview record is created in database
5. Verify first question is generated

#### TC-006: Answer Submission
**Test Steps:**
1. Send POST request to `/api/interviews/submit-answer`:
   ```http
   POST /api/interviews/submit-answer
   Content-Type: application/json
   Authorization: Bearer {token}
   
   {
     "session_id": "{session_id}",
     "user_id": "test_user_123",
     "question_id": "{question_id}",
     "answer": "package main\n\nfunc main() {\n    fmt.Println(\"Hello\")\n}"
   }
   ```
2. Verify HTTP 200 response
3. Verify answer record is saved in database
4. Verify evaluation is processed

### 2.2 SSE Event Testing

#### TC-012 to TC-014: SSE Event Verification
**Prerequisites:**
- Active interview session
- SSE client connected to `/api/interviews/events/{session_id}`

**Test Steps:**
1. Establish SSE connection:
   ```http
   GET /api/interviews/events/{session_id}
   Accept: text/event-stream
   Authorization: Bearer {token}
   ```
2. Submit an answer as in TC-006
3. Verify evaluation result event is received:
   ```
   event: evaluation
   data: {"session_id":"{session_id}","score":85,"feedback":"Good implementation..."}
   ```
4. Verify next question event is received:
   ```
   event: question
   data: {"session_id":"{session_id}","question_id":"{new_question_id}","content":"..."}
   ```
5. Submit answers until 5 questions are completed
6. Verify interview completion event is received:
   ```
   event: completion
   data: {"session_id":"{session_id}","status":"completed","total_score":82}
   ```

## 3. Edge Case Testing

### 3.1 Input Validation Testing

#### TC-015: Extremely Long Answers
**Test Steps:**
1. Generate extremely long answer (>10KB) with valid Go code
2. Submit answer via API
3. Verify system handles large input without crashing
4. Verify evaluation is still performed

#### TC-017: Empty Answers
**Test Steps:**
1. Submit empty string as answer
2. Verify appropriate error response
3. Check that validation error is logged

#### TC-019: Malformed JSON
**Test Steps:**
1. Send POST request with malformed JSON body
2. Verify HTTP 400 Bad Request response
3. Verify error message indicates JSON parsing issue

### 3.2 Interview Boundary Testing

#### TC-020: Maximum Question Limit
**Test Steps:**
1. Start new interview
2. Submit answers for 5 questions
3. Verify no 6th question is generated
4. Verify interview status is set to "completed"

#### TC-022: Concurrent Submissions
**Test Steps:**
1. Start new interview
2. Send 2 concurrent answer submissions for the same question
3. Verify only one answer is accepted
4. Check for appropriate conflict handling

## 4. Error Handling Testing

### 4.1 API Error Testing

#### TC-025: Invalid Session ID
**Test Steps:**
1. Send POST request with non-existent session ID
2. Verify HTTP 404 Not Found response
3. Verify error message indicates invalid session

#### TC-027: Unauthorized Access
**Test Steps:**
1. Send API request without authorization token
2. Verify HTTP 401 Unauthorized response
3. Send API request with invalid token
4. Verify HTTP 401 Unauthorized response

### 4.2 Service Failure Testing

#### TC-029: AI Service Unavailability
**Prerequisites:**
- Test environment with controllable AI service

**Test Steps:**
1. Stop mock AI service
2. Attempt to generate question or evaluate answer
3. Verify graceful error handling
4. Check that appropriate fallback is attempted
5. Verify error is logged properly

#### TC-032: Database Connection Failure
**Test Steps:**
1. Stop test database
2. Attempt API operations that require database
3. Verify error responses are appropriate
4. Check that connection retry logic works when database is restored

## 5. Performance Testing

### 5.1 Response Time Testing

#### Critical Path Performance
**Test Steps:**
1. Set up performance testing environment
2. Measure response time for resume upload
   - Should complete within 3 seconds
3. Measure response time for question generation
   - Should complete within 5 seconds
4. Measure response time for answer evaluation
   - Should complete within 10 seconds
5. Record average and 95th percentile response times

### 5.2 Load Testing

#### Concurrent Interviews
**Test Steps:**
1. Use load testing tool (e.g., k6, JMeter)
2. Configure to simulate 100+ concurrent users
3. Each user should:
   - Start a new interview
   - Answer 5 questions
   - Complete the interview
4. Run test for 10 minutes
5. Monitor system performance metrics:
   - CPU usage
   - Memory usage
   - Response times under load
   - Error rates

## 6. Test Cleanup

After each test, ensure proper cleanup:
1. Delete test interview sessions
2. Clear test answer and evaluation records
3. Reset mock service states
4. Close database connections
5. Clear Redis cache

## 7. Test Data Reset

For consistent testing, reset test data before each test run:
1. Run database migrations to reset schema
2. Seed fresh test data
3. Clear any persistent caches
4. Restart required services