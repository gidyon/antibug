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
	case strings.HasSuffix(info.FullMethod, "CreateAntimicrobial"),
		strings.HasSuffix(info.FullMethod, "UpdateAntimicrobial"):
		_, err := auth.AuthorizeGroup(ctx, auth.Physician, auth.Researcher, auth.LabTechnician, auth.Admin, auth.Pharmacist)
		if err != nil {
			return nil, err
		}
	case strings.HasSuffix(info.FullMethod, "DeleteAntimicrobial"):
		_, err := auth.AuthorizeGroup(ctx, auth.Physician, auth.Researcher, auth.Admin)
		if err != nil {
			return nil, err
		}
	}
	return handler(ctx, req)
}
