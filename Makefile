#!/usr/bin/make -rRw

.PHONY: setup build run test clean

APP		:= gitolite
MODULE 	:= woodpecker-gitolite-forge

SRCS	= $(wildcard *.go)

go.mod:
	go mod init $(MODULE)

go.sum:
	go mod tidy

build: go.mod go.sum bin $(SRCS)
	go build -o bin/$(APP) .

bin:
	mkdir -p bin

clean:
	rm -rf bin/
