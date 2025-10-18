BINARY := weight-tracker
SRC_DIR := src
OUT_DIR := out
WEBAPP_DIR := webapp
DOCKER_IMAGE := weight-tracker:local
PORT ?= 8080

.PHONY: all build run test fmt tidy docker-build docker-run clean help

all: build

help:
	@echo "Makefile for $(BINARY)"
	@echo
	@echo "Available targets:"
	@echo "  build        Build the Go binary into $(OUT_DIR)/$(BINARY)"
	@echo "  run          Run the server (go run)"
	@echo "  test         Run go tests"
	@echo "  fmt          Format source with gofmt"
	@echo "  tidy         Run go mod tidy"
	@echo "  docker-build Build a local Docker image ($(DOCKER_IMAGE))"
	@echo "  docker-run   Run the Docker image and forward port $(PORT)"
	@echo "  clean        Remove build outputs"

build:
	@echo "=> building $(BINARY)"
	@mkdir -p $(OUT_DIR)
	go build -o $(OUT_DIR)/$(BINARY) ./$(SRC_DIR)

run:
	@echo "=> running $(BINARY)"
	go run $(SRC_DIR)

test:
	@echo "=> running tests"
	go test ./...

fmt:
	@echo "=> formatting"
	gofmt -w .

tidy:
	@echo "=> tidy modules"
	go mod tidy

docker-build:
	@echo "=> building docker image $(DOCKER_IMAGE)"
	docker build -t $(DOCKER_IMAGE) .

docker-run:
	@echo "=> running docker image $(DOCKER_IMAGE) on :$(PORT)"
	docker run --rm -p $(PORT):$(PORT) $(DOCKER_IMAGE)

clean:
	@echo "=> cleaning"
	rm -rf $(OUT_DIR)/*
