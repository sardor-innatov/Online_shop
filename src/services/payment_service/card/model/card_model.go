package model

type Card struct {
	Id               int64
	UserId           int64
	StripeMethodId   string
	StripeCustomerId string
	Last4            string
	PaymentSystem    string
	Name             string
}
