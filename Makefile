.PHONY: setup build run test clean

APP := gitolite
MODULE := woodpecker-gitolite-forge

setup:
	@if [ ! -f go.mod ]; then go mod init $(MODULE); fi
	go mod tidy

build:
	mkdir -p bin
	go build -o bin/$(APP) .

run:
	go run .

test:
	go test ./...

clean:
	rm -rf bin/
