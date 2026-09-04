package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"ride-sharing/shared/types"
)

type previewTripRequest struct {
	UserID      string           `json:"userID"`
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	var req previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, "Missing userID", http.StatusBadRequest)
		return
	}

	jsonBody, _ := json.Marshal(req)
	reader := bytes.NewReader(jsonBody)

	// call trip

	res, err := http.Post("http://trip-service:8083/preview", "application/json", reader)
	if err != nil {
		http.Error(w, "Failed to create request to trip service", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	var resBody any
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		http.Error(w, "Failed to decode response from trip service", http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, res.StatusCode, resBody)
}

func writeJsonResponse(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}
