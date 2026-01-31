# 业务流程详解

> 以流程图的方式理解系统的核心业务逻辑

## 📋 学习任务

1. 理解每个核心业务流程的完整路径
2. 掌握服务层的业务逻辑编排方式
3. 理解数据在各层之间的流转
4. 理解外部服务（AI Agent、微信API）的集成方式

---

## 目录

- [1. 用户注册登录流程](#1-用户注册登录流程)
- [2. 微信登录流程](#2-微信登录流程)
- [3. 用户模型管理流程](#3-用户模型管理流程)
- [4. 简历上传解析流程](#4-简历上传解析流程)
- [5. 面试创建流程](#5-面试创建流程)
- [6. 面试进行流程](#6-面试进行流程)
- [7. 押题生成流程](#7-押题生成流程)
- [8. 评估报告生成流程](#8-评估报告生成流程)

---

## 1. 用户注册登录流程

### 1.1 用户注册流程

```mermaid
sequenceDiagram
    participant C as 客户端
    participant H as Handler
    participant S as UserManager
    participant DAO as UserDao
    participant MW as Middleware
    participant DB as MySQL

    C->>H: POST /api/user/register<br/>{username, email, password}
    H->>H: 绑定和验证请求参数
    H->>S: Register(ctx, req)
    
    Note over S: 业务逻辑处理
    S->>DAO: FindByUsernameOrEmail(username, email)
    DAO->>DB: SELECT * FROM user<br/>WHERE username=? OR email=?
    
    alt 用户已存在
        DB-->>DAO: 返回用户记录
        DAO-->>S: User
        S-->>H: Error: 用户名或邮箱已存在
        H-->>C: 400 Bad Request
    else 用户不存在
        DB-->>DAO: ErrRecordNotFound
        DAO-->>S: ErrRecordNotFound
        
        S->>S: 创建User对象<br/>role='user'
        Note over S: 注意：密码明文存储<br/>生产环境应使用bcrypt
        
        S->>DAO: Create(user)
        DAO->>DB: INSERT INTO user
        DB-->>DAO: Success, ID=123
        DAO-->>S: Success
        
        S->>MW: GenerateToken(userID, username, role)
        MW->>MW: 生成JWT token
        MW-->>S: token
        
        S->>S: buildLoginResponse(token, user)
        S-->>H: LoginResponse{token, userProfile}
        H-->>C: 200 OK<br/>{token, user}
    end
```

**关键点**:
1. **唯一性检查**: 注册前检查用户名和邮箱是否已存在
2. **默认角色**: 新用户默认角色为 'user'
3. **JWT生成**: 注册成功后立即生成token，实现自动登录
4. **密码安全**: 当前密码明文存储在PasswordHash字段，生产环境应使用bcrypt加密

### 1.2 用户登录流程

```mermaid
sequenceDiagram
    participant C as 客户端
    participant H as Handler
    participant S as UserManager
    participant DAO as UserDao
    participant MW as Middleware
    participant DB as MySQL

    C->>H: POST /api/user/login<br/>{email, password}
    H->>H: 绑定和验证请求参数
    H->>S: Login(ctx, req)
    
    S->>DAO: FindByEmail(email)
    DAO->>DB: SELECT * FROM user<br/>WHERE email=?
    
    alt 用户不存在
        DB-->>DAO: ErrRecordNotFound
        DAO-->>S: ErrRecordNotFound
        S-->>H: Error: 用户不存在
        H-->>C: 400 Bad Request
    else 用户存在
        DB-->>DAO: User
        DAO-->>S: User
        
        S->>S: 验证密码<br/>user.PasswordHash == req.Password
        
        alt 密码错误
            S-->>H: Error: 密码错误
            H-->>C: 401 Unauthorized
        else 密码正确
            S->>MW: GenerateToken(userID, username, role)
            MW-->>S: token
            
            S->>S: buildLoginResponse(token, user)
            S-->>H: LoginResponse{token, userProfile}
            H-->>C: 200 OK<br/>{token, user}
        end
    end
```

**关键点**:
1. **邮箱登录**: 使用邮箱作为登录凭证
2. **密码验证**: 直接比较明文密码（生产环境应使用bcrypt.CompareHashAndPassword）
3. **错误区分**: 区分"用户不存在"和"密码错误"两种错误情况
4. **Token生成**: 登录成功生成JWT token用于后续认证

---

## 2. 微信登录流程

### 2.1 微信OAuth 2.0 授权流程

```mermaid
sequenceDiagram
    participant U as 用户
    participant C as 客户端
    participant H as Handler
    participant S as UserManager
    participant WX as 微信API
    participant DAO as UserDao
    participant DB as MySQL

    U->>C: 点击"微信登录"
    C->>H: GET /api/user/wechat/login
    H->>S: WechatLogin(ctx)
    
    S->>S: 构建微信OAuth URL<br/>appid + redirect_uri + scope
    S-->>H: WechatLoginQRResponse{loginURL}
    H-->>C: 200 OK {loginURL}
    
    C->>U: 显示二维码或跳转
    U->>WX: 扫码授权
    
    WX->>C: 重定向到 redirect_uri?code=CODE&state=STATE
    C->>H: GET /api/user/wechat/callback?code=CODE
    H->>H: 验证参数
    H->>S: WechatCallback(ctx, req)
    
    Note over S: 步骤1: 获取access_token
    S->>WX: GET /sns/oauth2/access_token<br/>?appid=&secret=&code=&grant_type=
    WX-->>S: {access_token, openid, unionid}
    
    Note over S: 步骤2: 获取用户信息
    S->>WX: GET /sns/userinfo<br/>?access_token=&openid=
    WX-->>S: {openid, nickname, headimgurl, unionid}
    
    Note over S: 步骤3: 查询或创建用户
    S->>DAO: FindByWechatOpenID(openid)
    DAO->>DB: SELECT * FROM user<br/>WHERE wechat_open_id=?
    
    alt 用户已存在
        DB-->>DAO: User
        DAO-->>S: User
        
        S->>S: 更新用户信息<br/>(nickname, avatar, unionid)
        S->>DAO: UpdateByID(userID, updates)
        DAO->>DB: UPDATE user SET ...
        DB-->>DAO: Success
    else 用户不存在
        DB-->>DAO: ErrRecordNotFound
        DAO-->>S: ErrRecordNotFound
        
        S->>S: 生成用户名<br/>wechat_<openid前8位>
        S->>S: 创建新用户<br/>username, openid, unionid, nickname, avatar
        S->>DAO: Create(newUser)
        DAO->>DB: INSERT INTO user
        DB-->>DAO: Success, ID=456
        DAO-->>S: newUser
    end
    
    S->>S: GenerateToken(userID, username, role)
    S->>S: buildLoginResponse(token, user)
    S-->>H: LoginResponse{token, userProfile}
    H-->>C: 200 OK {token, user}
    C->>U: 登录成功，跳转到首页
```

**关键点**:
1. **OAuth 2.0流程**: 标准的授权码模式（Authorization Code）
2. **两步获取**: 先用code换access_token，再用access_token获取用户信息
3. **自动注册**: 首次微信登录自动创建账号
4. **UnionID**: 支持微信多平台（公众号、小程序、网站）用户统一
5. **信息更新**: 每次登录更新微信昵称和头像
6. **用户名生成**: 格式为 `wechat_<openid前8位>`

---

## 3. 用户模型管理流程

### 3.1 创建用户模型配置

```mermaid
sequenceDiagram
    participant C as 客户端
    participant H as Handler
    participant S as ModelManager
    participant COM as CommonUtil
    participant DAO as UserModelDao
    participant DB as MySQL

    C->>H: POST /api/user/model/create<br/>{name, model_key, protocol, base_url, api_key, ...}
    H->>H: 获取userID（JWT中间件）
    H->>H: 绑定和验证请求
    H->>S: CreateUserModel(ctx, userID, req)
    
    Note over S: 步骤1: 加密API密钥
    S->>COM: EncryptAPIKey(req.APIKey)
    COM->>COM: 获取密钥<br/>环境变量或默认值
    COM->>COM: AES-CBC加密
    COM-->>S: encryptedAPIKey
    
    Note over S: 步骤2: 处理可选字段
    S->>S: 设置默认值<br/>scope=7, status=1
    S->>S: 处理metaID, configJSON, defaultParams
    
    Note over S: 步骤3: 处理默认模型逻辑
    alt isDefault == 1
        S->>DAO: CancelDefaultUserModel(userID, 0)
        DAO->>DB: UPDATE user_model<br/>SET is_default=0<br/>WHERE user_id=?
        DB-->>DAO: Success
    end
    
    Note over S: 步骤4: 创建模型记录
    S->>S: 构建UserModel对象
    S->>DAO: CreateUserModel(newModel)
    DAO->>DB: INSERT INTO user_model
    DB-->>DAO: Success, ID=789
    DAO-->>S: Success
    
    Note over S: 步骤5: 设置启用状态
    alt status == 1
        S->>DAO: SetEnabledUserModel(userID, modelID)
        DAO->>DB: UPDATE user_model<br/>SET status=0 WHERE user_id=?<br/>UPDATE user_model<br/>SET status=1 WHERE id=?
        DB-->>DAO: Success
    end
    
    S-->>H: "success"
    H-->>C: 200 OK {"message": "success"}
```

**关键点**:
1. **API密钥加密**: 使用AES-CBC加密存储，密钥从环境变量获取
2. **默认模型管理**: 设为默认时自动取消其他模型的默认状态
3. **启用状态管理**: status=1时将其他模型设为禁用
4. **事务一致性**: 默认状态和启用状态的变更在DAO层使用事务保证

### 3.2 更新用户模型配置

```mermaid
sequenceDiagram
    participant C as 客户端
    participant H as Handler
    participant S as ModelManager
    participant COM as CommonUtil
    participant DAO as UserModelDao
    participant DB as MySQL

    C->>H: PUT /api/user/model/update<br/>{id, name, api_key?, ...}
    H->>S: UpdateUserModel(ctx, userID, req)
    
    Note over S: 步骤1: 获取现有模型
    S->>DAO: GetUserModelByID(userID, modelID)
    DAO->>DB: SELECT * FROM user_model<br/>WHERE id=? AND user_id=? AND deleted=0
    
    alt 模型不存在
        DB-->>DAO: ErrRecordNotFound
        DAO-->>S: NotFoundError
        S-->>H: Error: 模型不存在
        H-->>C: 404 Not Found
    else 模型存在
        DB-->>DAO: UserModel
        DAO-->>S: existingModel
        
        Note over S: 步骤2: 更新基本字段
        S->>S: 更新 name, model_key, protocol, base_url, provider_name
        
        Note over S: 步骤3: 处理可选字段
        alt req.IsSetAPIKey()
            S->>COM: EncryptAPIKey(req.APIKey)
            COM-->>S: encryptedAPIKey
            S->>S: existingModel.APIKeyEncrypted = encryptedAPIKey
        end
        
        S->>S: 处理 MetaID, DefaultParams, ConfigJSON, Scope, Status
        
        Note over S: 步骤4: 更新UpdatedAt
        S->>S: existingModel.UpdatedAt = time.Now().UnixMilli()
        
        Note over S: 步骤5: 保存到数据库
        S->>DAO: UpdateUserModel(existingModel)
        DAO->>DB: UPDATE user_model SET ... WHERE id=?
        DB-->>DAO: Success
        
        Note over S: 步骤6: 处理默认状态变更
        alt isDefault == 1
            S->>DAO: SetDefaultUserModel(userID, modelID)
            DAO->>DB: 事务：取消其他默认 + 设置当前为默认
            DB-->>DAO: Success
        else isDefault == 0
            S->>DAO: CancelDefaultUserModel(userID, modelID)
            DAO->>DB: UPDATE is_default=0
            DB-->>DAO: Success
        end
        
        Note over S: 步骤7: 处理启用状态
        alt status == 1
            S->>DAO: SetEnabledUserModel(userID, modelID)
            DAO->>DB: 事务：禁用其他 + 启用当前
            DB-->>DAO: Success
        end
        
        S-->>H: Success
        H-->>C: 200 OK
    end
```

**关键点**:
1. **部分更新**: 只更新提供的字段，未提供的字段保持不变
2. **API密钥可选更新**: 只在提供新密钥时才重新加密
3. **状态管理**: 独立处理默认状态和启用状态的变更
4. **权限验证**: 通过userID确保只能更新自己的模型

---

## 4. 简历上传解析流程

### 4.1 完整的简历上传解析流程

```mermaid
sequenceDiagram
    participant C as 客户端
    participant H as Handler
    participant RS as ResumeService
    participant RA as ResumeAgent
    participant LLM as OpenAI API
    participant TOOL as PDFToTextTool
    participant DAO as ResumeDao
    participant DB as MySQL

    C->>H: POST /api/resume/upload<br/>multipart/form-data {file}
    H->>H: 获取userID（JWT）
    H->>H: 解析multipart表单
    
    Note over H: 文件验证
    H->>H: 检查文件类型（仅PDF）
    H->>H: 检查文件大小（≤10MB）
    
    alt 验证失败
        H-->>C: 400 Bad Request<br/>Invalid file type or size
    else 验证通过
        Note over H: 保存文件到本地
        H->>H: 生成文件路径<br/>uploads/resumes/{timestamp}_{filename}
        H->>H: 保存文件到磁盘
        
        Note over H: 调用解析服务
        H->>RS: ParseResumeAndSave(ctx, userID, filePath, fileSize)
        
        Note over RS: 步骤1: 创建Resume Agent
        RS->>RA: NewResumeParserAgent(userID)
        RA->>RA: 创建OpenAI ChatModel
        RA->>RA: 配置Prompt模板<br/>输出JSON格式
        RA->>RA: 注册PDFToTextTool工具
        RA->>RA: 设置MaxIterations=20
        RA-->>RS: agent
        
        Note over RS: 步骤2: 运行Agent
        RS->>RS: 构建消息<br/>SystemMessage + UserMessage(filePath)
        RS->>RA: Run(ctx, messages)
        
        Note over RA: Agent执行循环
        RA->>TOOL: 调用PDFToTextTool<br/>提取PDF文本
        TOOL->>TOOL: 读取PDF文件
        TOOL->>TOOL: 转换为纯文本
        TOOL-->>RA: PDF文本内容
        
        RA->>LLM: 发送Prompt + PDF文本
        Note over LLM: AI分析简历<br/>提取结构化信息
        LLM-->>RA: JSON响应<br/>{name, contact, education, work_experience, tech_stack, ...}
        
        Note over RA: 检查输出格式
        alt 输出完整且格式正确
            RA-->>RS: JSON字符串
        else 输出不完整或需要追问
            RA->>LLM: 继续对话获取完整信息
            LLM-->>RA: 补充信息
            RA-->>RS: 完整JSON字符串
        end
        
        Note over RS: 步骤3: 解析JSON
        RS->>RS: 清理Markdown代码块标记<br/>```json ... ```
        RS->>RS: json.Unmarshal(content, &resumeInfo)
        
        alt JSON解析失败
            RS-->>H: Error: 解析失败
            H-->>C: 500 Internal Server Error
        else JSON解析成功
            Note over RS: 步骤4: 保存到数据库
            RS->>RS: 构建Resume对象<br/>Content = json.Marshal(resumeInfo)
            RS->>DAO: CreateResume(resume)
            DAO->>DB: INSERT INTO resume<br/>(user_id, content, file_name, file_size, file_type)
            DB-->>DAO: Success, resumeID=123
            DAO-->>RS: resumeID
            
            RS-->>H: resumeID, resumeInfo
            H-->>C: 200 OK<br/>{resume_id: 123, ...resumeInfo}
        end
    end
```

**关键点**:
1. **文件验证**: 只接受PDF格式，最大10MB
2. **本地存储**: 文件保存到 `uploads/resumes/` 目录
3. **AI解析**: 使用OpenAI API + PDFToTextTool工具链
4. **结构化输出**: AI返回JSON格式的简历信息
5. **智能追问**: Agent可以多轮对话获取完整信息（MaxIterations=20）
6. **数据存储**: 解析后的JSON存储在Content字段

### 4.2 简历信息结构

AI解析后返回的JSON结构：
```json
{
  "name": "张三",
  "contact": {
    "phone": "138****8888",
    "email": "zhang@example.com",
    "github": "https://github.com/zhang"
  },
  "education": [
    {
      "institution": "XX大学",
      "degree": "本科",
      "major": "计算机科学与技术",
      "graduation": "2020"
    }
  ],
  "work_experience": [
    {
      "company": "XX科技",
      "position": "后端开发工程师",
      "duration": "2020.07 - 2023.05",
      "description": "负责..."
    }
  ],
  "tech_stack": ["Java", "Spring Boot", "MySQL", "Redis"],
  "projects": [...],
  "skills": [...],
  "strengths": "熟悉分布式系统...",
  "potential_weaknesses": "缺少大规模系统经验...",
  "recommended_difficulty": "中等",
  "interview_focus_areas": ["Spring框架", "数据库优化"],
  "suggested_questions_directions": ["项目架构", "性能优化"]
}
```

---

## 5. 面试创建流程

### 5.1 创建面试记录

```mermaid
sequenceDiagram
    participant C as 客户端
    participant H as Handler
    participant S as InterviewManager
    participant DAO as InterviewRecordDao
    participant DB as MySQL

    C->>H: POST /api/mianshi/stream/start<br/>{type, difficulty, domain, resume_id?, company_name?, position_name?}
    H->>H: 绑定和验证请求
    H->>H: 获取userID（JWT）
    H->>S: CreateInterviewRecord(ctx, recordDTO)
    
    Note over S: 处理可选字段
    S->>S: 提取 company_name, position_name<br/>（如果为指针类型且非空）
    S->>S: 初始化 status='pending'
    S->>S: 设置 duration=0
    
    Note over S: 构建InterviewRecord对象
    S->>S: record = InterviewRecord{<br/>  UserID: userID,<br/>  Type: type,<br/>  Difficulty: difficulty,<br/>  Domain: domain,<br/>  CompanyName: companyName,<br/>  PositionName: positionName,<br/>  Status: 'pending',<br/>  Duration: 0<br/>}
    
    S->>DAO: CreateInterviewRecord(record)
    DAO->>DB: INSERT INTO interview_record<br/>(user_id, type, difficulty, domain, ...)
    DB-->>DAO: Success, recordID=456
    DAO-->>S: recordID
    
    S->>S: 记录日志<br/>[CreateInterviewRecord] 成功: ID=456, UserID=?
    S-->>H: recordID
    
    H->>H: 创建Session（会话管理）
    H->>H: 设置SSE响应头
    H->>H: 启动异步面试循环
    H-->>C: SSE Stream<br/>{type: "session_id", session_id, record_id}
```

**关键点**:
1. **初始状态**: 创建时状态为 'pending'
2. **可选字段处理**: 公司名和岗位名为可选
3. **记录ID**: 创建成功后立即返回recordID用于后续关联
4. **日志记录**: 详细记录创建过程便于追踪

---

## 6. 面试进行流程

### 6.1 SSE流式面试流程

```mermaid
sequenceDiagram
    participant C as 客户端
    participant H as Handler
    participant SM as SessionManager
    participant E as InterviewEngine
    participant AG as InterviewAgent
    participant LLM as OpenAI API
    participant S as InterviewService
    participant DAO as InterviewDialogueDao
    participant DB as MySQL

    Note over C,DB: 前置：已创建面试记录和会话
    
    C->>H: SSE连接已建立
    H->>SM: CreateSession(userID, recordID, ...)
    SM-->>H: session
    
    H->>E: RunInterviewLoop(ctx, session)
    
    Note over E: 面试循环开始
    loop 直到面试结束
        Note over E: 步骤1: 生成问题
        E->>AG: 创建或获取Interview Agent
        E->>AG: 生成面试问题
        AG->>LLM: 发送Prompt<br/>（包含简历、领域、难度等上下文）
        LLM-->>AG: 返回问题
        AG-->>E: question
        
        E->>C: SSE事件 {type: "question", content: question}
        
        Note over E: 步骤2: 等待用户回答
        E->>E: 阻塞等待session.AnswerChan
        
        C->>H: POST /api/mianshi/answer/submit<br/>{session_id, answer}
        H->>SM: SubmitAnswer(session_id, answer)
        SM->>SM: 找到session
        SM->>SM: session.AnswerChan <- answer
        SM-->>H: Success
        H-->>C: 200 OK
        
        E->>E: answer := <-session.AnswerChan
        
        Note over E: 步骤3: 处理特殊命令
        alt answer == "quit"
            E->>C: SSE事件 {type: "quit"}
            E->>E: 结束循环
        else answer == "continue"
            E->>C: SSE事件 {type: "continue"}
            Note over E: 继续下一个问题
        else 正常答案
            Note over E: 步骤4: 评估答案
            E->>AG: 评估用户答案
            AG->>LLM: 发送答案进行评估
            LLM-->>AG: 评估结果
            AG-->>E: evaluation
            
            E->>C: SSE事件 {type: "evaluation", content: evaluation}
            
            Note over E: 步骤5: 保存对话
            E->>S: SaveInterviewDialogue(userID, recordID, question, answer)
            S->>DAO: Create(dialogue)
            DAO->>DB: INSERT INTO interview_dialogue
            DB-->>DAO: Success
            
            Note over E: 步骤6: 决定是否追问
            E->>AG: 判断是否需要追问
            AG->>AG: 分析答案质量
            
            alt 需要追问
                AG-->>E: followUpQuestion
                E->>C: SSE事件 {type: "follow_up", content: followUpQuestion}
                Note over E: 返回步骤2等待回答
            else 不需要追问
                E->>E: questionCount++
                
                alt questionCount >= maxQuestions
                    E->>E: 结束面试
                    E->>C: SSE事件 {type: "complete"}
                else 继续下一题
                    Note over E: 返回步骤1生成新问题
                end
            end
        end
    end
    
    Note over E: 面试结束，生成报告
    E->>E: 触发评估报告生成
```

**关键点**:
1. **SSE流式通信**: 实时推送问题、评估、追问等事件
2. **会话管理**: 通过SessionManager管理多个并发面试
3. **异步处理**: 问题生成在后台异步进行
4. **智能追问**: 根据答案质量决定是否追问
5. **答案通道**: 使用channel实现异步答案接收
6. **特殊命令**: 支持quit（退出）和continue（跳过）命令

### 6.2 面试状态流转

```mermaid
stateDiagram-v2
    [*] --> pending: 创建面试
    pending --> in_progress: 开始答题<br/>（首次提交答案）
    in_progress --> in_progress: 继续答题
    in_progress --> completed: 完成所有题目<br/>或用户选择结束
    in_progress --> cancelled: 用户选择退出<br/>(action=quit)
    completed --> [*]: 生成评估报告
    cancelled --> [*]
    
    note right of in_progress
        状态可在Handler层更新
        也可在Engine中自动更新
    end note
```

---

## 7. 押题生成流程

### 7.1 完整的押题生成流程

```mermaid
sequenceDiagram
    participant C as 客户端
    participant H as Handler
    participant S as PredictionService
    participant RDao as ResumeDao
    participant PA as PredictionAgent
    participant LLM as OpenAI API
    participant PDao as PredictionDao
    participant DB as MySQL

    C->>H: POST /api/prediction/start<br/>{resume_id, prediction_type, language, job_title, difficulty, company_name?}
    H->>H: 获取userID（JWT）
    H->>H: 绑定和验证请求
    H->>S: Predict(ctx, req, userID)
    
    Note over S: Panic恢复机制
    S->>S: defer func() { recover() }
    
    Note over S: 步骤1: 获取简历内容
    S->>RDao: GetResumeByID(resumeID)
    RDao->>DB: SELECT * FROM resume WHERE id=?
    
    alt 简历不存在
        DB-->>RDao: ErrRecordNotFound
        RDao-->>S: Error
        S-->>H: Error: resume not found
        H-->>C: 404 Not Found
    else 简历存在
        DB-->>RDao: resume
        RDao-->>S: resume
        
        Note over S: 步骤2: 构建Prompt
        S->>S: prompt = 简历内容 + 押题要求<br/>- 类型：校招/社招<br/>- 语言：java/go<br/>- 岗位：前端/后端<br/>- 难度：入门/进阶<br/>- 目标公司（可选）
        
        Note over S: 步骤3: 创建Prediction Agent
        S->>PA: NewPredictionAgent(userID)
        PA->>PA: 创建OpenAI ChatModel
        PA->>PA: 配置Prompt模板<br/>输出JSON格式<br/>{questions: [{question, content, focus, ...}]}
        PA-->>S: agent
        
        Note over S: 步骤4: 运行Agent
        S->>PA: Run(ctx, messages)
        PA->>LLM: 发送Prompt + 简历内容 + 要求
        
        Note over LLM: AI生成押题<br/>基于简历和要求<br/>生成20+道题目
        LLM-->>PA: JSON响应（可能带Markdown标记）
        
        Note over PA: 收集完整响应
        loop 遍历事件流
            PA->>PA: 收集message content
        end
        PA-->>S: JSON字符串
        
        Note over S: 步骤5: 清理和解析JSON
        S->>S: cleanJSONContent()<br/>去除 ```json 和 ``` 标记
        S->>S: json.Unmarshal(content, &result)
        
        alt JSON解析失败
            S-->>H: Error: 解析失败
            H-->>C: 500 Internal Server Error
        else JSON解析成功
            Note over S: 步骤6: 保存押题记录
            S->>S: 构建PredictionRecord{<br/>  UserID, ResumeID,<br/>  Type, Language,<br/>  JobTitle, Difficulty,<br/>  Company<br/>}
            
            S->>PDao: CreatePredictionRecord(record)
            PDao->>DB: INSERT INTO prediction_record
            DB-->>PDao: Success, recordID=789
            PDao-->>S: recordID
            
            alt recordID == 0
                S-->>H: Error: 数据库错误
                H-->>C: 500 Internal Server Error
            else recordID > 0
                Note over S: 步骤7: 保存题目列表
                S->>S: 构建Questions列表<br/>处理FollowUp字段（可能是string或数组）
                
                loop 遍历所有题目
                    S->>S: questions.append(PredictionQuestion{<br/>  RecordID: recordID,<br/>  Question, Content, Focus,<br/>  ThinkingPath, ReferenceAnswer,<br/>  FollowUp, Sort: i+1<br/>})
                end
                
                S->>PDao: CreatePredictionQuestions(questions)
                PDao->>DB: INSERT INTO prediction_question<br/>(批量插入)
                DB-->>PDao: Success
                PDao-->>S: Success
                
                Note over S: 步骤8: 构建响应
                S->>S: 构建PredictResponse{<br/>  RecordID,<br/>  Questions: [...题目列表]<br/>}
                S-->>H: PredictResponse
                H-->>C: 200 OK<br/>{record_id, questions: [...]}
            end
        end
    end
```

**关键点**:
1. **Panic恢复**: 使用defer+recover捕获AI调用可能的panic
2. **Prompt构建**: 包含完整简历内容和详细要求
3. **JSON清理**: AI可能返回Markdown代码块，需要清理
4. **FollowUp处理**: 追问字段可能是字符串或数组，需要灵活处理
5. **批量保存**: 题目列表使用批量插入提高性能
6. **Sort字段**: 保持题目顺序（1, 2, 3, ...）
7. **详细日志**: 每个步骤都有详细日志便于排查问题

### 7.2 押题结果结构

AI生成的JSON结构：
```json
{
  "questions": [
    {
      "question": "请介绍一下Spring Boot的自动配置原理",
      "content": "Spring Boot核心特性",
      "focus": "自动配置机制、@EnableAutoConfiguration注解",
      "thinking_path": "1. 解释自动配置的概念\n2. 说明@SpringBootApplication注解\n3. 讲解条件注解的作用",
      "reference_answer": "Spring Boot的自动配置...",
      "follow_up": "那你能说说@Conditional注解的使用场景吗？"
    },
    // ... 更多题目
  ]
}
```

---

## 8. 评估报告生成流程

### 8.1 面试评估报告生成

```mermaid
sequenceDiagram
    participant E as InterviewEngine
    participant EA as EvaluationAgent
    participant LLM as OpenAI API
    participant DDao as DialogueDao
    participant IEDao as InterviewEvaluationDao
    participant ARDao as AnswerReportDao
    participant DB as MySQL

    Note over E: 面试结束，触发评估
    E->>E: 更新面试状态为completed
    E->>E: 计算面试耗时duration
    
    Note over E: 步骤1: 获取对话历史
    E->>DDao: GetInterviewDialoguesByUserIdAndRecordId(userID, recordID)
    DDao->>DB: SELECT * FROM interview_dialogues<br/>WHERE user_id=? AND report_id=?<br/>ORDER BY id ASC
    DB-->>DDao: dialogues[]
    DDao-->>E: dialogues
    
    Note over E: 步骤2: 构建评估Prompt
    E->>E: 组织对话内容<br/>Q1: ...\nA1: ...\nQ2: ...\nA2: ...
    E->>E: 添加评估维度要求<br/>- 技术能力<br/>- 表达能力<br/>- 问题理解<br/>- 实战经验<br/>- 学习潜力<br/>- 综合素质
    
    Note over E: 步骤3: 调用评估Agent
    E->>EA: NewEvaluationAgent(userID)
    EA->>EA: 创建ChatModel
    EA->>EA: 配置评估Prompt模板
    EA-->>E: agent
    
    E->>EA: Run(ctx, evaluationPrompt)
    EA->>LLM: 发送完整对话历史 + 评估要求
    
    Note over LLM: AI分析面试表现<br/>生成多维度评估
    LLM-->>EA: JSON响应{<br/>  overall_comment: "总体评价",<br/>  overall_score: 85.5,<br/>  dimensions: [<br/>    {dimension_name, evaluation, score},<br/>    ...<br/>  ]<br/>}
    EA-->>E: evaluationJSON
    
    Note over E: 步骤4: 解析评估结果
    E->>E: json.Unmarshal(evaluationJSON, &evaluationResult)
    
    Note over E: 步骤5: 保存评估报告
    E->>E: 构建InterviewEvaluation{<br/>  UserID, ReportID,<br/>  Comment: overall_comment,<br/>  Score: overall_score,<br/>  Dimensions: [{<br/>    DimensionName,<br/>    Evaluation,<br/>    Score<br/>  }, ...]<br/>}
    
    E->>IEDao: CreateEvaluation(evaluation)
    IEDao->>DB: INSERT INTO interview_evaluation<br/>(user_id, report_id, comment, score, dimensions)
    DB-->>IEDao: Success
    
    Note over E: 步骤6: 生成答题报告
    E->>E: 构建AnswerReport<br/>遍历dialogues，为每个问题生成详细分析
    
    loop 每个对话
        E->>EA: 评估单个答案
        EA->>LLM: 发送单个Q&A + 评估要求
        LLM-->>EA: {<br/>  score: 85,<br/>  key_points: "关键点",<br/>  strengths: "优势",<br/>  weaknesses: "不足",<br/>  suggestion: "建议"<br/>}
        EA-->>E: answerEvaluation
        
        E->>E: Records.append({<br/>  Order: i,<br/>  Content: question,<br/>  Comment: answerEvaluation,<br/>  Message: [{Order, Question, Answer}]<br/>})
    end
    
    E->>E: 构建AnswerReport{<br/>  UserID, ReportID,<br/>  Records: [...]<br/>}
    
    E->>ARDao: CreateAnswerReport(answerReport)
    ARDao->>DB: INSERT INTO answer_report<br/>(user_id, report_id, records)
    DB-->>ARDao: Success
    
    E->>E: 发送SSE完成事件<br/>{type: "complete", message: "评估完成"}
```

**关键点**:
1. **对话历史**: 评估基于完整的对话记录
2. **多维度评估**: 6个维度全面评价面试表现
3. **JSON存储**: Dimensions和Records使用JSON字段存储复杂结构
4. **两层报告**: 
   - InterviewEvaluation: 总体评分和维度评分
   - AnswerReport: 每个问题的详细分析
5. **AI驱动**: 完全由AI Agent分析生成，无固定规则

### 8.2 评估报告结构

**InterviewEvaluation结构**:
```json
{
  "id": 1,
  "user_id": 123,
  "report_id": 456,
  "comment": "整体表现良好，技术基础扎实...",
  "score": 85.5,
  "dimensions": [
    {
      "dimension_name": "技术能力",
      "evaluation": "对Spring框架理解深入，能够说明核心原理...",
      "score": 88.0
    },
    {
      "dimension_name": "表达能力",
      "evaluation": "回答逻辑清晰，表达流畅...",
      "score": 85.0
    },
    {
      "dimension_name": "问题理解",
      "evaluation": "能够准确理解问题要点...",
      "score": 82.0
    },
    {
      "dimension_name": "实战经验",
      "evaluation": "有实际项目经验，能够结合实例说明...",
      "score": 86.0
    },
    {
      "dimension_name": "学习潜力",
      "evaluation": "展现出良好的学习能力和技术热情...",
      "score": 87.0
    },
    {
      "dimension_name": "综合素质",
      "evaluation": "综合素质较好，具备团队协作意识...",
      "score": 84.0
    }
  ],
  "created_at": "2026-01-29T10:30:00Z",
  "updated_at": "2026-01-29T10:30:00Z"
}
```

**AnswerReport结构**:
```json
{
  "id": 1,
  "user_id": 123,
  "report_id": 456,
  "records": [
    {
      "order": 1,
      "content": "请介绍一下Spring Boot的自动配置原理",
      "comment": {
        "score": 85,
        "key_points": "自动配置、条件注解、SPI机制",
        "difficulty": "中等",
        "strengths": "能够说明核心原理和关键注解",
        "weaknesses": "对源码实现细节了解不够深入",
        "suggestion": "建议深入阅读源码，理解自动配置的实现机制",
        "know_points": "Spring Boot自动配置、@Conditional注解、spring.factories",
        "thinking": "从@EnableAutoConfiguration入手，说明自动配置的触发流程",
        "reference": "自动配置通过spring.factories加载配置类，配合@Conditional条件注解实现..."
      },
      "message": [
        {
          "order": 1,
          "question": "请介绍一下Spring Boot的自动配置原理",
          "answer": "Spring Boot的自动配置主要通过@EnableAutoConfiguration注解实现..."
        },
        {
          "order": 2,
          "question": "那你能说说@Conditional注解的使用场景吗？",
          "answer": "@Conditional注解用于根据条件决定是否创建Bean..."
        }
      ]
    },
    // ... 更多题目的详细分析
  ],
  "created_at": "2026-01-29T10:30:00Z",
  "updated_at": "2026-01-29T10:30:00Z"
}
```

---

## 9. 完整业务流程总览

### 9.1 系统核心流程交互图

```mermaid
graph TB
    Start([用户访问系统])
    
    subgraph "用户认证"
        Register[注册/登录]
        WechatLogin[微信登录]
        Token[获取JWT Token]
    end
    
    subgraph "模型配置"
        CheckModel{检查模型配置}
        ConfigModel[配置AI模型]
        EncryptKey[加密API密钥]
    end
    
    subgraph "简历管理"
        UploadResume[上传简历PDF]
        ParseResume[AI解析简历]
        SaveResume[保存简历信息]
        SetDefault[设置默认简历]
    end
    
    subgraph "押题功能"
        SelectResume[选择简历]
        SetRequirements[设置要求]
        GenerateQuestions[AI生成题目]
        SaveQuestions[保存押题记录]
        ViewQuestions[查看题目列表]
    end
    
    subgraph "面试功能"
        CreateInterview[创建面试]
        StartInterview[开始面试SSE]
        AskQuestion[AI提问]
        UserAnswer[用户回答]
        Evaluate[AI评估]
        FollowUp{需要追问?}
        NextQuestion{继续面试?}
    end
    
    subgraph "报告生成"
        GenerateEvaluation[生成评估报告]
        GenerateAnswerReport[生成答题报告]
        ViewReport[查看报告]
    end
    
    Start --> Register
    Register --> Token
    WechatLogin --> Token
    
    Token --> CheckModel
    CheckModel -->|未配置| ConfigModel
    ConfigModel --> EncryptKey
    EncryptKey --> CheckModel
    CheckModel -->|已配置| UploadResume
    
    UploadResume --> ParseResume
    ParseResume --> SaveResume
    SaveResume --> SetDefault
    
    SetDefault --> SelectResume
    SelectResume --> SetRequirements
    SetRequirements --> GenerateQuestions
    GenerateQuestions --> SaveQuestions
    SaveQuestions --> ViewQuestions
    
    SetDefault --> CreateInterview
    CreateInterview --> StartInterview
    StartInterview --> AskQuestion
    AskQuestion --> UserAnswer
    UserAnswer --> Evaluate
    Evaluate --> FollowUp
    FollowUp -->|是| AskQuestion
    FollowUp -->|否| NextQuestion
    NextQuestion -->|是| AskQuestion
    NextQuestion -->|否| GenerateEvaluation
    
    GenerateEvaluation --> GenerateAnswerReport
    GenerateAnswerReport --> ViewReport
    
    style Register fill:#e1f5ff
    style ConfigModel fill:#fff4e6
    style UploadResume fill:#e8f5e9
    style GenerateQuestions fill:#f3e5f5
    style StartInterview fill:#fce4ec
    style GenerateEvaluation fill:#fff9c4
```

### 9.2 数据流转图

```mermaid
graph LR
    User[用户]
    Handler[Handler层]
    Service[Service层]
    DAO[DAO层]
    DB[(MySQL)]
    Agent[AI Agent]
    LLM[OpenAI API]
    
    User -->|HTTP请求| Handler
    Handler -->|调用服务| Service
    Service -->|数据操作| DAO
    DAO -->|SQL查询| DB
    
    Service -->|调用Agent| Agent
    Agent -->|API请求| LLM
    LLM -->|AI响应| Agent
    Agent -->|结果| Service
    
    DB -->|数据| DAO
    DAO -->|Model| Service
    Service -->|DTO| Handler
    Handler -->|JSON响应| User
    
    style User fill:#e1f5ff
    style Service fill:#fff4e6
    style Agent fill:#f3e5f5
    style LLM fill:#fce4ec
```

---

## 10. 关键业务规则

### 10.1 用户模型管理规则

| 规则 | 说明 | 实现位置 |
|------|------|----------|
| API密钥加密 | 所有API密钥必须加密存储 | CommonUtil.EncryptAPIKey |
| 默认模型唯一 | 一个用户只能有一个默认模型 | SetDefaultUserModel时取消其他 |
| 启用模型唯一 | 一个用户只能有一个启用模型 | SetEnabledUserModel时禁用其他 |
| 软删除 | 删除模型使用软删除标记 | Deleted=1 |

### 10.2 简历管理规则

| 规则 | 说明 | 实现位置 |
|------|------|----------|
| 文件类型限制 | 仅支持PDF格式 | Handler层验证 |
| 文件大小限制 | 最大10MB | Handler层验证 |
| 默认简历唯一 | 一个用户只能有一个默认简历 | SetDefaultResume时取消其他 |
| AI解析 | 所有简历必须通过AI解析 | ResumeAgent + PDFToTextTool |
| JSON存储 | 解析结果以JSON存储在Content字段 | ResumeService |

### 10.3 面试管理规则

| 规则 | 说明 | 实现位置 |
|------|------|----------|
| 会话隔离 | 每次面试独立的Session | SessionManager |
| SSE通信 | 面试过程使用SSE流式推送 | InterviewEngine |
| 状态流转 | pending → in_progress → completed | 状态机管理 |
| 对话记录 | 所有Q&A必须保存 | InterviewDialogue表 |
| 智能追问 | AI决定是否追问 | InterviewAgent |

### 10.4 押题管理规则

| 规则 | 说明 | 实现位置 |
|------|------|----------|
| 基于简历 | 押题必须基于用户简历 | 查询Resume表 |
| AI生成 | 所有题目由AI生成 | PredictionAgent |
| JSON输出 | AI输出必须是标准JSON格式 | Prompt模板约束 |
| 批量保存 | 题目列表批量插入 | CreatePredictionQuestions |
| 保持顺序 | Sort字段保持题目顺序 | Sort=1,2,3,... |

### 10.5 评估报告规则

| 规则 | 说明 | 实现位置 |
|------|------|----------|
| 完整对话 | 评估基于完整对话历史 | 查询所有InterviewDialogue |
| 多维度 | 必须包含6个评估维度 | Dimensions JSON字段 |
| 双层报告 | 总体评估 + 逐题分析 | InterviewEvaluation + AnswerReport |
| AI驱动 | 完全由AI分析生成 | EvaluationAgent |
| JSON存储 | 复杂结构用JSON存储 | GORM serializer:json |

---

## ❓ 问题记录

### Q1: 为什么面试使用SSE而不是WebSocket？
**我的理解**: SSE（Server-Sent Events）是单向通信，服务器推送消息到客户端，客户端通过HTTP POST提交答案。面试场景的特点：1) **服务器推送为主**：问题、评估、追问都是服务器主动推送；2) **答案提交不频繁**：用户回答通过HTTP POST即可；3) **更简单**：SSE基于HTTP，不需要协议升级，易于实现和调试；4) **更稳定**：断线重连机制简单，浏览器原生支持。WebSocket适合双向高频通信，面试场景没有这个需求。

### Q2: 为什么押题和评估都要调用AI Agent？直接存储模板不行吗？
**我的理解**: AI Agent的优势：1) **个性化**：基于用户简历生成定制题目，每个人都不同；2) **灵活性**：可以根据不同领域、难度、公司生成针对性题目；3) **质量**：AI评估更全面、更客观、更具体；4) **扩展性**：新增面试领域不需要人工编写题库。缺点是成本较高、响应时间较长。可以考虑混合模式：常见题目用模板（快速响应），深度分析用AI（高质量输出）。

### Q3: 为什么要分开保存InterviewEvaluation和AnswerReport？
**我的理解**: 两层报告服务不同目的：1) **InterviewEvaluation（总体评估）**：快速了解整体表现，6个维度打分，适合列表展示和筛选；2) **AnswerReport（答题报告）**：详细的逐题分析，包含每题的优劣势和改进建议，适合学习和提升。分开存储的好处：a) 查询效率：列表页只需加载Evaluation；b) 数据大小：AnswerReport的JSON较大，按需加载；c) 关注分离：总体评分和详细分析是不同维度的信息。

### Q4: JSON字段存储复杂数据有什么优缺点？
**我的理解**: 
**优点**：
1) **灵活性**：结构可以动态变化，无需修改表结构
2) **简化查询**：一次查询获取完整信息，无需JOIN
3) **原子性**：整体读写，无需事务处理多表
4) **便于扩展**：新增字段不影响旧数据

**缺点**：
1) **查询限制**：无法直接按JSON内部字段查询排序（MySQL 5.7+支持JSON函数但性能较差）
2) **索引限制**：JSON字段无法建立常规索引
3) **数据验证**：需要在应用层验证JSON结构
4) **存储空间**：可能占用更多空间

**使用建议**：适合不需要按内部字段查询、结构灵活多变的数据（如评估维度、答题记录）。如果需要频繁按某字段查询，应拆分为独立表。

---

## ✅ 学习检查点

- [ ] 理解用户注册登录的完整流程
- [ ] 理解微信OAuth 2.0授权流程
- [ ] 理解用户模型配置的创建和更新逻辑
- [ ] 理解简历上传和AI解析的完整过程
- [ ] 理解面试创建和SSE流式面试流程
- [ ] 理解押题生成的AI调用过程
- [ ] 理解评估报告的两层结构
- [ ] 能够画出任意核心业务的流程图
- [ ] 理解数据在各层之间的流转方式
- [ ] 理解AI Agent的集成和调用方式

---

**上一步**: [Service层概述 (README.md)](README.md)

**相关文档**:
- [数据模型详解 (models.md)](../02-data-layer/models.md)
- [API层文档 (handler-pattern.md)](../04-api-layer/handler-pattern.md)
- [AI核心架构 (agent-architecture.md)](../05-ai-core/agent-architecture.md)