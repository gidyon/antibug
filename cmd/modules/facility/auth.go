package main

import (
	"context"
	"github.com/gidyon/antibug/internal/pkg/auth"
	"google.golang.org/grpc"
	"strings"
)

func authInterceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	switch {
	case strings.HasSuffix(info.FullMethod, "AddHospital"),
		strings.HasSuffix(info.FullMethod, "RemoveHospital"):
		_, err := auth.AuthorizeGroup(ctx, auth.Admin)
		if err != nil {
			return nil, err
		}
	}
	return handler(ctx, req)
}
