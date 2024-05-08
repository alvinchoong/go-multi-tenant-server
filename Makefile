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

sqlc:
	which sqlc || go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.26.0
	sqlc generate 

bench:
	DATABASE_URL=$(DATABASE_URL) go test ./... -bench=. 

test-api:
	newman run docs/multi-tenant.postman_collection.json
