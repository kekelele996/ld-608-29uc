package constructors

import (
	"time"

	"groundTurn/src/models"
	"groundTurn/src/types"
)

func BuildAuditLogView(l models.AuditLog) types.AuditLogView {
	return types.AuditLogView{
		ID:         l.ID,
		Actor:      l.Actor,
		Action:     l.Action,
		TargetType: l.TargetType,
		TargetID:   l.TargetID,
		Detail:     l.Detail,
		CreatedAt:  l.CreatedAt.UTC().Format(time.RFC3339),
	}
}
