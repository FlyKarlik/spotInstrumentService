package grpc_interceptor

import "github.com/FlyKarlik/spotInstrumentService/pkg/logger"

type GRPCInterceptor struct {
	logger logger.Logger
}

func New(logger logger.Logger) *GRPCInterceptor {
	return &GRPCInterceptor{
		logger: logger,
	}
}
