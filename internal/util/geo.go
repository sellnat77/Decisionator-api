package util

import (
	"github.com/sellnat77/Decisionator-api/internal/models"
)

func CoordLatSum(coords *models.MeetCoordinates) float64 {
	var result float64 = 0
	for _, test := range coords.Coordinates {
		result += test.Lat
	}
	return result
}
func CoordLonSum(coords *models.MeetCoordinates) float64 {
	var result float64 = 0
	for _, test := range coords.Coordinates {
		result += test.Lon
	}
	return result
}
