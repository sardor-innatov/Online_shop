package dto

type CategoryCreateDto struct {
	Name string `json:"name"`
}

type CategoryGetDto struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
