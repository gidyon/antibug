package main

import (
	"context"
	"github.com/gidyon/antibug/internal/pkg/auth"
	"google.golang.org/grpc"
)

func authInterceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	_, err := auth.AuthorizeGroup(ctx, auth.Physician, auth.Researcher, auth.LabTechnician, auth.Admin)
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}
