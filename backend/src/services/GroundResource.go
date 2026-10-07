package services

import (
	"slices"

	"groundTurn/src/constants"
	"groundTurn/src/constructors"
	"groundTurn/src/repositories"
	"groundTurn/src/types"
	"groundTurn/src/utils"
)

type GroundResourceService struct {
	repo *repositories.GroundResourceRepository
}

func NewGroundResourceService(repo *repositories.GroundResourceRepository) *GroundResourceService {
	return &GroundResourceService{repo: repo}
}

func (s *GroundResourceService) List() ([]types.GroundResourceView, error) {
	rows, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	views := make([]types.GroundResourceView, 0, len(rows))
	for i := range rows {
		views = append(views, constructors.BuildGroundResourceView(rows[i]))
	}
	return views, nil
}

func (s *GroundResourceService) Upsert(req types.GroundResourceUpsertRequest) (*types.GroundResourceView, error) {
	if !slices.Contains(constants.ResourceStatus, req.AvailabilityStatus) && req.AvailabilityStatus != "" {
		return nil, utils.EnumError("availability_status", req.AvailabilityStatus)
	}
	row := constructors.NewGroundResource(req)
	if err := s.repo.Create(&row); err != nil {
		return nil, err
	}
	view := constructors.BuildGroundResourceView(row)
	return &view, nil
}
