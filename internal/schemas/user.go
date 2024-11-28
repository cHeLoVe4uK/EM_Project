package schemas

type RequestLoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RequestRegisterUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
