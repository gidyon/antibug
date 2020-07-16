package main

import (
	"context"
	facility_service "github.com/gidyon/antibug/internal/modules/facility"
	"github.com/gidyon/antibug/pkg/api/facility"
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
	app.AddEndpoint("/api/antibug/facilities/health/ready", healthcheck.RegisterProbe(&healthcheck.ProbeOptions{
		Service:      app,
		Type:         healthcheck.ProbeReadiness,
		AutoMigrator: func() error { return nil },
	}))

	// Liveness health check
	app.AddEndpoint("/api/antibug/facilities/health/live", healthcheck.RegisterProbe(&healthcheck.ProbeOptions{
		Service:      app,
		Type:         healthcheck.ProbeLiveNess,
		AutoMigrator: func() error { return nil },
	}))

	// Start service
	app.Start(ctx, func() error {
		// Create facility tracing instance
		facilityAPI, err := facility_service.NewFacilityAPI(ctx, &facility_service.Options{
			SQLDB:         app.GormDB(),
			Logger:        app.Logger(),
			JWTSigningKey: os.Getenv("JWT_SIGNING_KEY"),
		})
		handleErr(err)

		facility.RegisterFacilityAPIServer(app.GRPCServer(), facilityAPI)
		handleErr(facility.RegisterFacilityAPIHandlerServer(ctx, app.RuntimeMux(), facilityAPI))

		return nil
	})
}

func handleErr(err error) {
	if err != nil {
		logrus.Fatalln(err)
	}
}
