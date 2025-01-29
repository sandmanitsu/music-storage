run:
	go run cmd/app/main.go

up:
	goose up

down:
	goose down

swag:
	swag init -g cmd/app/main.go