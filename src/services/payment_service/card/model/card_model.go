package model

type Card struct {
	Id            int64
	UserId        int64
	Number        string
	Name          string
	PaymentSystem string
	Balance       float32
}
