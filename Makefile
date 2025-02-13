include .env

migrate-up: 
	goose -dir $(MIGRATIONS_FILE) postgres $(DB_URL) up

migrate-down:
	goose -dir $(MIGRATIONS_FILE) postgres $(DB_URL) down

build:
	go build -o myapp $(PROJECT_PATH)

start:
	sudo docker-compose up --build -d
	
stop: 
	sudo docker-compose stop
