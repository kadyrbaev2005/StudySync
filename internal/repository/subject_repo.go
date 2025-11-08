package repository

import (
	"github.com/kadyrbayev2005/studysync/internal/models"

	"gorm.io/gorm"
)

type SubjectRepository struct {
	db *gorm.DB
}

func NewSubjectRepository(db *gorm.DB) *SubjectRepository {
	return &SubjectRepository{db}
}

func (r *SubjectRepository) Create(subject *models.Subject) error {
	return r.db.Create(subject).Error
}

func (r *SubjectRepository) GetAll() ([]models.Subject, error) {
	var subjects []models.Subject
	err := r.db.Find(&subjects).Error
	return subjects, err
}

func (r *SubjectRepository) GetByID(id uint) (models.Subject, error) {
	var subject models.Subject
	err := r.db.First(&subject, id).Error
	return subject, err
}

func (r *SubjectRepository) Update(id uint, data map[string]interface{}) error {
	return r.db.Model(&models.Subject{}).Where("id = ?", id).Updates(data).Error
}

func (r *SubjectRepository) Delete(id uint) error {
	return r.db.Delete(&models.Subject{}, id).Error
}
