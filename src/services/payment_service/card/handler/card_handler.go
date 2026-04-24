package handler

import (
	"net/http"
	"online_shop/src/common/http/route"
	"online_shop/src/common/middleware"
	"online_shop/src/services/payment_service/card/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CardHandler interface{}

type cardHandler struct {
	service service.CardService
}

func NewCardHandler(group route.Group, conn *pgxpool.Pool) {

	handler := cardHandler{
		service: service.NewCardService(conn),
	}

	cardGroup := group.Group("/card", middleware.RequireAuth, middleware.CheckPermission(conn))
	{
		cardGroup.Handle("POST", "", handler.Create)
		cardGroup.Handle("PUT", "/balance/{id}", handler.UpdateBalance)
	}
}

// @Summary      create
// @Tags         card
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        card   body   dto.CardCreateDto  true  "card"
// @Success      200  {object}  nil
// @Failed       400  {object}  error
// @Failed       404  {object}  error
// @Failed       500  {object}  error
// @Router       /card [POST]
func (h *cardHandler) Create(ctx route.Context) error {

	id, err := h.service.Create(ctx)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusOK, map[string]any{"id": id})
}

// @Summary      udpate balance
// @Tags         card
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path   int64  true  "id"
// @Param        card   body	dto.BalanceUpdateDto   true  "card"
// @Success      200  {object}  nil
// @Failed       400  {object}  error
// @Failed       404  {object}  error
// @Failed       500  {object}  error
// @Router       /card/balance/{id} [PUT]
func (h *cardHandler) UpdateBalance(ctx route.Context) error {

	err := h.service.UpdateBalance(ctx)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusOK, map[string]any{})
}
