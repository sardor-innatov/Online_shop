package dto

type PaymentRequest struct {
	UserId           string
	Amount           float64  // Сумма в МИНИМАЛЬНЫХ единицах (центах/копейках)
	Currency         string // "usd", "eur"
	StripeCustomerId string // cus_...
	PaymentMethodId  string // pm_...
	OrderId          string // Для метаданных в Stripe (чтобы потом найти заказ)
}

type PaymentResponse struct {
	UserId        string
	Success       bool
	TransactionId string // ID от Stripe (pi_...)
	Status        string // "succeeded", "requires_action", "failed"
	ClientSecret  string // Нужен фронтенду, если банк потребует СМС (3D Secure)
	ErrorMessage  string
}
