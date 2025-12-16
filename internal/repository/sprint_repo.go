package repository

import (
    "github.com/kadyrbayev2005/studysync/internal/models"
    "gorm.io/gorm"
)

type SprintRepository struct {
    db *gorm.DB
}

func NewSprintRepository(db *gorm.DB) *SprintRepository {
    return &SprintRepository{db}
}

func (r *SprintRepository) Create(sprint *models.Sprint) error {
    return r.db.Create(sprint).Error
}

func (r *SprintRepository) GetAll() ([]models.Sprint, error) {
    var sprints []models.Sprint
    err := r.db.Preload("Tasks").Find(&sprints).Error
    return sprints, err
}

func (r *SprintRepository) GetByID(id uint) (models.Sprint, error) {
    var sprint models.Sprint
    err := r.db.Preload("Tasks").First(&sprint, id).Error
    return sprint, err
}

func (r *SprintRepository) Update(id uint, data map[string]interface{}) error {
    return r.db.Model(&models.Sprint{}).Where("id = ?", id).Updates(data).Error
}

func (r *SprintRepository) Delete(id uint) error {
    return r.db.Delete(&models.Sprint{}, id).Error
}

func (r *SprintRepository) GetActiveSprints() ([]models.Sprint, error) {
    var sprints []models.Sprint
    err := r.db.Where("status = ?", "active").Preload("Tasks").Find(&sprints).Error
    return sprints, err
}