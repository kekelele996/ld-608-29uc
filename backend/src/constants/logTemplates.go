package constants

// 日志模板集中定义。字段变更时必须同步改调用处（每个写动作至少一条）。
// 占位符使用 text/template 风格 {{.X}}，由 utils.RenderLog 填充。
var LogTemplates = map[string]string{
	// FlightTurnaround（≥4）
	"TURNAROUND_REGISTER": "航班过站登记 {{.FlightNo}} 机位 {{.StandNo}} 计划离港 {{.DepartureTime}}",
	"TURNAROUND_UPDATE":   "航班 {{.FlightNo}} 状态 {{.FromStatus}} -> {{.ToStatus}}",
	"TURNAROUND_RELEASE":  "航班 {{.FlightNo}} 放行，实际离港 {{.DepartureTime}}",
	"TURNAROUND_DELAYED":  "航班 {{.FlightNo}} 标记延误，原因 {{.DelayReason}}",

	// GroundTask（≥4）
	"TASK_DISPATCH":      "派工 #{{.TaskID}} {{.TaskType}} 派发至班组 {{.TeamID}}，原计划截止 {{.Deadline}}",
	"TASK_ACCEPT":        "任务 #{{.TaskID}} 签收，签收时间 {{.AcceptedAt}}，当前生效截止 {{.EffectiveDeadline}}",
	"TASK_COMPLETE":      "任务 #{{.TaskID}} 完成，实际完成 {{.ActualFinish}}，生效截止 {{.EffectiveDeadline}}，超时 {{.Overdue}}",
	"TASK_BLOCK":         "任务 #{{.TaskID}} 阻塞，原因 {{.BlockerNote}}",
	"TASK_DEADLINE_VIEW": "任务 #{{.TaskID}} 截止口径：原计划 {{.Deadline}}，顺延 {{.OpenDelayMinutes}} 分钟，生效截止 {{.EffectiveDeadline}}",

	// GroundResource（≥4）
	"RESOURCE_CREATE":  "资源台账新增 {{.ResourceCode}}（{{.ResourceType}}）",
	"RESOURCE_UPDATE":  "资源 {{.ResourceCode}} 状态 -> {{.AvailabilityStatus}}",
	"RESOURCE_MAINT":   "资源 {{.ResourceCode}} 进入检修，计划复用于 {{.MaintenanceDueAt}}",
	"RESOURCE_OFFLINE": "资源 {{.ResourceCode}} 下线 {{.Reason}}",

	// ResourceBooking（≥4）
	"BOOKING_CREATE":   "预约 #{{.BookingID}} 资源 {{.ResourceID}} 窗口 {{.StartTime}}~{{.EndTime}}",
	"BOOKING_CONFIRM":  "预约 #{{.BookingID}} 确认",
	"BOOKING_RELEASE":  "预约 #{{.BookingID}} 释放",
	"BOOKING_CONFLICT": "预约冲突：资源 {{.ResourceID}} 与预约 #{{.ConflictID}} 时间重叠（{{.ConflictReason}}）",

	// DelayEvent（≥4）
	"DELAY_REGISTER":        "航班 #{{.TurnaroundID}} 登记延误 {{.DelayType}} {{.Minutes}} 分钟，未关闭延误累计顺延 {{.OpenDelayMinutes}} 分钟，责任 {{.ResponsibilityTeam}}",
	"DELAY_RESOLVE":         "延误 #{{.DelayID}} 关闭，关闭时间 {{.ResolvedAt}}，该 {{.Minutes}} 分钟不再顺延",
	"DELAY_ATTRIBUTE":       "延误 #{{.DelayID}} 归因变更为 {{.ResponsibilityTeam}}/{{.RootCause}}",
	"DELAY_DEADLINE_RECALC": "航班 #{{.TurnaroundID}} 延误变化触发任务截止重算：未关闭 {{.OpenCount}} 起，顺延 {{.OpenDelayMinutes}} 分钟",

	// 鉴权
	"AUTH_LOGIN": "用户 {{.Username}}（{{.Role}}）登录签发令牌",
}
