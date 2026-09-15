package grpc

import (
	"context"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pb.UnimplementedTripServiceServer
	service domain.TripService
}

func NewGRPCHanlder(server *grpc.Server, service domain.TripService) *gRPCHandler {
	handler := &gRPCHandler{
		service: service,
	}
	pb.RegisterTripServiceServer(server, handler)
	return handler
}

func (h *gRPCHandler) PreviewTrip(ctx context.Context, req *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {
	pickup := req.GetStartLocation()
	destination := req.GetEndLocation()

	pickupCoord := &types.Coordinate{Latitude: pickup.Latitude, Longitude: pickup.Longitude}
	destinationCoord := &types.Coordinate{Latitude: destination.Latitude, Longitude: destination.Longitude}

	route, err := h.service.GetRoute(ctx, pickupCoord, destinationCoord)
	if err != nil {
		log.Println("Error creating trip:", err)
		return nil, status.Errorf(codes.Internal, "failed to get route :%v", err)
	}

	userId := req.GetUserID()
	// 1.Estimate the ride fare price based on the route
	estimate := h.service.EstimatePackagesPriceWithRoute(route)

	// 2.Store the ride fares for the create trip
	fares, err := h.service.GenerateTripFares(ctx, estimate, userId, route)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate the ride fares :%v", err)
	}
	return &pb.PreviewTripResponse{
		Route:     route.ToProto(),
		RideFares: domain.ToRideFareProto(fares),
	}, nil
}

func (h *gRPCHandler) CreateTrip(ctx context.Context, req *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {
	fareID := req.GetRideFareID()
	userID := req.UserID
	rideFare, err := h.service.GetAndValidateFare(ctx, fareID, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to validate the fare: %v", err)
	}

	trip, err := h.service.CreateTrip(ctx, rideFare)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create the trip: %v", err)
	}

	//4. Add a comment at the end of function to publish an event on the Async Comms Module

	return &pb.CreateTripResponse{
		TripID: trip.ID.Hex(),
	}, nil
}
