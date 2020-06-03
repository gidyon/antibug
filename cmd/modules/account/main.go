package main

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/antibug/pkg/api/account"
	app_grpc_middleware "github.com/gidyon/micros/pkg/grpc/middleware"
	"github.com/gidyon/micros/utils/healthcheck"
	"google.golang.org/grpc"

	account_service "github.com/gidyon/antibug/internal/modules/account"

	"github.com/gidyon/micros"
	"github.com/gidyon/micros/pkg/config"

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

	// Recovery middleware
	recoveryUIs, recoverySIs := app_grpc_middleware.AddRecovery()
	unaryInterceptors = append(unaryInterceptors, recoveryUIs...)
	streamInterceptors = append(streamInterceptors, recoverySIs...)

	// Add interceptors to service
	app.AddGRPCStreamServerInterceptors(streamInterceptors...)
	app.AddGRPCUnaryServerInterceptors(unaryInterceptors...)

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

	// Start app
	app.Start(ctx, func() error {
		accountAPI, err := account_service.NewAccountAPI(ctx, &account_service.Options{
			SQLDB:      app.GormDB(),
			Logger:     app.Logger(),
			SigningKey: randomdata.RandStringRunes(32),
		})
		handleErr(err)

		account.RegisterAccountAPIServer(app.GRPCServer(), accountAPI)
		handleErr(account.RegisterAccountAPIHandlerServer(ctx, app.RuntimeMux(), accountAPI))

		return nil
	})
}

func handleErr(err error) {
	if err != nil {
		logrus.Fatalln(err)
	}
}
