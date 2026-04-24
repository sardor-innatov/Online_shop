package dto_test

type CategoryCreateDto struct {
	Name string 	`json:"name" faker:"first_name"`
}