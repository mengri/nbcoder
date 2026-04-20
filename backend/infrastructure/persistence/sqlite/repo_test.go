package sqlite

import (
	"context"
	"testing"

	"github.com/mengri/nbcoder/domain/project"
	"github.com/mengri/nbcoder/domain/requirement"
	"github.com/mengri/nbcoder/infrastructure/database"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *database.DatabaseManager {
	config := &database.DatabaseConfig{
		DSN:                ":memory:",
		MaxOpenConnections: 5,
		MaxIdleConnections: 2,
		ConnectionMaxLife:  60,
		DisableForeignKey:  false,
		DisableAutoMigrate: false,
		LogLevel:           4,
	}

	dbManager, err := database.NewDatabaseManager(config)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	err = database.InitSchema(dbManager.GetDB())
	if err != nil {
		t.Fatalf("Failed to init schema: %v", err)
	}

	return dbManager
}

func TestProjectRepo_Save(t *testing.T) {
	dbManager := setupTestDB(t)
	defer dbManager.Close()

	repo := NewProjectRepo(dbManager.GetDB())

	testProject := project.NewProject("test-id", "Test Project", "Description", "https://github.com/test/repo")

	err := repo.Save(testProject)
	if err != nil {
		t.Fatalf("Failed to save project: %v", err)
	}

	found, err := repo.FindByID(testProject.ID)
	if err != nil {
		t.Fatalf("Failed to find project: %v", err)
	}

	if found == nil {
		t.Fatal("Expected to find project")
	}

	if found.Name != testProject.Name {
		t.Errorf("Expected name %s, got %s", testProject.Name, found.Name)
	}
}

func TestProjectRepo_FindByID(t *testing.T) {
	dbManager := setupTestDB(t)
	defer dbManager.Close()

	repo := NewProjectRepo(dbManager.GetDB())

	testProject := project.NewProject("test-id-2", "Test Project 2", "Description", "https://github.com/test/repo")
	repo.Save(testProject)

	found, err := repo.FindByID(testProject.ID)
	if err != nil {
		t.Fatalf("Failed to find project: %v", err)
	}

	if found == nil {
		t.Fatal("Expected to find project")
	}

	if found.ID != testProject.ID {
		t.Errorf("Expected ID %s, got %s", testProject.ID, found.ID)
	}
}

func TestProjectRepo_FindAll(t *testing.T) {
	dbManager := setupTestDB(t)
	defer dbManager.Close()

	repo := NewProjectRepo(dbManager.GetDB())

	project1 := project.NewProject("test-id-3", "Test Project 3", "Description 3", "https://github.com/test/repo3")
	project2 := project.NewProject("test-id-4", "Test Project 4", "Description 4", "https://github.com/test/repo4")

	repo.Save(project1)
	repo.Save(project2)

	projects, err := repo.FindAll()
	if err != nil {
		t.Fatalf("Failed to find all projects: %v", err)
	}

	if len(projects) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(projects))
	}
}

func TestProjectRepo_Update(t *testing.T) {
	dbManager := setupTestDB(t)
	defer dbManager.Close()

	repo := NewProjectRepo(dbManager.GetDB())

	testProject := project.NewProject("test-id-5", "Test Project 5", "Description", "https://github.com/test/repo")
	repo.Save(testProject)

	testProject.Update("Updated Name", "Updated Description", "https://github.com/test/updated-repo")
	err := repo.Update(testProject)
	if err != nil {
		t.Fatalf("Failed to update project: %v", err)
	}

	found, err := repo.FindByID(testProject.ID)
	if err != nil {
		t.Fatalf("Failed to find updated project: %v", err)
	}

	if found.Name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got %s", found.Name)
	}
}

func TestProjectRepo_Delete(t *testing.T) {
	dbManager := setupTestDB(t)
	defer dbManager.Close()

	repo := NewProjectRepo(dbManager.GetDB())

	testProject := project.NewProject("test-id-6", "Test Project 6", "Description", "https://github.com/test/repo")
	repo.Save(testProject)

	err := repo.Delete(testProject.ID)
	if err != nil {
		t.Fatalf("Failed to delete project: %v", err)
	}

	found, err := repo.FindByID(testProject.ID)
	if err != nil {
		t.Fatalf("Failed to find project: %v", err)
	}

	if found != nil {
		t.Error("Expected project to be deleted")
	}
}

func TestProjectRepo_FindByStatus(t *testing.T) {
	dbManager := setupTestDB(t)
	defer dbManager.Close()

	repo := NewProjectRepo(dbManager.GetDB())

	project1 := project.NewProject("test-id-7", "Test Project 7", "Description", "https://github.com/test/repo")
	project2 := project.NewProject("test-id-8", "Test Project 8", "Description", "https://github.com/test/repo")

	project2.Archive()
	repo.Save(project1)
	repo.Save(project2)

	activeProjects, err := repo.FindByStatus(project.ProjectActive)
	if err != nil {
		t.Fatalf("Failed to find active projects: %v", err)
	}

	if len(activeProjects) != 1 {
		t.Errorf("Expected 1 active project, got %d", len(activeProjects))
	}
}

func TestCardRepo_Save(t *testing.T) {
	dbManager := setupTestDB(t)
	defer dbManager.Close()

	repo := NewCardRepo(dbManager.GetDB())

	testCard := requirement.NewCard("test-card-1", "Test Card", "Description", "Original", "test-project-id")

	err := repo.Save(testCard)
	if err != nil {
		t.Fatalf("Failed to save card: %v", err)
	}

	found, err := repo.FindByID(testCard.ID)
	if err != nil {
		t.Fatalf("Failed to find card: %v", err)
	}

	if found == nil {
		t.Fatal("Expected to find card")
	}

	if found.Title != testCard.Title {
		t.Errorf("Expected title %s, got %s", testCard.Title, found.Title)
	}
}

func TestRepositories(t *testing.T) {
	dbManager := setupTestDB(t)
	defer dbManager.Close()

	repos := NewRepositories(dbManager.GetDB())

	if repos.Project == nil {
		t.Error("Expected non-nil Project repository")
	}

	if repos.Card == nil {
		t.Error("Expected non-nil Card repository")
	}

	if repos.Task == nil {
		t.Error("Expected non-nil Task repository")
	}

	if repos.Pipeline == nil {
		t.Error("Expected non-nil Pipeline repository")
	}

	if repos.Document == nil {
		t.Error("Expected non-nil Document repository")
	}

	if repos.Notification == nil {
		t.Error("Expected non-nil Notification repository")
	}
}

func TestTransactionWithRepository(t *testing.T) {
	dbManager := setupTestDB(t)
	defer dbManager.Close()

	tm := database.NewTransactionManager(dbManager.GetDB())
	repo := NewProjectRepo(dbManager.GetDB())

	err := tm.Execute(t.Context(), func(ctx context.Context, txDB *gorm.DB) error {
		testProject := project.NewProject("test-id-9", "Test Project 9", "Description", "https://github.com/test/repo")
		return repo.Save(testProject)
	})

	if err != nil {
		t.Fatalf("Failed to execute transaction: %v", err)
	}

	found, err := repo.FindByID("test-id-9")
	if err != nil {
		t.Fatalf("Failed to find project: %v", err)
	}

	if found == nil {
		t.Fatal("Expected to find project after transaction")
	}
}
