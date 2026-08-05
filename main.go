package main

import (
	"log"
	"net"

	pb "github.com/LimeOnTop/interverse-contracts/candidate/gen"
	"github.com/LimeOnTop/interverse-candidate/internal/config"
	"github.com/LimeOnTop/interverse-candidate/internal/database"
	"github.com/LimeOnTop/interverse-candidate/internal/handler"
	"github.com/LimeOnTop/interverse-candidate/internal/repository"
	"github.com/LimeOnTop/interverse-candidate/internal/service"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repository
	candidateRepo := repository.NewCandidateRepository(db)

	// Initialize service
	candidateService := service.NewCandidateService(*candidateRepo, cfg.UserServiceURL)

	// Initialize gRPC server
	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCandidateServiceServer(grpcServer, handler.NewCandidateHandler(candidateService))

	log.Printf("Candidate Service starting on port %s", cfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
