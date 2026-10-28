#!/usr/bin/make -rRw

APP 		:= gitolite
MODULE 		:= woodpecker-gitolite-forge
GO 			?= go
DOCKER 		?= docker 
GO_IMAGE 	?= golang:1.27
GOOS 		?= $(shell $(GO) env GOOS)
GOARCH 		?= $(shell $(GO) env GOARCH)
CGO_ENABLED ?= 0
BIN_DIR 	:= bin/$(GOOS)-$(GOARCH)
BIN 		:= $(BIN_DIR)/$(APP)
GOFILES 	:= $(wildcard *.go)

.DEFAULT_GOAL := build

.PHONY: all build test run clean deps fmt vet init arm64 amd64 docker-build

all: init deps fmt vet build-all

init:
	@test -f go.mod || $(GO) mod init $(MODULE) $(GO) mod tidy

deps:
	$(GO) mod download

build: $(BIN)

build-all: arm64 amd64

$(BIN): $(GOFILES) go.mod go.sum
	mkdir -p $(BIN_DIR)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -o $@ .

arm64:
	$(MAKE) build GOOS=linux GOARCH=arm64

amd64:
	$(MAKE) build GOOS=linux GOARCH=amd64

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

clean:
	rm -rf bin
