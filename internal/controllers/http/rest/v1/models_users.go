package v1

type CreateUserRequest struct {
	Username string `json:"username" example:"Username" validate:"required"`
	Email    string `json:"email" example:"example@gmail.com" validate:"required,email"`
	Password string `json:"password" example:"secret1234!" validate:"required"`
}

type LoginUserRequest struct {
	Email    string `json:"email" example:"example@gmail.com" validate:"required,email"`
	Password string `json:"password" example:"secret1234!" validate:"required"`
}

type CreateUserResponse struct {
	ID string `json:"id" example:"UUID"`
}

type LoginUserResponse struct {
	Token string `json:"token" example:"JWT token"`
}
