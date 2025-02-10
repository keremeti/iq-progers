// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/keremeti/iq-progers/config"

	v1 "github.com/keremeti/iq-progers/internal/handler/http/v1"
	"github.com/keremeti/iq-progers/internal/service"
	"github.com/keremeti/iq-progers/internal/service/repo"
	"github.com/keremeti/iq-progers/pkg/httpserver"
	"github.com/keremeti/iq-progers/pkg/logger"
	"github.com/keremeti/iq-progers/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Env)
	l.Debug("app - Run - Logger started")

	l.Info("app - Run", "Config", cfg.Env.ToString())

	err := applyMigrations(cfg.Goose)
	if err != nil {
		l.Warn(fmt.Sprintf("app - Run - applyMigrations: %v", err))
	}
	l.Debug("app - Run - Migrations done")

	pg, err := postgres.New(cfg.PG.Url, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Warn(fmt.Sprintf("app - Run - postgres.New: %v", err))
	}
	defer pg.Close()
	l.Debug("app - Run - Postgres started")

	balRepo := repo.NewBalanceRepo(pg)
	tranRepo := repo.NewTransactionRepo(pg)
	transRepo := repo.NewTransactionsRepo(pg)
	l.Debug("app - Run - Repositories started")

	tubSer := service.NewTopUpBalanceService(
		l,
		tranRepo,
	)
	gtSer := service.NewGetTransactionsService(
		l,
		transRepo,
	)
	tfmSer := service.NewTransferMoneyService(
		l,
		tranRepo,
		balRepo,
	)
	l.Debug("app - Run - Services started")

	handler := gin.New()
	v1.NewHandler(handler, l, cfg, tubSer, tfmSer, gtSer)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	l.Debug("app - Run - HTTP Server started")

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Sprintf("app - Run - httpServer.Notify: %v", err))
	}

	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Sprintf("app - Run - httpServer.Shutdown: %v", err))
	}
}
