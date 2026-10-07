package repositories

import (
	"sync"

	"groundTurn/src/models"
)

// 内存种子与 frontend/src/mocks/seedData.ts 保持一致。
var delayEventMu sync.RWMutex

var delayEventRows = []models.DelayEvent{
	{ID: 1, TurnaroundID: 1, DelayType: "REFUEL", Minutes: 30, RootCause: "加油车调度延迟", ResponsibilityTeam: "机务一队", ResolvedAt: ""},
	{ID: 2, TurnaroundID: 1, DelayType: "CATERING", Minutes: 15, RootCause: "配餐车晚到", ResponsibilityTeam: "配餐组", ResolvedAt: "2026-10-07T06:40:00Z"},
	{ID: 3, TurnaroundID: 2, DelayType: "BAGGAGE", Minutes: 45, RootCause: "行李分拣系统故障", ResponsibilityTeam: "行李组", ResolvedAt: ""},
}

func ListDelayEvents() []models.DelayEvent {
	delayEventMu.RLock()
	defer delayEventMu.RUnlock()
	rows := make([]models.DelayEvent, len(delayEventRows))
	copy(rows, delayEventRows)
	return rows
}

func FindDelayEvent(id int) (models.DelayEvent, bool) {
	delayEventMu.RLock()
	defer delayEventMu.RUnlock()
	for _, row := range delayEventRows {
		if row.ID == id {
			return row, true
		}
	}
	return models.DelayEvent{}, false
}

func NextDelayEventID() int {
	delayEventMu.RLock()
	defer delayEventMu.RUnlock()
	maxID := 0
	for _, row := range delayEventRows {
		if row.ID > maxID {
			maxID = row.ID
		}
	}
	return maxID + 1
}

func CreateDelayEvent(event models.DelayEvent) {
	delayEventMu.Lock()
	defer delayEventMu.Unlock()
	delayEventRows = append(delayEventRows, event)
}

func UpdateDelayEvent(event models.DelayEvent) {
	delayEventMu.Lock()
	defer delayEventMu.Unlock()
	for i, row := range delayEventRows {
		if row.ID == event.ID {
			delayEventRows[i] = event
			return
		}
	}
}
