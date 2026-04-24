package dto

import "time"

type ProductCreateDto struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Stock       int     `json:"stocke"`
	CategoryId  int64   `json:"categoryId"`
}

type ProductGetDto struct {
	Id           int64
	Name         string
	Description  string
	Price        float32
	Stock        int
	CategoryName string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
