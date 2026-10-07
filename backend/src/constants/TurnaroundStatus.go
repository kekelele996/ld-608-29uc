package constants

// TurnaroundStatus 航班过站状态枚举。
// 出现位置：models、types、constructors、logTemplates、errorMessages、
// services、看板筛选器、StatusBadge/TurnaroundTimeline。
var TurnaroundStatus = []string{"ARRIVING", "ON_STAND", "IN_SERVICE", "READY", "DEPARTED", "DELAYED"}

const (
	TurnaroundArriving  = "ARRIVING"
	TurnaroundOnStand   = "ON_STAND"
	TurnaroundInService = "IN_SERVICE"
	TurnaroundReady     = "READY"
	TurnaroundDeparted  = "DEPARTED"
	TurnaroundDelayed   = "DELAYED"
)

// DelayType 延误事件分类（与任务类型部分重合，另含天气/流量）。
var DelayType = []string{"CLEANING", "CATERING", "BAGGAGE", "REFUEL", "WATER_SERVICE", "PUSHBACK", "WEATHER", "ATC_FLOW"}
