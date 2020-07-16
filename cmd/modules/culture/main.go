package main

import (
	"context"
	culture_service "github.com/gidyon/antibug/internal/modules/culture"
	"github.com/gidyon/antibug/pkg/api/culture"
	app_grpc_middleware "github.com/gidyon/micros/pkg/grpc/middleware"
	"github.com/gidyon/micros/utils/healthcheck"
	"google.golang.org/grpc"
	"os"

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
	app.AddEndpoint("/api/antibug/cultures/health/ready", healthcheck.RegisterProbe(&healthcheck.ProbeOptions{
		Service:      app,
		Type:         healthcheck.ProbeReadiness,
		AutoMigrator: func() error { return nil },
	}))

	// Liveness health check
	app.AddEndpoint("/api/antibug/cultures/health/live", healthcheck.RegisterProbe(&healthcheck.ProbeOptions{
		Service:      app,
		Type:         healthcheck.ProbeLiveNess,
		AutoMigrator: func() error { return nil },
	}))

	// Start module
	app.Start(ctx, func() error {
		// Create culture tracing instance
		cultureAPI, err := culture_service.NewCultureAPI(ctx, &culture_service.Options{
			SQLDB:      app.GormDB(),
			Logger:     app.Logger(),
			SigningKey: os.Getenv("JWT_SIGNING_KEY"),
		})
		handleErr(err)

		culture.RegisterCultureAPIServer(app.GRPCServer(), cultureAPI)
		handleErr(culture.RegisterCultureAPIHandlerServer(ctx, app.RuntimeMux(), cultureAPI))

		return nil
	})
}

func handleErr(err error) {
	if err != nil {
		logrus.Fatalln(err)
	}
}
