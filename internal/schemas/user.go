package schemas

type RequestUserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RequestUserRegister struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
