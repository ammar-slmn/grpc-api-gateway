package main

import (
	"context"
	"log"
	"net"

	userpb "grp-api-gateway/proto/user" // Import the generated gRPC code from proto/user

	"google.golang.org/grpc"
)

// userServiceServer implements the UserService gRPC server
type userServiceServer struct {
	userpb.UnimplementedUserServiceServer
}

// GetUser handles incoming GetUser requests from clients
func (s *userServiceServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	// Return Mock response for the GetUser request
	return &userpb.GetUserResponse{
		Id:    req.Id,
		Name:  "Ammar Salman",
		Email: "ammar@email.com",
	}, nil
}

func main() {
	// Start listening for incoming TCP connections on port 50051
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	// Create a new gRPC server instance
	grpcServer := grpc.NewServer()

	// Register our user service with the gRPC server
	userpb.RegisterUserServiceServer(grpcServer, &userServiceServer{})

	log.Println("UserService is running on port :50051")

	// Start serving requests
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
