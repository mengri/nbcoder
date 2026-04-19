package migrations

import (
	"embed"
	"fmt"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	ErrMigrationAlreadyApplied = fmt.Errorf("migration already applied")
	ErrMigrationNotFound       = fmt.Errorf("migration not found")
)

type Migration struct {
	ID        string    `gorm:"type:varchar(100);primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	AppliedAt time.Time `gorm:"autoCreateTime" json:"applied_at"`
}

func (Migration) TableName() string {
	return "schema_migrations"
}

type Migrator struct {
	db *gorm.DB
	fs *embed.FS
}

func NewMigrator(db *gorm.DB, fs *embed.FS) *Migrator {
	return &Migrator{
		db: db,
		fs: fs,
	}
}

func (m *Migrator) Init() error {
	return m.db.AutoMigrate(&Migration{})
}

func (m *Migrator) Status() ([]Migration, error) {
	var migrations []Migration
	err := m.db.Order("applied_at ASC").Find(&migrations).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get migration status: %w", err)
	}
	return migrations, nil
}

func (m *Migrator) Up() error {
	appliedMigrations, err := m.Status()
	if err != nil {
		return err
	}

	appliedMap := make(map[string]bool)
	for _, m := range appliedMigrations {
		appliedMap[m.ID] = true
	}

	files, err := m.listMigrationFiles()
	if err != nil {
		return err
	}

	sort.Strings(files)

	for _, file := range files {
		if strings.HasSuffix(file, ".up.sql") {
			migrationID := strings.TrimSuffix(file, ".up.sql")

			if appliedMap[migrationID] {
				continue
			}

			err := m.applyMigration(migrationID, file)
			if err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", migrationID, err)
			}
		}
	}

	return nil
}

func (m *Migrator) Down() error {
	appliedMigrations, err := m.Status()
	if err != nil {
		return err
	}

	if len(appliedMigrations) == 0 {
		return nil
	}

	lastMigration := appliedMigrations[len(appliedMigrations)-1]
	downFile := lastMigration.ID + ".down.sql"

	err = m.applyDownMigration(lastMigration.ID, downFile)
	if err != nil {
		return fmt.Errorf("failed to revert migration %s: %w", lastMigration.ID, err)
	}

	err = m.db.Delete(&Migration{}, "id = ?", lastMigration.ID).Error
	if err != nil {
		return fmt.Errorf("failed to delete migration record: %w", err)
	}

	return nil
}

func (m *Migrator) CreateMigration(id, name string) error {
	migration := Migration{
		ID:   id,
		Name: name,
	}

	result := m.db.Create(&migration)
	if result.Error != nil {
		return fmt.Errorf("failed to create migration record: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrMigrationAlreadyApplied
	}

	return nil
}

func (m *Migrator) listMigrationFiles() ([]string, error) {
	if m.fs == nil {
		return []string{}, nil
	}

	var files []string
	entries, err := m.fs.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("failed to read migration directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && (strings.HasSuffix(entry.Name(), ".up.sql") || strings.HasSuffix(entry.Name(), ".down.sql")) {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

func (m *Migrator) applyMigration(id, filename string) error {
	content, err := m.readMigrationFile(filename)
	if err != nil {
		return err
	}

	if content == "" {
		return fmt.Errorf("migration file %s is empty", filename)
	}

	err = m.executeSQL(content)
	if err != nil {
		return err
	}

	return m.CreateMigration(id, filename)
}

func (m *Migrator) applyDownMigration(id, filename string) error {
	content, err := m.readMigrationFile(filename)
	if err != nil {
		return err
	}

	if content == "" {
		return fmt.Errorf("migration file %s is empty", filename)
	}

	return m.executeSQL(content)
}

func (m *Migrator) readMigrationFile(filename string) (string, error) {
	if m.fs == nil {
		return "", fmt.Errorf("no embedded filesystem")
	}

	content, err := m.fs.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read migration file %s: %w", filename, err)
	}

	return string(content), nil
}

func (m *Migrator) executeSQL(sql string) error {
	db, err := m.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	_, err = db.Exec(sql)
	if err != nil {
		return fmt.Errorf("failed to execute SQL: %w", err)
	}

	return nil
}

func (m *Migrator) IsApplied(id string) (bool, error) {
	var count int64
	err := m.db.Model(&Migration{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check migration status: %w", err)
	}
	return count > 0, nil
}
