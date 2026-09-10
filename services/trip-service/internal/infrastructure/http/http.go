package http

import (
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/types"
)

type previewTripRequest struct {
	UserID      string           `json:"userID"`
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

type TripServiceHandler struct {
	Service domain.TripService
}

func (s *TripServiceHandler) HandleTripPreview(w http.ResponseWriter, r *http.Request) {
	var req previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// fare := &domain.RideFareModel{
	// 	UserID: "user123s",
	// }

	// t, err := s.Service.CreateTrip(r.Context(), fare)
	t, err := s.Service.GetRoute(r.Context(), &req.Pickup, &req.Destination)
	if err != nil {
		log.Println("Error creating trip:", err)
	}
	writeJsonResponse(w, http.StatusOK, t)
}

func writeJsonResponse(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}
