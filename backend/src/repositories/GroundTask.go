package repositories

import (
	"sync"

	"groundTurn/src/models"
)

// 内存种子与 frontend/src/mocks/seedData.ts 保持一致，保证 API 与 mock 回退行为一致。
var groundTaskMu sync.RWMutex

var groundTaskRows = []models.GroundTask{
	{ID: 1, TurnaroundID: 1, TaskType: "CLEANING", TeamID: 1, PlannedStart: "2026-10-07T01:30:00Z", Deadline: "2026-10-07T02:00:00Z", ActualFinish: "", Status: "PENDING", BlockerNote: ""},
	{ID: 2, TurnaroundID: 1, TaskType: "BAGGAGE", TeamID: 2, PlannedStart: "2026-10-07T05:40:00Z", Deadline: "2026-10-07T06:20:00Z", ActualFinish: "2026-10-07T06:55:00Z", Status: "SIGNED", BlockerNote: ""},
	{ID: 3, TurnaroundID: 1, TaskType: "CATERING", TeamID: 3, PlannedStart: "2026-10-07T06:00:00Z", Deadline: "2026-10-07T06:50:00Z", ActualFinish: "2026-10-07T07:05:00Z", Status: "DONE", BlockerNote: ""},
	{ID: 4, TurnaroundID: 2, TaskType: "REFUEL", TeamID: 1, PlannedStart: "2026-10-07T07:00:00Z", Deadline: "2026-10-07T08:00:00Z", ActualFinish: "", Status: "PENDING", BlockerNote: ""},
	{ID: 5, TurnaroundID: 2, TaskType: "WATER_SERVICE", TeamID: 2, PlannedStart: "2026-10-07T06:30:00Z", Deadline: "2026-10-07T07:00:00Z", ActualFinish: "2026-10-07T07:20:00Z", Status: "SIGNED", BlockerNote: ""},
	{ID: 6, TurnaroundID: 3, TaskType: "PUSHBACK", TeamID: 3, PlannedStart: "2026-10-07T09:30:00Z", Deadline: "2026-10-07T10:00:00Z", ActualFinish: "", Status: "BLOCKED", BlockerNote: "牵引车故障待修"},
}

func ListGroundTasks() []models.GroundTask {
	groundTaskMu.RLock()
	defer groundTaskMu.RUnlock()
	rows := make([]models.GroundTask, len(groundTaskRows))
	copy(rows, groundTaskRows)
	return rows
}

func FindGroundTask(id int) (models.GroundTask, bool) {
	groundTaskMu.RLock()
	defer groundTaskMu.RUnlock()
	for _, row := range groundTaskRows {
		if row.ID == id {
			return row, true
		}
	}
	return models.GroundTask{}, false
}

func UpdateGroundTask(task models.GroundTask) {
	groundTaskMu.Lock()
	defer groundTaskMu.Unlock()
	for i, row := range groundTaskRows {
		if row.ID == task.ID {
			groundTaskRows[i] = task
			return
		}
	}
}
