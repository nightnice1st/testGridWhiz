package grpc

import (
	"context"

	"github.com/nightnice1st/testGridWhiz/internal/auth/usecase"
	pb "github.com/nightnice1st/testGridWhiz/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authUsecase *usecase.AuthUsecase
}

func NewAuthHandler(authUsecase *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user, err := h.authUsecase.Register(req.Email, req.Password, req.Name)
	if err != nil {
		return &pb.RegisterResponse{
			Success: false,
			Message: err.Error(),
		}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.RegisterResponse{
		Success: true,
		Message: "User registered successfully",
		UserId:  user.ID,
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := h.authUsecase.Login(req.Email, req.Password)
	if err != nil {
		return &pb.LoginResponse{
			Success: false,
			Message: err.Error(),
		}, status.Error(codes.Unauthenticated, err.Error())
	}

	return &pb.LoginResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
	}, nil
}

func (h *AuthHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	err := h.authUsecase.Logout(req.Token)
	if err != nil {
		return &pb.LogoutResponse{
			Success: false,
			Message: err.Error(),
		}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.LogoutResponse{
		Success: true,
		Message: "Logout successful",
	}, nil
}
