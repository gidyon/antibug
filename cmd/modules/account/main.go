package main

import (
	"context"
	"github.com/gidyon/antibug/pkg/api/account"
	app_grpc_middleware "github.com/gidyon/micros/pkg/grpc/middleware"
	"github.com/gidyon/micros/utils/healthcheck"
	"google.golang.org/grpc"

	account_service "github.com/gidyon/antibug/internal/modules/account"

	"github.com/gidyon/config"
	"github.com/gidyon/micros"

	"github.com/Sirupsen/logrus"
)

func main() {
	cfg, err := config.New()
	handleErr(err)

	ctx := context.Background()

	app, err := micros.NewService(ctx, cfg, nil)
	handleErr(err)

	unaryInterceptors := make([]grpc.UnaryServerInterceptor, 0)
	streamInterceptors := make([]grpc.StreamServerInterceptor, 0)

	// Unary interceptor for authentication
	contractAuth := grpc.UnaryServerInterceptor(authInterceptor)
	unaryInterceptors = append(unaryInterceptors, contractAuth)

	// Recovery middleware
	recoveryUIs, recoverySIs := app_grpc_middleware.AddRecovery()
	unaryInterceptors = append(unaryInterceptors, recoveryUIs...)
	streamInterceptors = append(streamInterceptors, recoverySIs...)

	// Initialize grpc server
	handleErr(app.InitGRPC(ctx))

	// Readiness health check
	app.AddEndpoint("/api/antibug/accounts/readyq/", healthcheck.RegisterProbe(&healthcheck.ProbeOptions{
		Service:      app,
		Type:         healthcheck.ProbeReadiness,
		AutoMigrator: func() error { return nil },
	}))

	// Liveness health check
	app.AddEndpoint("/api/antibug/accounts/liveq/", healthcheck.RegisterProbe(&healthcheck.ProbeOptions{
		Service:      app,
		Type:         healthcheck.ProbeLiveNess,
		AutoMigrator: func() error { return nil },
	}))

	// Create account tracing instance
	accountAPI, err := account_service.NewAccountAPI(ctx, &account_service.Options{
		SQLDB:  app.GormDB(),
		Logger: app.Logger(),
	})
	handleErr(err)

	account.RegisterAccountAPIServer(app.GRPCServer(), accountAPI)
	handleErr(account.RegisterAccountAPIHandlerServer(ctx, app.RuntimeMux(), accountAPI))

	handleErr(app.Run(ctx))
}

func handleErr(err error) {
	if err != nil {
		logrus.Fatalln(err)
	}
}
