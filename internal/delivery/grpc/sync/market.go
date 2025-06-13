package grpc_sync_handler

import (
	"context"

	pb "github.com/FlyKarlik/proto/spot_instrument_service/gen/spot_instrument_service/proto"
	"github.com/FlyKarlik/spotInstrumentService/internal/mapper"
	shared_context "github.com/FlyKarlik/spotInstrumentService/pkg/context"
	"github.com/FlyKarlik/spotInstrumentService/pkg/validate"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g *GRPCSyncHandler) ViewMarkets(
	ctx context.Context,
	req *pb.ViewMarketsRequest,
) (*pb.ViewMarketsResponse, error) {
	const layer = "delivery"
	const method = "ViewMarkets"

	ctx, span := g.tracer.Start(ctx, "GRPCSyncHandler.ViewMarkets")
	defer span.End()

	xRequestID := shared_context.XRequestIDFromContext(ctx)
	span.SetAttributes(
		attribute.String("x-request-id", xRequestID),
		attribute.Int("user_roles_count", len(req.GetUserRoles())),
	)

	g.logger.Info(layer, method, "view markets gRPC request received",
		"x_request_id", xRequestID,
		"user_roles", req.GetUserRoles(),
	)

	domainReq := mapper.FromProtoViewMarketsRequest(req)

	if err := validate.Validate(domainReq); err != nil {
		g.logger.Error(layer, method, "invalid view markets request", err)
		span.RecordError(err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := g.usecase.ViewMarkets(ctx, mapper.FromProtoViewMarketsRequest(req))
	if err != nil {
		g.logger.Error(layer, method, "failed to process view markets request",
			err,
			"x_request_id", xRequestID,
		)
		span.RecordError(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoResp := mapper.ToProtoViewMarketResponse(resp)

	span.SetAttributes(
		attribute.Int("markets.returned_count", len(protoResp.GetMarkets())),
	)
	g.logger.Info(layer, method, "markets successfully processed and returned",
		"x_request_id", xRequestID,
		"markets_count", len(protoResp.GetMarkets()),
	)

	return protoResp, nil
}
