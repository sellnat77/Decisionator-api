package models

type Coordinate struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type MeetCoordinates struct {
	Coordinates []Coordinate `json:"coords"`
}
