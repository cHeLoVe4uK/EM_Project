package models

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Claims struct {
	UserID   string
	Username string
}
