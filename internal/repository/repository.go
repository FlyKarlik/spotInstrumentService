package repository

import (
	"context"

	"github.com/FlyKarlik/spotInstrumentService/internal/domain"
	in_memory_repo "github.com/FlyKarlik/spotInstrumentService/internal/repository/in_memory"
	"github.com/FlyKarlik/spotInstrumentService/pkg/logger"
)

type IMarketRepository interface {
	GetMarkets(ctx context.Context) ([]domain.Market, error)
}

type Repository interface {
	IMarketRepository
}

type repositoryImpl struct {
	IMarketRepository
}

func New(l logger.Logger) *repositoryImpl {
	return &repositoryImpl{
		IMarketRepository: in_memory_repo.NewInMemoryMarketRepo(l),
	}
}
