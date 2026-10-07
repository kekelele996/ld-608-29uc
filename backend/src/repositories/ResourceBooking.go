package repositories

import (
	"time"

	"groundTurn/src/models"

	"gorm.io/gorm"
)

type ResourceBookingRepository struct{ DB *gorm.DB }

func NewResourceBookingRepository(db *gorm.DB) *ResourceBookingRepository {
	return &ResourceBookingRepository{DB: db}
}

func (r *ResourceBookingRepository) List() ([]models.ResourceBooking, error) {
	var rows []models.ResourceBooking
	err := r.DB.Order("start_time asc").Find(&rows).Error
	return rows, err
}

func (r *ResourceBookingRepository) Get(id uint) (*models.ResourceBooking, error) {
	var row models.ResourceBooking
	if err := r.DB.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// FindOverlap 查询同一资源、未释放且时间窗重叠的预约（冲突检测）。
func (r *ResourceBookingRepository) FindOverlap(resourceID uint, start, end time.Time, excludeID uint) ([]models.ResourceBooking, error) {
	var rows []models.ResourceBooking
	err := r.DB.Where("resource_id = ? AND id <> ? AND booking_status <> ?", resourceID, excludeID, "RELEASED").
		Where("start_time < ? AND end_time > ?", end, start).
		Find(&rows).Error
	return rows, err
}

func (r *ResourceBookingRepository) Create(row *models.ResourceBooking) error {
	return r.DB.Create(row).Error
}

func (r *ResourceBookingRepository) UpdateStatus(id uint, status, conflictReason string) error {
	return r.DB.Model(&models.ResourceBooking{}).Where("id = ?", id).
		Updates(map[string]interface{}{"booking_status": status, "conflict_reason": conflictReason}).Error
}
