package types

import "time"

type GroundResourceView struct {
	ID                 uint       `json:"id"`
	ResourceCode       string     `json:"resource_code"`
	ResourceType       string     `json:"resource_type"`
	Location           string     `json:"location"`
	AvailabilityStatus string     `json:"availability_status"`
	StatusText         string     `json:"status_text"`
	MaintenanceDueAt   *time.Time `json:"maintenance_due_at"`
	OwnerTeam          string     `json:"owner_team"`
}

type GroundResourceUpsertRequest struct {
	ResourceCode       string     `json:"resource_code" binding:"required"`
	ResourceType       string     `json:"resource_type" binding:"required"`
	Location           string     `json:"location"`
	AvailabilityStatus string     `json:"availability_status"`
	MaintenanceDueAt   *time.Time `json:"maintenance_due_at"`
	OwnerTeam          string     `json:"owner_team"`
}
