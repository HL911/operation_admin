package database

import (
	"context"
	"fmt"
	"time"

	"operation_admin/backend/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// pingTimeout 定义数据库初始化后的连通性校验超时时间。
const pingTimeout = 5 * time.Second

// New 负责初始化 Gorm 与 PostgreSQL 连接池。
func New(cfg config.DatabaseConfig) (*gorm.DB, error) {
	// dialector 描述 Gorm 访问 PostgreSQL 所使用的驱动。
	dialector := postgres.Open(cfg.DSN())

	// db 负责承载后续所有数据库访问能力。
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("打开数据库连接失败: %w", err)
	}

	// sqlDB 暴露底层连接池能力，用于设置连接参数和执行 Ping。
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层连接池失败: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// pingContext 用于限制数据库连通性校验时间，避免启动阶段长期阻塞。
	pingContext, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err := sqlDB.PingContext(pingContext); err != nil {
		return nil, fmt.Errorf("数据库连通性校验失败: %w", err)
	}

	return db, nil
}

// Close 负责释放底层数据库连接池。
func Close(db *gorm.DB) error {
	// sqlDB 代表 Gorm 持有的标准库连接池实例。
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取底层连接池失败: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("关闭数据库连接失败: %w", err)
	}

	return nil
}
