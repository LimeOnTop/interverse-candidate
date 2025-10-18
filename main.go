package main

import (
	"log"
	"net"

	pb "github.com/inter-verse/services/candidate-service"
	"github.com/inter-verse/services/candidate-service/internal/config"
	"github.com/inter-verse/services/candidate-service/internal/database"
	"github.com/inter-verse/services/candidate-service/internal/handler"
	"github.com/inter-verse/services/candidate-service/internal/repository"
	"github.com/inter-verse/services/candidate-service/internal/service"
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
