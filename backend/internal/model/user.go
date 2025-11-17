package model

import (
	"time"

	"gorm.io/gorm"
)

var UserDao _User

// User 用户模型
type (
	_User struct {
	}
	User struct {
		ID            uint           `json:"id" gorm:"primaryKey"`
		Username      string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
		Email         string         `json:"email" gorm:"uniqueIndex;size:100;not null"`
		PasswordHash  string         `json:"-" gorm:"size:255;not null"`
		Role          string         `json:"role" gorm:"size:20;default:'user'"`
		WechatOpenID  *string        `json:"wechat_open_id" gorm:"uniqueIndex;size:100"`  // 微信OpenID
		WechatUnionID *string        `json:"wechat_union_id" gorm:"uniqueIndex;size:100"` // 微信UnionID
		Nickname      string         `json:"nickname" gorm:"size:100"`                    // 微信昵称
		Avatar        string         `json:"avatar" gorm:"size:255"`                      // 微信头像
		CreatedAt     time.Time      `json:"created_at"`
		UpdatedAt     time.Time      `json:"updated_at"`
		DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
		Resumes       []Resume       `json:"resumes,omitempty" gorm:"foreignKey:UserID"`
		Interviews    []Interview    `json:"interviews,omitempty" gorm:"foreignKey:UserID"`
	}
)

// Resume 简历模型
type Resume struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"index;not null"`
	Name      string         `json:"name" gorm:"size:100;not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	Analysis  string         `json:"analysis,omitempty" gorm:"type:text"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	User      User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// Interview 面试模型
type Interview struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"index;not null"`
	ResumeID   uint           `json:"resume_id" gorm:"index"`
	JobTitle   string         `json:"job_title" gorm:"size:100;not null"`
	Difficulty string         `json:"difficulty" gorm:"size:20;not null"`
	Type       string         `json:"type" gorm:"size:50;not null"` // comprehensive, specialized, resume-based
	Status     string         `json:"status" gorm:"size:20;default:'pending'"`
	StartTime  *time.Time     `json:"start_time,omitempty"`
	EndTime    *time.Time     `json:"end_time,omitempty"`
	Duration   int            `json:"duration,omitempty"` // 单位：秒
	Score      float64        `json:"score,omitempty"`
	Evaluation string         `json:"evaluation,omitempty" gorm:"type:text"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	User       User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Resume     Resume         `json:"resume,omitempty" gorm:"foreignKey:ResumeID"`
	Questions  []Question     `json:"questions,omitempty" gorm:"foreignKey:InterviewID"`
}

// Question 面试问题模型
type Question struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	InterviewID uint           `json:"interview_id" gorm:"index;not null"`
	Content     string         `json:"content" gorm:"type:text;not null"`
	Type        string         `json:"type" gorm:"size:50;not null"` // technical, behavioral, coding
	Answer      string         `json:"answer,omitempty" gorm:"type:text"`
	Score       float64        `json:"score,omitempty"`
	Feedback    string         `json:"feedback,omitempty" gorm:"type:text"`
	Order       int            `json:"order" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Interview   Interview      `json:"interview,omitempty" gorm:"foreignKey:InterviewID"`
}

// QuestionBank 题库模型
type QuestionBank struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	Content         string         `json:"content" gorm:"type:text;not null"`
	Type            string         `json:"type" gorm:"size:50;not null"`
	Category        string         `json:"category" gorm:"size:50;not null"`
	Difficulty      string         `json:"difficulty" gorm:"size:20;not null"`
	Tags            string         `json:"tags,omitempty" gorm:"size:255"`
	ReferenceAnswer string         `json:"reference_answer,omitempty" gorm:"type:text"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// EvaluationCriteria 评估标准模型
type EvaluationCriteria struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Description string         `json:"description" gorm:"type:text"`
	Weight      float64        `json:"weight" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type QuestionBankEntrance struct {
	ID           uint           `json:"id" gorm:"primaryKey;comment:主键"`                           // 对应 bigint 主键
	CategoryName string         `json:"category_name" gorm:"size:32;comment:分类名称"`                 // varchar(32)
	CategoryType *int8          `json:"category_type,omitempty" gorm:"type:tinyint;comment:分类的类型"` // tinyint(4)，可选
	ImageURL     string         `json:"image_url,omitempty" gorm:"size:64;comment:图标链接"`           // varchar(64)，可选
	ParentID     *uint          `json:"parent_id,omitempty" gorm:"comment:父级id"`                   // 父分类ID，可选（层级关联）
	CreatedBy    string         `json:"created_by,omitempty" gorm:"size:32;comment:创建人"`           // varchar(32)，可选
	CreatedTime  time.Time      `json:"created_time" gorm:"comment:创建时间"`                          // datetime，自动填充
	UpdatedBy    string         `json:"updated_by,omitempty" gorm:"size:32;comment:更新人"`           // varchar(32)，可选
	UpdatedTime  time.Time      `json:"updated_time" gorm:"comment:更新时间"`                          // datetime，自动填充
	IsDeleted    *int8          `json:"is_deleted,omitempty" gorm:"type:tinyint;comment:是否删除"`     // tinyint(2)，0=未删/1=已删，可选
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`                                            // 软删除字段，和 Resume 保持一致
}

func (User) TableName() string {
	return "user"
}

// Create 创建用户记录
func (u *_User) Create(user *User) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	return getDB().Create(user).Error
}

// FindByUsernameOrEmail 根据用户名或邮箱查询用户
func (u *_User) FindByUsernameOrEmail(username, email string) (*User, error) {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	var user User
	err := getDB().
		Where("username = ? OR email = ?", username, email).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail 根据邮箱查询用户
func (u *_User) FindByEmail(email string) (*User, error) {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	var user User
	err := getDB().
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 根据ID查询用户
func (u *_User) FindByID(id uint) (*User, error) {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	var user User
	err := getDB().
		Where("id = ?", id).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateByID 根据ID更新用户字段
func (u *_User) UpdateByID(id uint, updates map[string]interface{}) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	if len(updates) == 0 {
		return nil
	}
	return getDB().
		Model(&User{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// FindByWechatOpenID 根据微信 OpenID 查询用户
func (u *_User) FindByWechatOpenID(openID string) (*User, error) {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	var user User
	err := getDB().
		Where("wechat_open_id = ?", openID).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
