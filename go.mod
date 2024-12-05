module github.com/cHeLoVe4uK/EM_Project

go 1.23.1

replace github.com/cHeLoVe4uK/EM_Project/internal/domain/models => ./internal/domain/models

require (
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
