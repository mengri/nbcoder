package database

import (
	"context"
	"os"
	"testing"

	"gorm.io/gorm"
)

type TestModel struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

func setupTestDB(t *testing.T) *DatabaseManager {
	config := &DatabaseConfig{
		DSN:                ":memory:",
		MaxOpenConnections: 5,
		MaxIdleConnections: 2,
		ConnectionMaxLife:  60,
		DisableForeignKey:  true,
		DisableAutoMigrate: false,
		LogLevel:           4,
	}

	dbManager, err := NewDatabaseManager(config)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	return dbManager
}

func TestNewDatabaseManager(t *testing.T) {
	dbManager := setupTestDB(t)
	defer dbManager.Close()

	if dbManager == nil {
		t.Fatal("Expected non-nil database manager")
	}

	if dbManager.GetDB() == nil {
		t.Fatal("Expected non-nil database connection")
	}

	if err := dbManager.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}

func TestDatabaseManager_Ping(t *testing.T) {
	dbManager := setupTestDB(t)
	defer dbManager.Close()

	err := dbManager.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}

func TestDatabaseManager_Close(t *testing.T) {
	dbManager := setupTestDB(t)

	err := dbManager.Close()
	if err != nil {
		t.Fatalf("Failed to close database: %v", err)
	}
}

func TestDefaultDatabaseConfig(t *testing.T) {
	config := DefaultDatabaseConfig()

	if config.DSN == "" {
		t.Error("Expected non-empty DSN")
	}

	if config.MaxOpenConnections != 25 {
		t.Errorf("Expected MaxOpenConnections to be 25, got %d", config.MaxOpenConnections)
	}

	if config.MaxIdleConnections != 5 {
		t.Errorf("Expected MaxIdleConnections to be 5, got %d", config.MaxIdleConnections)
	}
}

func TestTransactionManager_Execute(t *testing.T) {
	dbManager := setupTestDB(t)
	defer dbManager.Close()

	tm := NewTransactionManager(dbManager.GetDB())

	err := tm.Execute(context.Background(), func(ctx context.Context, db *gorm.DB) error {
		return nil
	})

	if err != nil {
		t.Fatalf("Failed to execute transaction: %v", err)
	}
}

func TestTransactionManager_RollbackOnError(t *testing.T) {
	dbManager := setupTestDB(t)
	defer dbManager.Close()

	tm := NewTransactionManager(dbManager.GetDB())

	err := tm.Execute(context.Background(), func(ctx context.Context, db *gorm.DB) error {
		return gorm.ErrRecordNotFound
	})

	if err == nil {
		t.Fatal("Expected error when transaction fails")
	}
}

func TestDatabaseManager_GetDB(t *testing.T) {
	dbManager := setupTestDB(t)
	defer dbManager.Close()

	db := dbManager.GetDB()

	if db == nil {
		t.Fatal("Expected non-nil database connection")
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get SQL DB: %v", err)
	}
	defer sqlDB.Close()

	if sqlDB == nil {
		t.Fatal("Expected non-nil SQL DB")
	}
}

func TestInitSchema(t *testing.T) {
	dbManager := setupTestDB(t)
	defer dbManager.Close()

	db := dbManager.GetDB()

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get SQL DB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	if sqlDB.Stats().MaxOpenConnections != 5 {
		t.Errorf("Expected MaxOpenConnections to be 5, got %d", sqlDB.Stats().MaxOpenConnections)
	}

	if sqlDB.Stats().Idle != 0 {
		t.Errorf("Expected Idle to be 0, got %d", sqlDB.Stats().Idle)
	}
}

func cleanupTestDB(dbPath string) {
	if dbPath != ":memory:" {
		os.Remove(dbPath)
	}
}
