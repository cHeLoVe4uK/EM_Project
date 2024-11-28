package models

type User struct {
	ID       string `bson:"_id,omitempty"`
	Email    string `bson:"email"`
	Username string `bson:"username"`
	Password string `bson:"password"`
}
