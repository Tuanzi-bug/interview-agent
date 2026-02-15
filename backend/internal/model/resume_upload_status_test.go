package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB 设置测试数据库
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// 设置全局 getDB 函数
	SetDBGetter(func() *gorm.DB {
		return db
	})

	// 自动迁移
	err = db.AutoMigrate(&ResumeUploadStatus{})
	assert.NoError(t, err)

	return db
}

// TestResumeUploadStatusTableName 测试表名
func TestResumeUploadStatusTableName(t *testing.T) {
	status := ResumeUploadStatus{}
	assert.Equal(t, "resume_upload_status", status.TableName())
}

// TestCreateUploadStatus 测试创建上传状态
func TestCreateUploadStatus(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	status := &ResumeUploadStatus{
		UploadID: "test-uuid-123",
		UserID:   1,
		FilePath: "/uploads/test.pdf",
		Status:   "pending",
		Progress: 0,
	}

	err := ResumeUploadStatusDao.CreateUploadStatus(status)
	assert.NoError(t, err)
	assert.NotZero(t, status.ID)
}

// TestGetByUploadID 测试根据UploadID获取状态
func TestGetByUploadID(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// 创建测试数据
	status := &ResumeUploadStatus{
		UploadID: "test-uuid-456",
		UserID:   2,
		FilePath: "/uploads/test2.pdf",
		Status:   "extracting",
		Progress: 50,
	}
	err := ResumeUploadStatusDao.CreateUploadStatus(status)
	assert.NoError(t, err)

	// 测试查询
	found, err := ResumeUploadStatusDao.GetByUploadID("test-uuid-456")
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, "test-uuid-456", found.UploadID)
	assert.Equal(t, uint(2), found.UserID)
	assert.Equal(t, "extracting", found.Status)
	assert.Equal(t, 50, found.Progress)
}

// TestGetByUploadIDNotFound 测试查询不存在的记录
func TestGetByUploadIDNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	found, err := ResumeUploadStatusDao.GetByUploadID("non-existent-uuid")
	assert.Error(t, err)
	assert.Nil(t, found)
}

// TestUpdateStatus 测试更新状态
func TestUpdateStatus(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// 创建测试数据
	status := &ResumeUploadStatus{
		UploadID: "test-uuid-789",
		UserID:   3,
		FilePath: "/uploads/test3.pdf",
		Status:   "pending",
		Progress: 0,
	}
	err := ResumeUploadStatusDao.CreateUploadStatus(status)
	assert.NoError(t, err)

	// 更新状态
	updates := map[string]interface{}{
		"status":   "completed",
		"progress": 100,
		"stage":    "Finished",
	}
	err = ResumeUploadStatusDao.UpdateStatus("test-uuid-789", updates)
	assert.NoError(t, err)

	// 验证更新结果
	found, err := ResumeUploadStatusDao.GetByUploadID("test-uuid-789")
	assert.NoError(t, err)
	assert.Equal(t, "completed", found.Status)
	assert.Equal(t, 100, found.Progress)
	assert.Equal(t, "Finished", found.Stage)
}

// TestUpdateStatusWithError 测试更新失败场景
func TestUpdateStatusWithError(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// 创建测试数据
	status := &ResumeUploadStatus{
		UploadID: "test-uuid-error",
		UserID:   4,
		FilePath: "/uploads/test4.pdf",
		Status:   "extracting",
		Progress: 30,
	}
	err := ResumeUploadStatusDao.CreateUploadStatus(status)
	assert.NoError(t, err)

	// 更新为失败状态
	updates := map[string]interface{}{
		"status":    "failed",
		"progress":  30,
		"error_msg": "Failed to extract PDF content",
	}
	err = ResumeUploadStatusDao.UpdateStatus("test-uuid-error", updates)
	assert.NoError(t, err)

	// 验证错误信息
	found, err := ResumeUploadStatusDao.GetByUploadID("test-uuid-error")
	assert.NoError(t, err)
	assert.Equal(t, "failed", found.Status)
	assert.Equal(t, "Failed to extract PDF content", found.ErrorMsg)
}

// TestStatusFieldsValidation 测试状态字段约束
func TestStatusFieldsValidation(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	status := &ResumeUploadStatus{
		UploadID:        "test-uuid-complete",
		UserID:          5,
		FilePath:        "/uploads/complete.pdf",
		ResumeID:        100,
		Status:          "completed",
		Progress:        100,
		Stage:           "All done",
		ExtractDuration: 1500,
		AnalyzeDuration: 3000,
	}

	err := ResumeUploadStatusDao.CreateUploadStatus(status)
	assert.NoError(t, err)

	// 验证所有字段
	found, err := ResumeUploadStatusDao.GetByUploadID("test-uuid-complete")
	assert.NoError(t, err)
	assert.Equal(t, uint64(100), found.ResumeID)
	assert.Equal(t, int64(1500), found.ExtractDuration)
	assert.Equal(t, int64(3000), found.AnalyzeDuration)
	assert.NotZero(t, found.CreatedAt)
	assert.NotZero(t, found.UpdatedAt)
}

// TestTimestamps 测试时间戳自动设置
func TestTimestamps(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	status := &ResumeUploadStatus{
		UploadID: "test-uuid-timestamp",
		UserID:   6,
		FilePath: "/uploads/timestamp.pdf",
		Status:   "pending",
	}

	// 创建时记录时间
	beforeCreate := time.Now()
	err := ResumeUploadStatusDao.CreateUploadStatus(status)
	assert.NoError(t, err)
	afterCreate := time.Now()

	// 验证创建时间
	assert.True(t, status.CreatedAt.After(beforeCreate) || status.CreatedAt.Equal(beforeCreate))
	assert.True(t, status.CreatedAt.Before(afterCreate) || status.CreatedAt.Equal(afterCreate))

	// 等待一点时间再更新
	time.Sleep(10 * time.Millisecond)

	// 更新记录
	updates := map[string]interface{}{
		"status": "extracting",
	}
	err = ResumeUploadStatusDao.UpdateStatus("test-uuid-timestamp", updates)
	assert.NoError(t, err)

	// 验证更新时间变化
	found, err := ResumeUploadStatusDao.GetByUploadID("test-uuid-timestamp")
	assert.NoError(t, err)
	assert.True(t, found.UpdatedAt.After(found.CreatedAt))
}
