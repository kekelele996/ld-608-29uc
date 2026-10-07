package constants

// 错误消息模板集中定义，格式化参数由各 service 填充。
const (
	AuthRequiredMessage      = "missing or malformed token"
	AuthInvalidMessage       = "invalid username or password"
	ForbiddenMessage         = "role %s is not allowed to %s"
	RateLimitedMessage       = "too many requests, slow down"
	ValidationFailedMessage  = "validation failed: %s"
	RecordNotFoundMessage    = "%s#%v not found"
	InvalidEnumMessage       = "invalid %s value: %s"
	ConflictMessage          = "resource %s is already booked in the requested window"
	DeadlineNotExtendedMsg   = "delay minutes must be positive to extend deadlines"
	TaskNotAcceptableMessage = "task %d is %s, only PENDING/IN_PROGRESS/BLOCKED tasks can be accepted"
	InternalErrorMessage     = "internal server error"
)
