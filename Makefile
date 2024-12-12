build:
	@echo "Building app..."
	@go build -o ./.bin/app ./cmd/app

run:build
	@./.bin/app

docker:
	@echo "Starting services in docker..."
	@docker compose -f docker-compose.local.yaml up --build -d 

stop:
	@echo "Stopping services in docker..."
	@docker compose -f docker-compose.local.yaml stop

swagger:
	@echo "Generating swagger doc..."
	@go install github.com/swaggo/swag/cmd/swag@v1.8.12
	@swag fmt
	@swag init -g cmd/app/main.go -o api/swagger

unit_test:
	@echo "Running unit tests..."
	@go test -v ./internal/... -cover -race
