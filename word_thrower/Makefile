
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

BINARY_NAME=randgen
CWD=$(shell pwd)
SRC=$(CWD)/src

build:
	$(GOBUILD) -o bin/$(BINARY_NAME) -v $(SRC)/main.go