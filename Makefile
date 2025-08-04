IMAGE_NAME      := smart-pack

.PHONY: prepare lint local-lint test test-short codegen proto mockgen \
        coverage coverage-it-html build start \
        build-docker push-docker migrate-up migrate-down run-docker run-docker-compose

prepare:
	go mod download

lint:
	make prepare && ./bin/golangci-lint run --timeout 5m0s -c ./.golangci.yaml ./...

local-lint:
	golangci-lint run --timeout 5m0s

test:
	go test ./... -v

test-short:
	go test -short ./... -v

codegen:
	go get -d github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen && \
	oapi-codegen --config=.oapi-codegen.yml api/api.yml && \
	go mod tidy

mockgen:
	go generate ./...

coverage:
	make prepare && go test -race -coverprofile=coverage.txt -covermode=atomic ./...

coverage-it-html:
	go tool cover -html=cover-it.out

build:
	make prepare && go build .

start:
	make prepare && go run main.go api

build-docker:
	docker build --no-cache -t $(IMAGE_NAME):latest .

push-docker:
	docker push $(IMAGE_NAME):latest

run-docker:
	docker run -p 8080:8080 $(IMAGE_NAME):latest

run-docker-compose:
	docker-compose up -d

migrate-up:
	migrate -database "postgres://smartpack:smartpack@localhost:5432/smartpack?sslmode=disable" -path resources/db up

migrate-down:
	migrate -database "postgres://smartpack:smartpack@localhost:5432/smartpack?sslmode=disable" -path resources/db down