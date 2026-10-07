package constants

// 错误码集中定义。service 与 controller 分别包装，禁止全局一处吞掉。
const (
	AuthRequired        = "AUTH_REQUIRED"
	AuthInvalid         = "AUTH_INVALID"
	Forbidden           = "FORBIDDEN"
	RateLimited         = "RATE_LIMITED"
	ValidationFailed    = "VALIDATION_FAILED"
	RecordNotFound      = "RECORD_NOT_FOUND"
	InvalidEnum         = "INVALID_ENUM"
	Conflict            = "RESOURCE_CONFLICT"
	DeadlineNotExtended = "DEADLINE_NOT_EXTENDED"
	TaskNotAcceptable   = "TASK_NOT_ACCEPTABLE"
	InternalError       = "INTERNAL_ERROR"
)

// 目标类型，供日志与错误共用。
const (
	TargetTurnaround = "FLIGHT_TURNAROUND"
	TargetTask       = "GROUND_TASK"
	TargetResource   = "GROUND_RESOURCE"
	TargetBooking    = "RESOURCE_BOOKING"
	TargetDelay      = "DELAY_EVENT"
)
