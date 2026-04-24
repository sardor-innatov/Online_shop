package dto

type CardCreateDto struct {
	Number        string `json:"number"`
	Name          string `json:"name"`
	PaymentSystem string `json:"paymentSystem"`
}

type BalanceUpdateDto struct {
	Amaunt float64 `json:"amaunt"`
}
