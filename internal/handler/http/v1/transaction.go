package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/keremeti/iq-progers/internal/service"
	"github.com/keremeti/iq-progers/pkg/logger"
)

type transactionHandler struct {
	l                    *logger.Logger
	topUpBalanceService  service.ITopUpBalanceService
	transferMoneyService service.ITransferMoneyService
}

func newTransactionHandler(handler *gin.RouterGroup, l *logger.Logger,
	tub service.ITopUpBalanceService,
	tm service.ITransferMoneyService) {
	r := &transactionHandler{l, tub, tm}

	h := handler.Group("")
	{
		h.POST("top-up", r.topUpBalance)
		h.POST("transfer", r.transferMoney)
	}
}

// @Summary		Top up balance
// @Description	Top up balance
// @Tags			transactions
// @Accept			json
// @Produce		json
// @Success		200	{object}	v1.WrapperResponse{data=v1.transactionResponse}
// @Failure		400	{object}	v1.WrapperResponse
// @Failure		404	{object}	v1.WrapperResponse
// @Failure		500	{object}	v1.WrapperResponse
// @Router			/v1/transactions/top-up [post]
func (h *transactionHandler) topUpBalance(ctx *gin.Context) {
	req := rechargeRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.l.Debug("transactionHandler - topUpBalance - ctx.ShouldBindJSON: ", "err", err.Error())
		err := fmt.Errorf("необходимо передать данные для пополнения")
		var wraper = NewErrorWrapper(400, err)
		ctx.JSONP(http.StatusOK, wraper)
		return
	}
	entity, err := req.ToDomain()
	if err != nil {
		h.l.Warn("transactionHandler - topUpBalance - req.ToDomain: ", "err", err.Error())
		var wraper = NewErrorWrapper(500, err)
		ctx.JSON(http.StatusOK, wraper)
		return
	}

	transaction, err := h.topUpBalanceService.Execute(ctx, entity)
	if err != nil {
		h.l.Warn("transactionHandler - topUpBalance - topUpBalanceService.Execute: ", "err", err.Error())
		var wraper = NewErrorWrapper(500, err)
		ctx.JSON(http.StatusOK, wraper)
		return
	}
	h.l.Info("transactionHandler - topUpBalance - topUpBalanceService.Execute completed")

	var wraper = NewDataWrapper(newTransactionResponse(transaction))
	ctx.JSON(http.StatusOK, wraper)
}

// @Summary		Transfer money
// @Description	Transfer money between users
// @Tags			transactions
// @Accept			json
// @Produce		json
// @Param       request body v1.remittanceRequest true "Данные для перевода средств"
// @Success		200	{object}	v1.WrapperResponse
// @Failure		400	{object}	v1.WrapperResponse
// @Failure		404	{object}	v1.WrapperResponse
// @Failure		500	{object}	v1.WrapperResponse
// @Router			/v1/transactions/transfer [post]
func (h *transactionHandler) transferMoney(ctx *gin.Context) {
	req := remittanceRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.l.Debug("transactionHandler - transferMoney - ctx.ShouldBindJSON: ", "err", err.Error())
		err := fmt.Errorf("необходимо передать данные для перевода")
		var wraper = NewErrorWrapper(400, err)
		ctx.JSONP(http.StatusOK, wraper)
		return
	}
	entity, err := req.ToDomain()
	if err != nil {
		h.l.Warn("transactionHandler - transferMoney - req.ToDomain: ", "err", err.Error())
		err := fmt.Errorf("необходимо указать цифрами сумму перевода")
		var wraper = NewErrorWrapper(400, err)
		ctx.JSON(http.StatusOK, wraper)
		return
	}

	err = h.transferMoneyService.Execute(ctx, entity)
	if err != nil {
		h.l.Warn("transactionHandler - transferMoney - transferMoneyService.Execute: ", "err", err.Error())
		var wraper = NewErrorWrapper(400, err)
		ctx.JSON(http.StatusOK, wraper)
		return
	}
	h.l.Info("transactionHandler - transferMoney - transferMoneyService.Execute completed")

	var wraper = New200Wrapper()
	ctx.JSON(http.StatusOK, wraper)
}
