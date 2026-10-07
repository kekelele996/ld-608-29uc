package repositories

import (
	"groundTurn/src/models"

	"gorm.io/gorm"
)

type FlightTurnaroundRepository struct{ DB *gorm.DB }

func NewFlightTurnaroundRepository(db *gorm.DB) *FlightTurnaroundRepository {
	return &FlightTurnaroundRepository{DB: db}
}

func (r *FlightTurnaroundRepository) List() ([]models.FlightTurnaround, error) {
	var rows []models.FlightTurnaround
	err := r.DB.Order("arrival_time asc").Find(&rows).Error
	return rows, err
}

func (r *FlightTurnaroundRepository) Get(id uint) (*models.FlightTurnaround, error) {
	var row models.FlightTurnaround
	if err := r.DB.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *FlightTurnaroundRepository) Create(row *models.FlightTurnaround) error {
	return r.DB.Create(row).Error
}

func (r *FlightTurnaroundRepository) UpdateStatus(id uint, status, delayReason string) error {
	return r.DB.Model(&models.FlightTurnaround{}).Where("id = ?", id).
		Updates(map[string]interface{}{"turnaround_status": status, "delay_reason": delayReason}).Error
}
