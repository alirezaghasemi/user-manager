package request

type CreateUserRequest struct {
	Name   string `json:"name" validate:"required,min=2,max=30"`
	Family string `json:"family" validate:"required,min=2,max=30"`
	Email  string `json:"email" validate:"required,email"`
	Age    int    `json:"age" validate:"required,gte=18,lte=120"`
}

type UpdateUserRequest struct {
	Name   *string `json:"name,omitempty" validate:"omitempty,min=2,max=30"`
	Family *string `json:"family,omitempty" validate:"omitempty,min=2,max=30"`
	Email  *string `json:"email,omitempty" validate:"omitempty,email"`
	Age    *int    `json:"age,omitempty" validate:"omitempty,gte=18,lte=120"`
}

type DeleteUserRequest struct {
}
