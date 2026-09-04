package main

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"time"
)

func main() {
	ctx := context.Background()
	inmemRepo := repository.NewInmemRepository()
	svc := service.NewService(inmemRepo)
	fare := &domain.RideFareModel{
		UserID: "user123",
	}

	t, err := svc.CreateTrip(ctx, fare)
	if err != nil {
		panic(err)
	}
	println("Trip created with ID:", t.ID.Hex())

	// for {
	// 	time.Sleep(time.Second)
	// }
}
