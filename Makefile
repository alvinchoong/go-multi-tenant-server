SHELL = /bin/bash -u -e -o pipefail

# `make` applies env vars from `.env`
include .env

server-run-app:
	which air || go install github.com/cosmtrek/air@latest
	air --build.delay=1000 \
		--build.cmd "go build -o bin/server main.go" \
		--build.bin "./bin/server" \
		--build.include_ext "go,tpl,tmpl,html,js,css" \
		--build.exclude_dir "tmp,vendor,testdata" \
