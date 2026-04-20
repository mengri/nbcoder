package database

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseConfig struct {
	DSN                  string
	MaxOpenConnections   int
	MaxIdleConnections   int
	ConnectionMaxLife    time.Duration
	MaxIdleTime          time.Duration
	DisableForeignKey    bool
	DisableAutoMigrate   bool
	LogLevel             logger.LogLevel
}

type DatabaseManager struct {
	db     *gorm.DB
	config *DatabaseConfig
}

func NewDatabaseManager(config *DatabaseConfig) (*DatabaseManager, error) {
	if config == nil {
		config = DefaultDatabaseConfig()
	}

	dbDir := filepath.Dir(config.DSN)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := gorm.Open(sqlite.Open(config.DSN), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: config.DisableForeignKey,
		Logger: logger.Default.LogMode(config.LogLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(config.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(config.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(config.ConnectionMaxLife)
	sqlDB.SetConnMaxIdleTime(config.MaxIdleTime)

	return &DatabaseManager{
		db:     db,
		config: config,
	}, nil
}

func DefaultDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		DSN:                ".NBCoder-global/nbcoder.db",
		MaxOpenConnections: 25,
		MaxIdleConnections: 5,
		ConnectionMaxLife:  time.Hour,
		MaxIdleTime:        10 * time.Minute,
		DisableForeignKey:  false,
		DisableAutoMigrate: false,
		LogLevel:           logger.Info,
	}
}

func (dm *DatabaseManager) GetDB() *gorm.DB {
	return dm.db
}

func (dm *DatabaseManager) Close() error {
	sqlDB, err := dm.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	return sqlDB.Close()
}

func (dm *DatabaseManager) Ping() error {
	sqlDB, err := dm.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	return sqlDB.Ping()
}
