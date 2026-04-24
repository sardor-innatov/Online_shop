package dto

type OrderCreateDto struct {
	Items  []OrderItemCreateDto `json:"items"`
	CardId int64                `json:"cardId"`
}

type OrderGetDto struct {
	OrderId    int64         `json:"orderId"`
	UserId     int64         `json:"userId"`
	Items      []OrderGetDto `json:"items"`
	TotalPrice float32       `json:"totalPrice"`
	Status     string        `json:"status"`
}

type OrderItemCreateDto struct {
	ProductId int64 `json:"productId"`
	Quantity  int   `json:"quantity"`
}

type OrderItemGetDto struct {
	ProductId   int64   `json:"productId"`
	ProductName string  `json:"productName"`
	Quantity    int     `json:"quantity"`
	Price       float32 `json:"price"`
}
