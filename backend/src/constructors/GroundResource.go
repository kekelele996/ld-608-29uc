package constructors

import (
	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/types"
)

func NewGroundResource(req types.GroundResourceUpsertRequest) models.GroundResource {
	status := req.AvailabilityStatus
	if status == "" {
		status = constants.ResourceAvailable
	}
	return models.GroundResource{
		ResourceCode:       req.ResourceCode,
		ResourceType:       req.ResourceType,
		Location:           req.Location,
		AvailabilityStatus: status,
		MaintenanceDueAt:   req.MaintenanceDueAt,
		OwnerTeam:          req.OwnerTeam,
	}
}

func BuildGroundResourceView(r models.GroundResource) types.GroundResourceView {
	return types.GroundResourceView{
		ID:                 r.ID,
		ResourceCode:       r.ResourceCode,
		ResourceType:       r.ResourceType,
		Location:           r.Location,
		AvailabilityStatus: r.AvailabilityStatus,
		StatusText:         constants.StatusText["ResourceStatus"][r.AvailabilityStatus],
		MaintenanceDueAt:   r.MaintenanceDueAt,
		OwnerTeam:          r.OwnerTeam,
	}
}
