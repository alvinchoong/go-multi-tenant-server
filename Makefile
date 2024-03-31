SHELL = /bin/bash -u -e -o pipefail

# `make` applies env vars from `.env`
include .env

up:
	docker-compose up -d --remove-orphans

down:
	docker-compose down

server-run-app:
	which air || go install github.com/cosmtrek/air@latest
	$(shell cat .env | egrep -v '^#' | xargs -0) \
	air --build.delay=1000 \
		--build.cmd "go build -o bin/server cmd/server/main.go" \
		--build.bin "./bin/server" \
		--build.include_ext "go" \
		--build.exclude_dir "tmp,vendor,testdata" \

wait-for-pg:
	@while ! pg_isready -q -d $(DATABASE_URL); do \
		echo "Waiting for PostgreSQL to be available..."; \
		sleep 1; \
	done

migrate:
	@make go-migrate DATABASE_URL=$(DATABASE_POOL_SU_URL)
	@make go-migrate DATABASE_URL=$(DATABASE_SILO_SU_URL)

go-migrate: wait-for-pg
	which migrate || go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	migrate -path ./database/migrations -database "$(DATABASE_URL)?sslmode=disable" up

seed: migrate
	psql $(DATABASE_SILO_RW_URL) -f database/seed/silo.sql

db-console-pool: DATABASE_URL=$(DATABASE_POOL_SU_URL) # set to DATABASE_POOL_RW_URL to test RLS
db-console-pool:
	psql $(DATABASE_URL)

db-console-silo: DATABASE_URL=$(DATABASE_SILO_SU_URL) # set to DATABASE_SILO_RW_URL to test RLS
db-console-silo:
	psql $(DATABASE_URL)

bench:
	DATABASE_POOL_RW_URL=$(DATABASE_POOL_RW_URL) \
	go test ./... -bench=. 
