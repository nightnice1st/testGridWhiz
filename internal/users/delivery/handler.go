package grpc

import (
	"context"

	"github.com/nightnice1st/testGridWhiz/internal/users/domain"
	pb "github.com/nightnice1st/testGridWhiz/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userUsecase domain.UserUsecase
}

func NewUserHandler(userUsecase domain.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	users, total, err := h.userUsecase.ListUsers(
		int(req.Page),
		int(req.Limit),
		req.NameFilter,
		req.EmailFilter,
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbUsers := make([]*pb.User, len(users))
	for i, user := range users {
		pbUsers[i] = &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return &pb.ListUsersResponse{
		Users: pbUsers,
		Total: int32(total),
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

func (h *UserHandler) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	// Get userID from context (set by auth middleware) or use requested ID
	userID := req.UserId
	if userID == "" {
		if ctxUserID, ok := ctx.Value("userID").(string); ok {
			userID = ctxUserID
		} else {
			return nil, status.Error(codes.InvalidArgument, "user ID required")
		}
	}

	user, err := h.userUsecase.GetProfile(userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.GetProfileResponse{
		User: &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		},
	}, nil
}

func (h *UserHandler) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	// Get userID from context - users can only update their own profile
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	// Ensure user is updating their own profile
	if req.UserId != "" && req.UserId != userID {
		return nil, status.Error(codes.PermissionDenied, "can only update own profile")
	}

	user, err := h.userUsecase.UpdateProfile(userID, req.Name)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.UpdateProfileResponse{
		Success: true,
		Message: "Profile updated successfully",
		User: &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		},
	}, nil
}

func (h *UserHandler) DeleteProfile(ctx context.Context, req *pb.DeleteProfileRequest) (*pb.DeleteProfileResponse, error) {
	// Get userID from context - users can only delete their own profile
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	// Ensure user is deleting their own profile
	if req.UserId != "" && req.UserId != userID {
		return nil, status.Error(codes.PermissionDenied, "can only delete own profile")
	}

	err := h.userUsecase.DeleteProfile(userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteProfileResponse{
		Success: true,
		Message: "Profile deleted successfully",
	}, nil
}
