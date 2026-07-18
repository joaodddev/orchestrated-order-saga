.PHONY: up migrate-up migrate-down run

up:
	cd .. && docker-compose up -d mysql-order kafka

migrate-up:
	migrate -database "mysql://root:root@tcp(localhost:3309)/order_db" \
		-path internal/infrastructure/persistence/mysql/migrations up

migrate-down:
	migrate -database "mysql://root:root@tcp(localhost:3309)/order_db" \
		-path internal/infrastructure/persistence/mysql/migrations down

run:
	go run cmd/worker/main.go