package main

import (
	"context"
	facility_service "github.com/gidyon/antibug/internal/modules/facility"
	"github.com/gidyon/antibug/pkg/api/facility"
	app_grpc_middleware "github.com/gidyon/micros/pkg/grpc/middleware"
	"github.com/gidyon/micros/utils/healthcheck"
	"google.golang.org/grpc"

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
	app.AddEndpoint("/api/antibug/facilities/readyq/", healthcheck.RegisterProbe(&healthcheck.ProbeOptions{
		Service:      app,
		Type:         healthcheck.ProbeReadiness,
		AutoMigrator: func() error { return nil },
	}))

	// Liveness health check
	app.AddEndpoint("/api/antibug/facilities/liveq/", healthcheck.RegisterProbe(&healthcheck.ProbeOptions{
		Service:      app,
		Type:         healthcheck.ProbeLiveNess,
		AutoMigrator: func() error { return nil },
	}))

	// Create facility tracing instance
	facilityAPI, err := facility_service.NewFacilityAPI(ctx, &facility_service.Options{
		SQLDB:  app.GormDB(),
		Logger: app.Logger(),
	})
	handleErr(err)

	facility.RegisterFacilityAPIServer(app.GRPCServer(), facilityAPI)
	handleErr(facility.RegisterFacilityAPIHandlerServer(ctx, app.RuntimeMux(), facilityAPI))

	handleErr(app.Run(ctx))
}

func handleErr(err error) {
	if err != nil {
		logrus.Fatalln(err)
	}
}
