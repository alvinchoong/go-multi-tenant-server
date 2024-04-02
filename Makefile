SHELL = /bin/bash -u -e -o pipefail

# `make` applies env vars from `.env`
include .env

up:
	docker-compose up -d --remove-orphans

down:
	docker-compose down

run:
	which air || go install github.com/cosmtrek/air@latest
	$(shell cat .env | egrep -v '^#' | xargs -0) \
	air --build.delay=1000 \
		--build.cmd "go build -o bin/server cmd/server/main.go" \
		--build.bin "./bin/server" \
		--build.include_ext "go" \
		--build.exclude_dir "tmp,vendor,testdata" \

seed: migrate
	psql postgres://su:password@127.0.0.1:5433/appdb -f database/seed/silo.sql

migrate:
	@make go-migrate DATABASE_URL=postgres://su:password@127.0.0.1:5432/appdb?sslmode=disable
	@make go-migrate DATABASE_URL=postgres://su:password@127.0.0.1:5433/appdb?sslmode=disable

go-migrate: wait-for-pg
	which migrate || go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	migrate -path ./database/migrations -database "$(DATABASE_URL)" up

wait-for-pg:
	@while ! pg_isready -q -d $(DATABASE_URL); do \
		echo "Waiting for PostgreSQL to be available..."; \
		sleep 1; \
	done

bench:
	DATABASE_POOL_RW_URL=$(DATABASE_POOL_RW_URL) \
	go test ./... -bench=. 

reset-db:
	docker-compose down db-pool db-silo
	docker-compose up db-pool db-silo --force-recreate --build --detach
	make seed

test-api:
	newman run docs/multi-tenant.postman_collection.json
