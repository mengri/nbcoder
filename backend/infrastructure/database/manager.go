package database

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/mengri/nbcoder/infrastructure/persistence/sqlite"
	sqlitedriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBManager interface {
	GetGlobalDB() *gorm.DB
	GetProjectDB(projectName string) (*gorm.DB, error)
}

type DatabaseConfig struct {
	DSN                  string
	MaxOpenConnections   int
	MaxIdleConnections   int
	ConnectionMaxLife    time.Duration
	MaxIdleTime          time.Duration
	DisableForeignKey    bool
	DisableAutoMigrate   bool
	LogLevel             logger.LogLevel
	ProjectBaseDir       string
}

type DatabaseManager struct {
	db            *gorm.DB
	config        *DatabaseConfig
	projectDBs    map[string]*gorm.DB
	projectDBsMux sync.RWMutex
}

func NewDatabaseManager(config *DatabaseConfig) (*DatabaseManager, error) {
	if config == nil {
		config = DefaultDatabaseConfig()
	}

	dbDir := filepath.Dir(config.DSN)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := gorm.Open(sqlitedriver.Open(config.DSN), &gorm.Config{
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
		db:         db,
		config:     config,
		projectDBs: make(map[string]*gorm.DB),
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
		ProjectBaseDir:     "./projects",
	}
}

func (dm *DatabaseManager) GetDB() *gorm.DB {
	return dm.db
}

func (dm *DatabaseManager) GetGlobalDB() *gorm.DB {
	return dm.db
}

var _ sqlite.DBProvider = (*DatabaseManager)(nil)

func (dm *DatabaseManager) GetProjectDB(projectName string) (*gorm.DB, error) {
	dm.projectDBsMux.RLock()
	db, exists := dm.projectDBs[projectName]
	dm.projectDBsMux.RUnlock()

	if exists {
		return db, nil
	}

	dm.projectDBsMux.Lock()
	defer dm.projectDBsMux.Unlock()

	dbPath := filepath.Join(dm.config.ProjectBaseDir, projectName, "nbcoder.db")
	projectDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create project database directory: %w", err)
	}

	db, err := gorm.Open(sqlitedriver.Open(dbPath), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: dm.config.DisableForeignKey,
		Logger: logger.Default.LogMode(dm.config.LogLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to project database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get project database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(dm.config.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(dm.config.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(dm.config.ConnectionMaxLife)
	sqlDB.SetConnMaxIdleTime(dm.config.MaxIdleTime)

	dm.projectDBs[projectName] = db

	return db, nil
}

func (dm *DatabaseManager) Close() error {
	sqlDB, err := dm.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	dm.projectDBsMux.Lock()
	defer dm.projectDBsMux.Unlock()

	for projectName, projectDB := range dm.projectDBs {
		if psqlDB, err := projectDB.DB(); err == nil {
			psqlDB.Close()
		} else {
			fmt.Printf("failed to get project database instance for %s: %v\n", projectName, err)
		}
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
