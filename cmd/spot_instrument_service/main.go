package main

import (
	"github.com/FlyKarlik/spotInstrumentService/config"
	"github.com/FlyKarlik/spotInstrumentService/internal/app/spot_instrument_service"
	"github.com/FlyKarlik/spotInstrumentService/pkg/logger"
	"github.com/FlyKarlik/spotInstrumentService/pkg/validate"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	if err := validate.Validate(cfg); err != nil {
		panic(err)
	}

	logger, err := logger.New(cfg.SpotInstrumentService.LogLevel)
	if err != nil {
		panic(err)
	}

	spotInstrumentSerivce := spot_instrument_service.New(cfg, logger)
	if err := spotInstrumentSerivce.Start(); err != nil {
		panic(err)
	}
}
