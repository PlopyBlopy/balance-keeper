package main

import (
	"github.com/PlopyBlopy/balance-keeper-service/internal/shared/config"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	var log *zap.Logger
	if cfg.IsProd {
		log, err = zap.NewProduction()
		if err != nil {
			panic(err)
		}
	} else {
		log, err = zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
	}
	defer log.Sync()

	err = app(log)
	if err != nil {
		log.Panic("app error", zap.Error(err))
	}
}

func app(log *zap.Logger) error {
	return nil
}
