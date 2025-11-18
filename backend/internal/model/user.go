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
	}
)

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
