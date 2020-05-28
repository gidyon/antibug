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
	case strings.HasSuffix(info.FullMethod, "CreatePathogen"),
		strings.HasSuffix(info.FullMethod, "UpdatePathogen"),
		strings.HasSuffix(info.FullMethod, "DeletePathogen"):
		_, err := auth.AuthorizeGroup(ctx, auth.Physician, auth.Researcher, auth.LabTechnician, auth.Admin)
		if err != nil {
			return nil, err
		}
	}
	return handler(ctx, req)
}
