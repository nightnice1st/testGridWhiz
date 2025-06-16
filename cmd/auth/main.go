package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	authDelivery "github.com/nightnice1st/testGridWhiz/internal/auth/delivery"
	authRepo "github.com/nightnice1st/testGridWhiz/internal/auth/repository"
	authUsecase "github.com/nightnice1st/testGridWhiz/internal/auth/usecase"
	"github.com/nightnice1st/testGridWhiz/internal/pkg/config"
	"github.com/nightnice1st/testGridWhiz/internal/pkg/database"
	"github.com/nightnice1st/testGridWhiz/internal/pkg/ratelimit"
	userRepo "github.com/nightnice1st/testGridWhiz/internal/users/repository"
	pb "github.com/nightnice1st/testGridWhiz/pb"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to MongoDB
	mongoClient, err := database.Connect(cfg.MongoDBURI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer mongoClient.Disconnect(nil)

	db := mongoClient.Database("mydb")

	// Initialize repositories
	userRepository := userRepo.NewUserRepository(db)
	authRepository := authRepo.NewAuthRepository(db)

	// Initialize rate limiter
	rateLimiter := ratelimit.NewRateLimiter(cfg.RateLimitAttempts, cfg.RateLimitWindow)

	// Initialize use cases
	authUseCase := authUsecase.NewAuthUsecase(
		userRepository,
		authRepository,
		cfg.JWTSecret,
		cfg.JWTExpiry,
		rateLimiter,
	)

	// Initialize gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.AuthServicePort))
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	grpcServer := grpc.NewServer()

	// Register service
	authHandler := authDelivery.NewAuthHandler(authUseCase)
	pb.RegisterAuthServiceServer(grpcServer, authHandler)

	// Register reflection service for development
	reflection.Register(grpcServer)

	log.Printf("Auth service starting on port %s", cfg.AuthServicePort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to serve:", err)
	}
}
