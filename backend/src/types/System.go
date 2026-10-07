package types

// LoginRequest 登录获取 JWT。
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginView 登录响应。
type LoginView struct {
	Token       string `json:"token"`
	Username    string `json:"username"`
	Role        string `json:"role"`
	RoleText    string `json:"role_text"`
	TeamCode    string `json:"team_code"`
	DisplayName string `json:"display_name"`
}

// AuditLogView 操作日志。
type AuditLogView struct {
	ID         uint   `json:"id"`
	Actor      string `json:"actor"`
	Action     string `json:"action"`
	TargetType string `json:"target_type"`
	TargetID   string `json:"target_id"`
	Detail     string `json:"detail"`
	CreatedAt  string `json:"created_at"`
}
