package model

import "time"

type Order struct {
	Id         int64
	UserId     int64
	TotalPrice float32
	Status     orderStatus
	CreatedAt  time.Time
}

type orderStatus string

const (
	PENDING   orderStatus = "pending"
	PAID      orderStatus = "paid"
	FAILED    orderStatus = "failed"
	CANCELLED orderStatus = "cancelled"
)

type OrderItem struct {
	OrderId   int64
	ProductId int64
	Quantity  int
	Price     float32
}
