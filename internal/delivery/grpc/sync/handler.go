package grpc_sync_handler

import (
	pb "github.com/FlyKarlik/proto/spot_instrument_service/gen/spot_instrument_service/proto"
	"github.com/FlyKarlik/spotInstrumentService/internal/usecase"
	"github.com/FlyKarlik/spotInstrumentService/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type GRPCSyncHandler struct {
	logger  logger.Logger
	usecase usecase.Usecase
	tracer  trace.Tracer
	pb.UnimplementedSpotInstrumentServiceServer
}

func New(logger logger.Logger, usecase usecase.Usecase) *GRPCSyncHandler {
	return &GRPCSyncHandler{
		logger:  logger,
		usecase: usecase,
		tracer:  otel.Tracer("spot-instrument-service/grpc-sync-handler"),
	}
}
