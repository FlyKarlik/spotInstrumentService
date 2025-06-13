package usecase

import (
	"context"

	"github.com/FlyKarlik/spotInstrumentService/internal/domain"
	"github.com/FlyKarlik/spotInstrumentService/internal/errs"
	"github.com/FlyKarlik/spotInstrumentService/internal/repository"
	shared_context "github.com/FlyKarlik/spotInstrumentService/pkg/context"
	"github.com/FlyKarlik/spotInstrumentService/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type marketUsecase struct {
	logger logger.Logger
	repo   repository.Repository
	tracer trace.Tracer
}

func newMarketUsecase(logger logger.Logger, repo repository.Repository) *marketUsecase {
	return &marketUsecase{
		logger: logger,
		repo:   repo,
		tracer: otel.Tracer("spot-instrument-service/usecase"),
	}
}

func (m *marketUsecase) ViewMarkets(
	ctx context.Context,
	req domain.ViewMarketsRequest,
) (domain.ViewMarketsResponse, error) {
	const layer = "usecase"
	const method = "ViewMarkets"

	ctx, span := m.tracer.Start(ctx, "MarketUsecase.ViewMarkets")
	defer span.End()

	xRequestID := shared_context.XRequestIDFromContext(ctx)
	span.SetAttributes(
		attribute.String("x-request-id", xRequestID),
		attribute.Int("user_roles_count", len(req.UserRoles)),
	)

	m.logger.Info(layer, method, "view markets request received",
		"x_request_id", xRequestID,
		"user_roles", req.UserRoles,
	)

	allMarkets, err := m.repo.GetMarkets(ctx)
	if err != nil {
		m.logger.Error(layer, method, "failed to get markets",
			err,
			"x-request-id", xRequestID,
		)
		return domain.ViewMarketsResponse{}, errs.ErrUnknown
	}

	var filtered []domain.Market
	roleSet := make(map[domain.UserRoleEnum]struct{}, len(req.UserRoles))
	for _, role := range req.UserRoles {
		roleSet[role] = struct{}{}
	}

	for _, market := range allMarkets {
		if !*market.Enabled || market.DeletedAt != nil {
			continue
		}

		for _, allowedRole := range market.AllowedRoles {
			if _, ok := roleSet[allowedRole]; ok {
				filtered = append(filtered, market)
				break
			}
		}
	}

	span.SetAttributes(
		attribute.Int("markets.returned_count", len(filtered)),
	)
	m.logger.Info(layer, method, "markets filtered and returned",
		"x_request_id", xRequestID,
		"markets_count", len(filtered),
	)

	return domain.ViewMarketsResponse{
		Markets: filtered,
	}, nil
}
