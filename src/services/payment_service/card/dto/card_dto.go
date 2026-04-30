package dto

type CardCreateDto struct {
	StripeMethodId   string `json:"stripeMethodId"`
	Last4            string `json:"last4"`
	StripeCustomerId string `json:"stripeCustomerId"`
	Name             string `json:"name"`
	PaymentSystem    string `json:"paymentSystem"`
}

type BalanceUpdateDto struct {
	Amaunt float64 `json:"amaunt"`
}

type InitiateSetupRequestDto struct {
	Email string `json:"email"` // Нужно Stripe для создания Customer, если его нет
}

type InitiateSetupResponseDto struct {
	ClientSecret     string `json:"client_secret"`
	StripeCustomerId string `json:"stripe_customer_id"`
}
