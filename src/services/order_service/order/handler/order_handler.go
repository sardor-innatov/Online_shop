package handler

import (
	"net/http"
	"online_shop/src/common/config"
	"online_shop/src/common/http/route"
	"online_shop/src/common/middleware"
	"online_shop/src/services/order_service/order/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderHandler interface{}

type orderHandler struct {
	service service.OrderService
}

func NewOrderHandler(conn *pgxpool.Pool, group route.Group) {

	handler := orderHandler{
		service: service.NewOrderService(conn),
	}

	orderGroup := group.Group("/order", middleware.RequireAuth, middleware.CheckPermission(conn))
	{
		orderGroup.Handle("POST", "", handler.Create, middleware.IdempotencyMiddleware(config.GetRedis()))
	}

}

// @Summary      create
// @Tags         order
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        order   body   dto.OrderCreateDto  true  "order"
// @Success      200  {object}  nil
// @Failed       400  {object}  error
// @Failed       404  {object}  error
// @Failed       500  {object}  error
// @Router       /order [POST]
func (h *orderHandler) Create(ctx route.Context) error {

	id, err := h.service.Create(ctx)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusOK, map[string]any{"id": id})
}
