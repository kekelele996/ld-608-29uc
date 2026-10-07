package repositories

import (
	"groundTurn/src/models"

	"gorm.io/gorm"
)

type GroundResourceRepository struct{ DB *gorm.DB }

func NewGroundResourceRepository(db *gorm.DB) *GroundResourceRepository {
	return &GroundResourceRepository{DB: db}
}

func (r *GroundResourceRepository) List() ([]models.GroundResource, error) {
	var rows []models.GroundResource
	err := r.DB.Order("resource_code asc").Find(&rows).Error
	return rows, err
}

func (r *GroundResourceRepository) Get(id uint) (*models.GroundResource, error) {
	var row models.GroundResource
	if err := r.DB.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *GroundResourceRepository) Create(row *models.GroundResource) error {
	return r.DB.Create(row).Error
}
func (r *GroundResourceRepository) Update(row *models.GroundResource) error {
	return r.DB.Save(row).Error
}
