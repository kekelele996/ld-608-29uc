package services

import (
	"groundTurn/src/models"
	"groundTurn/src/repositories"
)

func ListFlightTurnarounds() []models.FlightTurnaround {
	return repositories.ListFlightTurnarounds()
}
