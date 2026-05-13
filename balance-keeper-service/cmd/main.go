package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/PlopyBlopy/balance-keeper-service/internal/adapters/rest"
	"github.com/PlopyBlopy/balance-keeper-service/internal/shared/config"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
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

	err = app(*cfg, log)
	if err != nil {
		log.Panic("app error", zap.Error(err))
	}
}

func app(cfg config.Config, log *zap.Logger) error {
	ctx := context.Background()

	errChan := make(chan error)

	// postgres connection
	pgCfg, err := pgxpool.ParseConfig(cfg.DBConnString)
	if err != nil {
		return err
	}

	pgCfg.MaxConns = cfg.MaxConns
	pgCfg.MaxConns = cfg.MinConns

	pool, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		return err
	}

	// http server
	if cfg.IsProd {
		gin.SetMode(gin.ReleaseMode)
	}

	router := rest.NewRouter(1, pool)

	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler: router,
	}

	log.Info("Start listening to HTTP",
		zap.String("host", cfg.Host),
		zap.String("port", cfg.Port),
	)

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			errChan <- err
		}
	}()

	return nil
}
