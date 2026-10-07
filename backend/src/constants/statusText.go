package constants

// 状态文案集中存放，供构造器/日志/响应格式化引用。
var StatusText = map[string]map[string]string{
	"TurnaroundStatus": {
		"ARRIVING": "进近中", "ON_STAND": "靠机", "IN_SERVICE": "保障中",
		"READY": "就绪", "DEPARTED": "已离港", "DELAYED": "延误",
	},
	"GroundTaskStatus": {
		"PENDING": "待签收", "ACCEPTED": "已签收", "IN_PROGRESS": "进行中",
		"BLOCKED": "阻塞", "COMPLETED": "已完成",
	},
	"GroundTaskType": {
		"CLEANING": "客舱清洁", "CATERING": "餐食保障", "BAGGAGE": "行李装卸",
		"REFUEL": "加油", "WATER_SERVICE": "清水污水", "PUSHBACK": "推机",
	},
	"ResourceStatus": {
		"AVAILABLE": "可用", "BOOKED": "已预约", "MAINTENANCE": "检修中", "OFFLINE": "下线",
	},
	"BookingStatus": {
		"HELD": "占用", "CONFIRMED": "已确认", "RELEASED": "已释放", "CONFLICT": "冲突",
	},
	"Role": {
		"DISPATCHER": "地勤调度", "TEAM": "班组", "RESOURCE_MANAGER": "资源管理员", "SUPERVISOR": "运行督导",
	},
}
