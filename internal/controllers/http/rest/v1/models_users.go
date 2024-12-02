package v1

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateUserResponse struct {
	ID string `json:"id"`
}

type LoginUserResponse struct {
	Token string `json:"token"`
}
