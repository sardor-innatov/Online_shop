package model

import "time"

type Product struct {
	Id          int64
	Name        string
	Description string
	Price       float32
	Stock       int
	CategoryId  int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
	