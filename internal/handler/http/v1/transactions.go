package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/keremeti/iq-progers/internal/service"
	"github.com/keremeti/iq-progers/pkg/logger"
)

type transactionsHandler struct {
	l   *logger.Logger
	gts service.IGetTransactionsService
}

func newTransactionsHandler(handler *gin.RouterGroup, l *logger.Logger,
	gts service.IGetTransactionsService) {
	r := &transactionsHandler{l, gts}
	h := handler.Group("")
	{
		h.POST("", r.get)
	}
}

// @Summary		Get transactions
// @Description	Get transactions by filter
// @Tags			transactions
// @Accept			json
// @Produce		json
// @Param       limit query int false "Количество записей на одной странице"
// @Param       page query int false "Страница списка"
// @Param       request body v1.filterRequest true "Фильтр для транзакций"
// @Success		200	{object}	v1.WrapperResponse{data=[]v1.transactionResponse}
// @Failure		400	{object}	v1.WrapperResponse
// @Failure		404	{object}	v1.WrapperResponse
// @Failure		500	{object}	v1.WrapperResponse
// @Router			/v1/transactions 	[post]
func (h *transactionsHandler) get(ctx *gin.Context) {
	req := filterRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.l.Debug("transactionsHandler - get - ctx.ShouldBindJSON: ", "err", err.Error())
		err := fmt.Errorf("необходимо передать фильтр для транзакций")
		var wraper = NewErrorWrapper(400, err)
		ctx.JSONP(http.StatusOK, wraper)
		return
	}
	query := ctx.Request.URL.Query()
	limit := 10
	if query["limit"] != nil {
		l, err := strconv.Atoi(query["limit"][0])
		if err == nil {
			limit = l
		}
	}
	page := 1
	if query["page"] != nil {
		p, err := strconv.Atoi(query["page"][0])
		if err == nil {
			page = p
		}
	}

	transactions, err := h.gts.Execute(ctx, req.ToDomain(), limit, page)
	if err != nil {
		h.l.Warn(err.Error())
		var wraper = NewErrorWrapper(500, err)
		ctx.JSON(http.StatusOK, wraper)
		return
	}
	res := []transactionResponse{}
	for i := 0; i < len(transactions.Slice); i++ {
		t := newTransactionResponse(transactions.Slice[i])
		res = append(res, t)
	}
	var wraper = NewDataWrapper(res)
	j, err := json.Marshal(NewPaginationResponse(&transactions))
	if err != nil {
		h.l.Info(err.Error())
		var wraper = NewErrorWrapper(400, err)
		ctx.JSON(http.StatusOK, wraper)
		return
	}
	ctx.Header("X-Pagination", string(j))
	ctx.JSON(http.StatusOK, wraper)
}
