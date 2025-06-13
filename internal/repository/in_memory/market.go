package in_memory_repo

import (
	"context"
	"sync"
	"time"

	"github.com/FlyKarlik/spotInstrumentService/internal/domain"
	shared_context "github.com/FlyKarlik/spotInstrumentService/pkg/context"
	"github.com/FlyKarlik/spotInstrumentService/pkg/generics"
	"github.com/FlyKarlik/spotInstrumentService/pkg/logger"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type marketInMemoryRepo struct {
	mu     sync.RWMutex
	data   map[uuid.UUID]domain.Market
	logger logger.Logger
	tracer trace.Tracer
}

func NewInMemoryMarketRepo(l logger.Logger) *marketInMemoryRepo {
	marketRepo := &marketInMemoryRepo{
		data:   make(map[uuid.UUID]domain.Market),
		logger: l,
		tracer: otel.Tracer("spot-instrument-service/repo"),
	}

	now := time.Now()
	markets := []domain.Market{
		{
			ID:           generics.Pointer(uuid.New()),
			Name:         generics.Pointer("BTC-USDT"),
			Enabled:      generics.Pointer(true),
			DeletedAt:    nil,
			AllowedRoles: domain.UserRolesEnum{domain.UserRoleEnumTrader, domain.UserRoleEnumAdmin},
		},
		{
			ID:           generics.Pointer(uuid.New()),
			Name:         generics.Pointer("DOGE-USDT"),
			Enabled:      generics.Pointer(true),
			DeletedAt:    nil,
			AllowedRoles: domain.UserRolesEnum{domain.UserRoleEnumViewer},
		},
		{
			ID:           generics.Pointer(uuid.New()),
			Name:         generics.Pointer("ETH-USDT"),
			Enabled:      generics.Pointer(false),
			DeletedAt:    nil,
			AllowedRoles: domain.UserRolesEnum{domain.UserRoleEnumTrader},
		},
		{
			ID:           generics.Pointer(uuid.New()),
			Name:         generics.Pointer("SOL-USDT"),
			Enabled:      generics.Pointer(true),
			DeletedAt:    &now,
			AllowedRoles: domain.UserRolesEnum{domain.UserRoleEnumAdmin},
		},
	}

	for _, market := range markets {
		marketRepo.data[*market.ID] = market
	}

	return marketRepo
}

func (m *marketInMemoryRepo) GetMarkets(ctx context.Context) ([]domain.Market, error) {
	const layer = "repo"
	const method = "GetMarkets"

	ctx, span := m.tracer.Start(ctx, "MarketInMemoryRepo.GetMarkets")
	defer span.End()

	xRequestID := shared_context.XRequestIDFromContext(ctx)
	span.SetAttributes(attribute.String("x-request-id", xRequestID))

	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []domain.Market
	for _, market := range m.data {
		result = append(result, market)
	}

	span.SetAttributes(attribute.Int("markets.total_count", len(result)))

	m.logger.Info(layer, method, "retrieved markets",
		"x_request_id", xRequestID,
		"markets_total", len(result),
	)

	return result, nil
}
