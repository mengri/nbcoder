package database

import (
	"context"
	"database/sql"
	"fmt"

	"gorm.io/gorm"
)

type TransactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) *TransactionManager {
	return &TransactionManager{
		db: db,
	}
}

type TransactionalFunc func(ctx context.Context, db *gorm.DB) error

func (tm *TransactionManager) Execute(ctx context.Context, fn TransactionalFunc) error {
	return tm.ExecuteInTx(ctx, nil, fn)
}

func (tm *TransactionManager) ExecuteInTx(ctx context.Context, opts *sql.TxOptions, fn TransactionalFunc) error {
	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx, tx)
	}, opts)
}

func (tm *TransactionManager) WithDB(db *gorm.DB) *TransactionManager {
	return &TransactionManager{
		db: db,
	}
}

func (tm *TransactionManager) GetDB() *gorm.DB {
	return tm.db
}

func (tm *TransactionManager) BeginTx(opts *sql.TxOptions) *gorm.DB {
	return tm.db.Begin(opts)
}

func (tm *TransactionManager) Rollback(tx *gorm.DB) error {
	if tx.Error == nil {
		return tx.Rollback().Error
	}
	return tx.Error
}

func (tm *TransactionManager) Commit(tx *gorm.DB) error {
	if tx.Error == nil {
		return tx.Commit().Error
	}
	return fmt.Errorf("cannot commit transaction: %w", tx.Error)
}
