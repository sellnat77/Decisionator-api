package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sellnat77/Decisionator-api/internal/models"
	"github.com/sellnat77/Decisionator-api/internal/util"
)

func Meet(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "Hello from sign in")
}

func GuestMeet(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var inputCoords models.MeetCoordinates

	err := json.NewDecoder(req.Body).Decode(&inputCoords)

	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	numCoords := len(inputCoords.Coordinates)

	midLat := util.CoordLatSum(&inputCoords) / float64(numCoords)
	midLon := util.CoordLonSum(&inputCoords) / float64(numCoords)

	fmt.Println(midLat, midLon)
	midpoint := models.Coordinate{midLat, midLon}
	payload := json.NewEncoder(w).Encode(midpoint)

	fmt.Fprintf(w, payload)
}
