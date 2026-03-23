.PHONY: all build run test clean frontend frontend-install frontend-dev

all: build

build:
	go build -o bin/relay ./cmd/server

run:
	go run ./cmd/server -demo

run-prod:
	go run ./cmd/server -config config.json

test:
	go test ./...

clean:
	rm -rf bin/

frontend-install:
	cd frontend && npm install

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build
