package handler

import (
	"net/http"
	"online_shop/src/common/http/route"
	"online_shop/src/services/payment_service/payment/service"
)

type PaymentHandler interface{}

type paymentHandler struct {
	service service.PaymentService
}

func NewPaymentHandler(group route.Group) {

	handler := paymentHandler{
		service: service.NewPaymentService(),
	}

	paymentGroup := group.Group("/payment")
	{
		paymentGroup.Handle("POST", "/webhook", handler.WebhookHandler)
	}
}

// @Summary      webhook
// @Tags         payment
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        payment   body   dto.PaymentRequest  true  "payment"
// @Success      200  {object}  nil
// @Failed       400  {object}  error
// @Failed       404  {object}  error
// @Failed       500  {object}  error
// @Router       /payment/webhook [POST]
func (h *paymentHandler) WebhookHandler(ctx route.Context) error {

	err := h.service.CheckWebhook(ctx)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusOK, map[string]any{})
}
