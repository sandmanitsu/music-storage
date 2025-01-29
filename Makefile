start:
	docker-compose up -d
	goose up
	swag init -g cmd/app/main.go
	go run cmd/app/main.go

run:
	go run cmd/app/main.go

up:
	goose up

down:
	goose down

swag:
	swag init -g cmd/app/main.go