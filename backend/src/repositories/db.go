package repositories

import (
	"time"

	"groundTurn/src/config"
	"groundTurn/src/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect 打开 GORM 连接并自动建表（init.sql 已先建库/表，AutoMigrate 负责补齐列）。
func Connect(cfg config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger:  logger.Default.LogMode(logger.Warn),
		NowFunc: func() time.Time { return time.Now().UTC() },
	})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(
		&models.FlightTurnaround{},
		&models.GroundTask{},
		&models.GroundResource{},
		&models.ResourceBooking{},
		&models.DelayEvent{},
		&models.User{},
		&models.AuditLog{},
	); err != nil {
		return nil, err
	}
	return db, nil
}
