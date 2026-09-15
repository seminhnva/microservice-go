package main

import (
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	var req previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.UserID == "" {
		http.Error(w, "Missing userID", http.StatusBadRequest)
		return
	}

	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer tripService.Close()
	// tripService.Client.PreviewTrip()

	// call trip

	tripPreview, err := tripService.Client.PreviewTrip(r.Context(), req.ToProto())
	if err != nil {
		log.Printf("Failed to preview a trip:%v", err)
		http.Error(w, "Failed to preview trip", http.StatusInternalServerError)
		return
	}
	response := contracts.APIResponse{Data: tripPreview}
	writeJsonResponse(w, http.StatusOK, response)
}

func handleTripStart(w http.ResponseWriter, r *http.Request) {
	var req startTripRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if req.UserID == "" {
		http.Error(w, "Missing userID", http.StatusBadRequest)
		return
	}
	if req.RideFareID == "" {
		http.Error(w, "Missing RideFareID", http.StatusBadRequest)
		return
	}
	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer tripService.Close()
	tripStart, err := tripService.Client.CreateTrip(r.Context(), req.ToProto())
	if err != nil {
		log.Printf("Failed to create a trip :%v", err)
		http.Error(w, "Failed to preview trip", http.StatusInternalServerError)
		return
	}
	response := contracts.APIResponse{Data: tripStart}
	writeJsonResponse(w, http.StatusOK, response)

}

func writeJsonResponse(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}
