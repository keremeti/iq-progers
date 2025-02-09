// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/keremeti/iq-progers/config"
	"github.com/keremeti/iq-progers/internal/service"
	"github.com/keremeti/iq-progers/pkg/logger"
)

func NewHandler(handler *gin.Engine, l *logger.Logger, cfg *config.Config,
	tub service.ITopUpBalanceService,
	tm service.ITransferMoneyService,
	gs service.IGetTransactionsService) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// Routers
	h := handler.Group("v1/transactions")
	{
		newTransactionHandler(h, l, tub, tm)
		newTransactionsHandler(h, l, gs)
	}
}
