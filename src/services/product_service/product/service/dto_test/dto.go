package dto_test

type ProductCreateDto struct {
	Name        string  `json:"name" faker:"word"` 
	Description string  `json:"description" faker:"sentence"` 
	Price       float32 `json:"price" faker:"amount"` 
	Stock       int     `json:"stock" faker:"boundary_start=1, boundary_end=100"`
	CategoryId  int64   `json:"categoryId" faker:"boundary_start=11, boundary_end=11"`
}