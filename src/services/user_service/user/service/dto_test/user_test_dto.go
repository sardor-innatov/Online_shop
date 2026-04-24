package dto_test

type UserCreateDto struct {
	FirstName string `json:"firstName" faker:"first_name"`
	LastName  string `json:"lastName" faker:"last_name"`
	Email     string `json:"email" faker:"email"`
	Password  string `json:"password" faker:"-"`
}

type UserUpdateDto struct {
	FirstName string `json:"firstName" faker:"first_name"`
	LastName  string `json:"lastName" faker:"last_name"`
	Email     string `json:"email" faker:"email"`
}

type UserGetDto struct {
	Id        int64  `json:"id" faker:"-"`
	FirstName string `json:"firstName" faker:"first_name"`
	LastName  string `json:"lastName" faker:"last_name"`
	Email     string `json:"email" faker:"email"`
}
