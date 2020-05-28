package main

import (
	"context"
	"github.com/gidyon/antibug/internal/pkg/auth"
	"github.com/gidyon/antibug/internal/pkg/errs"
	"github.com/gidyon/antibug/pkg/api/account"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"strings"
)

func authInterceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	// Handle validation
	switch {
	case strings.HasSuffix(info.FullMethod, "ActivateAccount"):
		activateReq, ok := req.(*account.ActivateAccountRequest)
		if !ok {
			return nil, errs.WrapMessage(
				codes.InvalidArgument, "bad activate account request",
			)
		}
		if activateReq.ByAdmin {
			_, err := auth.AuthorizeGroup(ctx, auth.Admin)
			if err != nil {
				return nil, err
			}
		} else {
			err := auth.AuthenticateCtx(ctx)
			if err != nil {
				return nil, err
			}
		}
	case strings.HasSuffix(info.FullMethod, "CreateAccount"):
	case strings.HasSuffix(info.FullMethod, "Login"):
	default:
		err := auth.AuthenticateCtx(ctx)
		if err != nil {
			return nil, err
		}
	}
	return handler(ctx, req)
}
