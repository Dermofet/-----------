package middleware

import (
	"context"
	"go-test-grpc-http/internal/entity"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey = string

const userIDKey contextKey = "user-id"

func NewAuthMiddleware() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Aborted, "metadata not found")
		}

		authorization := md.Get("authorization")
		if len(authorization) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "can't find authorization")
		}

		tokenString := strings.TrimPrefix(authorization[0], "Bearer ")
		if len(tokenString) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token format")
		}

		id, err := entity.ParseToken(tokenString)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		c := context.WithValue(ctx, userIDKey, id.Id)

		return handler(c, req)
	}
}
