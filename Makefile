SHELL = /bin/bash -u -e -o pipefail

# `make` applies env vars from `.env`
include .env

up:
	docker-compose up -d --remove-orphans

down:
	docker-compose down

run:
	$(shell cat .env | egrep -v '^#' | xargs -0) \
	go run cmd/server/main.go

dev:
	which air || go install github.com/cosmtrek/air@latest
	$(shell cat .env | egrep -v '^#' | xargs -0) \
	air --build.delay=1000 \
		--build.cmd "go build -o bin/server cmd/server/main.go" \
		--build.bin "./bin/server" \
		--build.include_ext "go" \
		--build.exclude_dir "tmp,vendor,testdata" \

# seed: migrate
# 	psql postgres://su:password@127.0.0.1:5433/appdb -f database/seed/silo.sql

DATABASE_URL=postgres://su:password@127.0.0.1:5432/appdb
migrate: wait-for-pg
	which migrate || go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	migrate -path ./database/migrations -database "$(DATABASE_URL)?sslmode=disable" up

wait-for-pg:
	@while ! pg_isready -q -d $(DATABASE_URL); do \
		echo "Waiting for PostgreSQL to be available..."; \
		sleep 1; \
	done

sqlc:
	which sqlc || go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.26.0
	sqlc generate 

bench:
	DATABASE_URL=$(DATABASE_URL) go test ./... -bench=. 

reset-db:
	docker-compose down db-pool
	docker-compose up db-pool --force-recreate --build --detach
	make migrate

test-api:
	newman run docs/multi-tenant.postman_collection.json
