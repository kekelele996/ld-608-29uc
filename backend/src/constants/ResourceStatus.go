package constants

// ResourceStatus 保障资源可用性枚举。
var ResourceStatus = []string{"AVAILABLE", "BOOKED", "MAINTENANCE", "OFFLINE"}

const (
	ResourceAvailable   = "AVAILABLE"
	ResourceBooked      = "BOOKED"
	ResourceMaintenance = "MAINTENANCE"
	ResourceOffline     = "OFFLINE"
)

// BookingStatus 资源预约状态。
var BookingStatus = []string{"HELD", "CONFIRMED", "RELEASED", "CONFLICT"}

const (
	BookingHeld      = "HELD"
	BookingConfirmed = "CONFIRMED"
	BookingReleased  = "RELEASED"
	BookingConflict  = "CONFLICT"
)

// RBAC 角色：地勤调度/班组/资源管理员/运行督导。
var Roles = []string{"DISPATCHER", "TEAM", "RESOURCE_MANAGER", "SUPERVISOR"}

const (
	RoleDispatcher      = "DISPATCHER"
	RoleTeam            = "TEAM"
	RoleResourceManager = "RESOURCE_MANAGER"
	RoleSupervisor      = "SUPERVISOR"
)
