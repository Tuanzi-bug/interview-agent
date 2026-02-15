package repository

import (
	"fmt"
	"log"
	"time"

	"ai-eino-interview-agent/internal/config"
	"ai-eino-interview-agent/internal/model"
	"ai-eino-interview-agent/internal/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库实例
var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase(dbConfig config.DatabaseConfig) error {
	// 配置GORM日志
	logLevel := logger.Info

	// 连接数据库，带重试机制
	var db *gorm.DB
	var err error
	maxRetries := 10
	retryDelay := 1 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(mysql.Open(dbConfig.DSN), &gorm.Config{
			Logger: logger.Default.LogMode(logLevel),
		})
		if err == nil {
			break
		}

		if i < maxRetries-1 {
			log.Printf("数据库连接尝试 %d/%d 失败: %v. 将在 %v 后重试...",
				i+1, maxRetries, err, retryDelay)
			time.Sleep(retryDelay)
			retryDelay *= 2
			if retryDelay > 30*time.Second {
				retryDelay = 30 * time.Second
			}
		} else {
			return fmt.Errorf("连接数据库失败，已尝试 %d 次: %w", maxRetries, err)
		}
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	if dbConfig.ConnMaxLifetime != "" {
		connMaxLifetime := utils.ParseDurationWithDefault(dbConfig.ConnMaxLifetime, 0, "conn_max_lifetime")
		if connMaxLifetime > 0 {
			sqlDB.SetConnMaxLifetime(connMaxLifetime)
		}
	}

	// 设置全局DB实例
	DB = db

	// 设置 model 包的 DB 获取函数
	model.SetDBGetter(GetDB)

	// 自动迁移数据库表结构
	err = migrateDatabase()
	if err != nil {
		return err
	}

	log.Println("数据库连接成功并完成迁移")
	return nil
}

// migrateDatabase 执行数据库迁移
func migrateDatabase() error {
	return DB.AutoMigrate(
		&model.User{},
		&model.UserModel{},
		&model.InterviewRecord{},
		&model.InterviewDialogue{},
		&model.InterviewEvaluation{},
		&model.AnswerReport{},
		&model.Resume{},
		&model.PredictionRecord{},
		&model.PredictionQuestion{},
	)
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
