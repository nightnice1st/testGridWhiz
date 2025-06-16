package grpc

import (
	"context"
	"strings"

	"github.com/nightnice1st/testGridWhiz/internal/auth/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(authUsecase *usecase.AuthUsecase) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Skip auth for auth service methods
		if strings.Contains(info.FullMethod, "AuthService") {
			return handler(ctx, req)
		}

		// Extract token from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader[0], " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization header format")
		}

		token := tokenParts[1]

		// Validate token
		claims, err := authUsecase.ValidateToken(token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		// Add user info to context
		ctx = context.WithValue(ctx, "userID", claims.UserID)
		ctx = context.WithValue(ctx, "email", claims.Email)

		return handler(ctx, req)
	}
}
