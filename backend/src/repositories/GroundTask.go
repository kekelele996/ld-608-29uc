package repositories

import (
	"groundTurn/src/models"

	"gorm.io/gorm"
)

type GroundTaskRepository struct{ DB *gorm.DB }

func NewGroundTaskRepository(db *gorm.DB) *GroundTaskRepository {
	return &GroundTaskRepository{DB: db}
}

func (r *GroundTaskRepository) List() ([]models.GroundTask, error) {
	var rows []models.GroundTask
	err := r.DB.Order("deadline asc").Find(&rows).Error
	return rows, err
}

func (r *GroundTaskRepository) Get(id uint) (*models.GroundTask, error) {
	var row models.GroundTask
	if err := r.DB.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *GroundTaskRepository) ListByTurnaround(ids []uint) (map[uint][]models.GroundTask, error) {
	result := map[uint][]models.GroundTask{}
	if len(ids) == 0 {
		return result, nil
	}
	var rows []models.GroundTask
	if err := r.DB.Where("turnaround_id IN ?", ids).Order("deadline asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	for i := range rows {
		result[rows[i].TurnaroundID] = append(result[rows[i].TurnaroundID], rows[i])
	}
	return result, nil
}

func (r *GroundTaskRepository) Create(row *models.GroundTask) error { return r.DB.Create(row).Error }

func (r *GroundTaskRepository) Update(row *models.GroundTask) error {
	return r.DB.Save(row).Error
}
