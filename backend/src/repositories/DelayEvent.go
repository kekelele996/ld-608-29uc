package repositories

import (
	"groundTurn/src/models"

	"gorm.io/gorm"
)

type DelayEventRepository struct{ DB *gorm.DB }

func NewDelayEventRepository(db *gorm.DB) *DelayEventRepository {
	return &DelayEventRepository{DB: db}
}

// List 全量取延误事件，顺延口径（只算 resolved_at IS NULL）在 utils 层计算。
func (r *DelayEventRepository) List() ([]models.DelayEvent, error) {
	var rows []models.DelayEvent
	err := r.DB.Order("created_at desc").Find(&rows).Error
	return rows, err
}

func (r *DelayEventRepository) Get(id uint) (*models.DelayEvent, error) {
	var row models.DelayEvent
	if err := r.DB.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// MapByTurnaround 一次取回各航班延误事件，避免在任务循环里 N+1。
func (r *DelayEventRepository) MapByTurnaround(ids []uint) (map[uint][]models.DelayEvent, error) {
	result := map[uint][]models.DelayEvent{}
	if len(ids) == 0 {
		return result, nil
	}
	var rows []models.DelayEvent
	if err := r.DB.Where("turnaround_id IN ?", ids).Order("created_at asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	for i := range rows {
		result[rows[i].TurnaroundID] = append(result[rows[i].TurnaroundID], rows[i])
	}
	return result, nil
}

func (r *DelayEventRepository) Create(row *models.DelayEvent) error { return r.DB.Create(row).Error }

func (r *DelayEventRepository) Resolve(id uint, resolvedAt interface{}) error {
	return r.DB.Model(&models.DelayEvent{}).Where("id = ?", id).
		Update("resolved_at", resolvedAt).Error
}

func (r *DelayEventRepository) Update(row *models.DelayEvent) error { return r.DB.Save(row).Error }
