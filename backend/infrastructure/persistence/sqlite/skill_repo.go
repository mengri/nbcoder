package sqlite

import (
	"fmt"

	"github.com/mengri/nbcoder/domain/agent"
	"github.com/mengri/nbcoder/infrastructure/database/models"
	"gorm.io/gorm"
)

type SkillRepo struct {
	dbProvider DBProvider
}

func NewSkillRepo(dbProvider DBProvider) agent.SkillRepo {
	return &SkillRepo{dbProvider: dbProvider}
}

func (r *SkillRepo) getDB() (*gorm.DB, error) {
	return r.dbProvider.GetGlobalDB(), nil
}

func (r *SkillRepo) Save(skill *agent.Skill) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	model := &models.Skill{
		ID:          skill.ID,
		Name:        skill.Name,
		Description: skill.Description,
		AgentType:   string(skill.AgentType),
		Config:      models.JSONMap(skill.Config),
		CreatedAt:   skill.CreatedAt,
		UpdatedAt:   skill.UpdatedAt,
	}

	result := db.Save(model)
	if result.Error != nil {
		return fmt.Errorf("failed to save skill: %w", result.Error)
	}
	return nil
}

func (r *SkillRepo) FindByID(id string) (*agent.Skill, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var model models.Skill
	result := db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find skill by id: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *SkillRepo) FindByAgentType(agentType agent.AgentType) ([]*agent.Skill, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []models.Skill
	result := db.Where("agent_type = ?", string(agentType)).Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find skills by agent type: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *SkillRepo) FindAll() ([]*agent.Skill, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []models.Skill
	result := db.Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find all skills: %w", result.Error)
	}

	return r.modelsToDomain(models), nil
}

func (r *SkillRepo) Update(skill *agent.Skill) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	model := &models.Skill{
		ID:          skill.ID,
		Name:        skill.Name,
		Description: skill.Description,
		AgentType:   string(skill.AgentType),
		Config:      models.JSONMap(skill.Config),
		CreatedAt:   skill.CreatedAt,
		UpdatedAt:   skill.UpdatedAt,
	}

	result := db.Model(&models.Skill{}).Where("id = ?", skill.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update skill: %w", result.Error)
	}
	return nil
}

func (r *SkillRepo) Delete(id string) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	result := db.Delete(&models.Skill{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete skill: %w", result.Error)
	}
	return nil
}

func (r *SkillRepo) FindByName(name string) (*agent.Skill, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var model models.Skill
	result := db.Where("name = ?", name).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find skill by name: %w", result.Error)
	}

	return r.modelToDomain(&model), nil
}

func (r *SkillRepo) modelToDomain(m *models.Skill) *agent.Skill {
	return &agent.Skill{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		AgentType:   agent.AgentType(m.AgentType),
		Config:      map[string]interface{}(m.Config),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (r *SkillRepo) modelsToDomain(models []models.Skill) []*agent.Skill {
	result := make([]*agent.Skill, len(models))
	for i, m := range models {
		result[i] = r.modelToDomain(&m)
	}
	return result
}
