package main

import (
	"context"
	antibiogram "github.com/gidyon/antibug/pkg/api/antibiogram"
	app_grpc_middleware "github.com/gidyon/micros/pkg/grpc/middleware"
	"github.com/gidyon/micros/utils/healthcheck"
	"google.golang.org/grpc"
	"os"

	antibiogram_service "github.com/gidyon/antibug/internal/modules/antibiogram"

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
	app.AddEndpoint("/api/antibug/antibiograms/health/ready", healthcheck.RegisterProbe(&healthcheck.ProbeOptions{
		Service:      app,
		Type:         healthcheck.ProbeReadiness,
		AutoMigrator: func() error { return nil },
	}))

	// Liveness health check
	app.AddEndpoint("/api/antibug/antibiograms/health/live", healthcheck.RegisterProbe(&healthcheck.ProbeOptions{
		Service:      app,
		Type:         healthcheck.ProbeLiveNess,
		AutoMigrator: func() error { return nil },
	}))

	// Start app
	app.Start(ctx, func() error {
		antibiogramAPI, err := antibiogram_service.NewAntibiogramAPIServer(ctx, &antibiogram_service.Options{
			SQLDB:         app.GormDB(),
			RedisDB:       app.RedisClient(),
			Logger:        app.Logger(),
			JWTSigningKey: os.Getenv("JWT_SIGNING_KEY"),
		})
		handleErr(err)

		antibiogram.RegisterAntibiogramAPIServer(app.GRPCServer(), antibiogramAPI)
		handleErr(antibiogram.RegisterAntibiogramAPIHandlerServer(ctx, app.RuntimeMux(), antibiogramAPI))

		return nil
	})
}

func handleErr(err error) {
	if err != nil {
		logrus.Fatalln(err)
	}
}
