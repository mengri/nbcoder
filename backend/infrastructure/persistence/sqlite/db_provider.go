package sqlite

import (
	"gorm.io/gorm"
)

type DBProvider interface {
	GetGlobalDB() *gorm.DB
	GetProjectDB(projectName string) (*gorm.DB, error)
}
