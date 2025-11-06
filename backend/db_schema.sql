-- 数据库表结构设计
-- 生成时间: 2024-01-15

-- 创建users表
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 创建resumes表（注意：interviews表有外键引用resumes表，所以resumes需要先创建）
CREATE TABLE resumes (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) REFERENCES users(id),
    file_url VARCHAR(255) NOT NULL,
    filename VARCHAR(255) NOT NULL,
    content TEXT,
    analyzed_data JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建interviews表
CREATE TABLE interviews (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) REFERENCES users(id),
    resume_id VARCHAR(36) REFERENCES resumes(id),
    interview_type VARCHAR(50) NOT NULL,
    category VARCHAR(50) NOT NULL,
    difficulty VARCHAR(50) NOT NULL,
    duration INT, -- 分钟
    overall_score FLOAT,
    status VARCHAR(20) DEFAULT 'created', -- created, in_progress, completed, cancelled
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

-- 创建questions表
CREATE TABLE questions (
    id VARCHAR(36) PRIMARY KEY,
    interview_id VARCHAR(36) REFERENCES interviews(id),
    content TEXT NOT NULL,
    question_type VARCHAR(50) NOT NULL,
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建answers表
CREATE TABLE answers (
    id VARCHAR(36) PRIMARY KEY,
    question_id VARCHAR(36) REFERENCES questions(id),
    content TEXT NOT NULL,
    score FLOAT,
    evaluation_data JSON,
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建reports表
CREATE TABLE reports (
    id VARCHAR(36) PRIMARY KEY,
    interview_id VARCHAR(36) REFERENCES interviews(id),
    overall_score FLOAT,
    skills_json JSON, -- 各技能评分
    strengths TEXT, -- 优点，JSON数组
    weaknesses TEXT, -- 不足，JSON数组
    suggestions TEXT, -- 建议，JSON数组
    detailed_feedback TEXT,
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引以优化查询性能
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_resumes_user_id ON resumes(user_id);
CREATE INDEX idx_interviews_user_id ON interviews(user_id);
CREATE INDEX idx_interviews_status ON interviews(status);
CREATE INDEX idx_questions_interview_id ON questions(interview_id);
CREATE INDEX idx_answers_question_id ON answers(question_id);
CREATE INDEX idx_reports_interview_id ON reports(interview_id);