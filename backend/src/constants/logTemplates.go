package constants

var LogTemplates = map[string][]string{
	"FlightTurnaround": {"航班过站创建", "航班过站更新", "航班过站状态变更", "航班过站导出"},
	"GroundTask":       {"地勤任务创建", "地勤任务更新", "地勤任务状态变更", "地勤任务导出", "地勤任务签收"},
	"GroundResource":   {"保障资源创建", "保障资源更新", "保障资源状态变更", "保障资源导出"},
	"ResourceBooking":  {"资源预约创建", "资源预约更新", "资源预约状态变更", "资源预约导出"},
	"DelayEvent":       {"延误事件创建", "延误事件更新", "延误事件状态变更", "延误事件导出", "延误事件登记", "延误事件关闭"},
}

// 写操作日志模板（含占位符），调用处用 log.Printf 输出
const LogGroundTaskSignOff = "[ground-turn] 地勤任务签收 task=%d signed_at=%s"
const LogDelayEventRegister = "[ground-turn] 延误事件登记 delay=%d turnaround=%d minutes=%d"
const LogDelayEventClose = "[ground-turn] 延误事件关闭 delay=%d closed_at=%s"
