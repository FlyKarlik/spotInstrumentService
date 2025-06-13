package usecase

import (
	"context"

	"github.com/FlyKarlik/spotInstrumentService/internal/domain"
	"github.com/FlyKarlik/spotInstrumentService/internal/repository"
	"github.com/FlyKarlik/spotInstrumentService/pkg/logger"
)

type IMarketUsecase interface {
	ViewMarkets(ctx context.Context, req domain.ViewMarketsRequest) (domain.ViewMarketsResponse, error)
}

type Usecase interface {
	IMarketUsecase
}

type usecaseImpl struct {
	IMarketUsecase
}

func New(l logger.Logger, repo repository.Repository) *usecaseImpl {
	return &usecaseImpl{
		IMarketUsecase: newMarketUsecase(l, repo),
	}
}
