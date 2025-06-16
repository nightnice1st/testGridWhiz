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
	userDelivery "github.com/nightnice1st/testGridWhiz/internal/users/delivery"
	userRepo "github.com/nightnice1st/testGridWhiz/internal/users/repository"
	userUsecase "github.com/nightnice1st/testGridWhiz/internal/users/usecase"
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

	// Initialize rate limiter for auth
	rateLimiter := ratelimit.NewRateLimiter(cfg.RateLimitAttempts, cfg.RateLimitWindow)

	// Initialize use cases
	userUseCase := userUsecase.NewUserUsecase(userRepository)
	authUseCase := authUsecase.NewAuthUsecase(
		userRepository,
		authRepository,
		cfg.JWTSecret,
		cfg.JWTExpiry,
		rateLimiter,
	)

	// Initialize gRPC server with auth interceptor
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.UserServicePort))
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authDelivery.AuthInterceptor(authUseCase)),
	)

	// Register service
	userHandler := userDelivery.NewUserHandler(userUseCase)
	pb.RegisterUserServiceServer(grpcServer, userHandler)

	// Register reflection service for development
	reflection.Register(grpcServer)

	log.Printf("User service starting on port %s", cfg.UserServicePort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to serve:", err)
	}
}
