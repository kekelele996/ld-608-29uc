package repositories

import (
	"groundTurn/src/models"

	"gorm.io/gorm"
)

type UserRepository struct{ DB *gorm.DB }

func NewUserRepository(db *gorm.DB) *UserRepository { return &UserRepository{DB: db} }

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var row models.User
	if err := r.DB.Where("username = ?", username).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *UserRepository) Count() (int64, error) {
	var n int64
	err := r.DB.Model(&models.User{}).Count(&n).Error
	return n, err
}

func (r *UserRepository) Create(row *models.User) error { return r.DB.Create(row).Error }

type AuditLogRepository struct{ DB *gorm.DB }

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository { return &AuditLogRepository{DB: db} }

func (r *AuditLogRepository) Create(row *models.AuditLog) error { return r.DB.Create(row).Error }

func (r *AuditLogRepository) List(limit int) ([]models.AuditLog, error) {
	var rows []models.AuditLog
	err := r.DB.Order("created_at desc").Limit(limit).Find(&rows).Error
	return rows, err
}
