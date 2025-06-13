package spot_instrument_service

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/FlyKarlik/proto/spot_instrument_service/gen/spot_instrument_service/proto"
	"github.com/FlyKarlik/spotInstrumentService/config"
	grpc_interceptor "github.com/FlyKarlik/spotInstrumentService/internal/delivery/grpc/interceptor"
	grpc_sync_handler "github.com/FlyKarlik/spotInstrumentService/internal/delivery/grpc/sync"
	"github.com/FlyKarlik/spotInstrumentService/internal/repository"
	"github.com/FlyKarlik/spotInstrumentService/internal/usecase"
	"github.com/FlyKarlik/spotInstrumentService/pkg/logger"
	"github.com/FlyKarlik/spotInstrumentService/pkg/metric"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

type SpotInstrumentService struct {
	cfg        *config.Config
	logger     logger.Logger
	grpcServer *grpc.Server
}

func New(cfg *config.Config, logger logger.Logger) *SpotInstrumentService {
	return &SpotInstrumentService{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *SpotInstrumentService) Start() error {
	const layer = "app"
	const method = "Start"

	s.logger.Info(layer, method, "starting service")

	s.mustSetupTracer()

	repo := s.mustSetupRepo()
	usecase := s.mustSetupUsecase(repo)

	go func() {
		s.logger.Infof(
			layer,
			method,
			"starting prometheus",
			"address: %s", s.cfg.Infrastructure.Prometheus.Address)
		if err := s.mustStartPrometheus(); err != nil {
			s.logger.Error(layer, method, "failed to start prometheus", err)
			os.Exit(1)
		}
		s.logger.Info(layer, method, "prometheus started successfully")
	}()

	go func() {
		s.logger.Infof(
			layer,
			method,
			"starting gRPC server",
			"address: %s", s.cfg.GRPCServer.Address)
		if err := s.mustStartGRPCServer(usecase); err != nil {
			s.logger.Error(layer, method, "failed to start grpc server", err)
			os.Exit(1)
		}
		s.logger.Info(layer, method, "grpc server stopped")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	s.logger.Info(layer, method, "waiting for shutdown signal")
	<-quit
	s.logger.Info(layer, method, "shutdown signal received")

	s.mustStopGRPCServer()
	s.logger.Info(layer, method, "service stopped gracefully")

	return nil
}

func (s *SpotInstrumentService) mustStartPrometheus() error {
	return metric.StartPrometheus(s.cfg)
}

func (s *SpotInstrumentService) mustSetupTracer() {
	const method = "mustSetupTracer"
	const layer = "app"

	s.logger.Info(layer, method, "setting up tracing")
}

func (s *SpotInstrumentService) mustSetupRepo() repository.Repository {
	const method = "mustSetupRepo"
	const layer = "app"

	s.logger.Info(layer, method, "setting up repository")
	return repository.New(s.logger)
}

func (s *SpotInstrumentService) mustSetupUsecase(repo repository.Repository) usecase.Usecase {
	const method = "mustSetuUsecase"
	const layer = "app"

	s.logger.Info(layer, method, "setting up usecase")
	return usecase.New(s.logger, repo)
}

func (s *SpotInstrumentService) mustStartGRPCServer(usecase usecase.Usecase) error {
	const layer = "app"
	const method = "mustStartGRPCServer"

	grpcInterceptor := grpc_interceptor.New(s.logger)
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcInterceptor.XRequestIDInterceptor(),
			grpcInterceptor.LoggerInterceptor(),
			grpcInterceptor.UnaryPanicRecoveryInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
		),
	)
	s.grpcServer = grpcServer

	grpcSyncHandler := grpc_sync_handler.New(s.logger, usecase)
	pb.RegisterSpotInstrumentServiceServer(grpcServer, grpcSyncHandler)

	grpc_prometheus.Register(grpcServer)

	lis, err := net.Listen("tcp", s.cfg.GRPCServer.Address)
	if err != nil {
		s.logger.Error(layer, method, "failed to listen tcp", err, "address", s.cfg.GRPCServer.Address)
		return err
	}

	s.logger.Info(layer, method, "grpc server listening", "address", s.cfg.GRPCServer.Address)
	err = grpcServer.Serve(lis)
	if err != nil {
		s.logger.Error(layer, method, "grpc server serve error", err)
		return err
	}

	return nil
}

func (o *SpotInstrumentService) mustStopGRPCServer() {
	const layer = "app"
	const method = "mustStopGRPCServer"

	o.logger.Info(layer, method, "stopping grpc server")
	o.grpcServer.GracefulStop()
	o.logger.Info(layer, method, "grpc server stopped gracefully")
}
