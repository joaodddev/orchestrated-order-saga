.PHONY: up down migrate-up migrate-down run

up:
	docker-compose up -d mysql-orchestrator kafka

migrate-up:
	migrate -database "mysql://root:root@tcp(localhost:3308)/orchestrator_db" \
		-path internal/infrastructure/persistence/mysql/migrations up

migrate-down:
	migrate -database "mysql://root:root@tcp(localhost:3308)/orchestrator_db" \
		-path internal/infrastructure/persistence/mysql/migrations down

run:
	go run cmd/worker/main.go