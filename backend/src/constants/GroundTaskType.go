package constants

// GroundTaskType 地勤任务类型枚举。
// 出现位置：models 校验、types、constructors、logTemplates、errorMessages、
// services 校验、前端筛选器与展示组件。
var GroundTaskType = []string{"CLEANING", "CATERING", "BAGGAGE", "REFUEL", "WATER_SERVICE", "PUSHBACK"}

// GroundTaskStatus 任务流转状态：PENDING -> ACCEPTED -> IN_PROGRESS -> COMPLETED；可置 BLOCKED。
var GroundTaskStatus = []string{"PENDING", "ACCEPTED", "IN_PROGRESS", "BLOCKED", "COMPLETED"}

const (
	TaskStatusPending    = "PENDING"
	TaskStatusAccepted   = "ACCEPTED"
	TaskStatusInProgress = "IN_PROGRESS"
	TaskStatusBlocked    = "BLOCKED"
	TaskStatusCompleted  = "COMPLETED"
)

// IsFinishedStatus 完成过的任务按实际完成时间比超时。
func IsFinishedStatus(status string) bool {
	return status == TaskStatusCompleted
}
