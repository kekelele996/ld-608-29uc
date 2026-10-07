package models

type DelayEvent struct {
	ID                 int    `json:"id"`
	TurnaroundID       int    `json:"turnaround_id"`
	DelayType          string `json:"delay_type"`
	Minutes            int    `json:"minutes"`
	RootCause          string `json:"root_cause"`
	ResponsibilityTeam string `json:"responsibility_team"`
	ResolvedAt         string `json:"resolved_at"`
}
