VERSION_FILE=VERSION
REV=$(shell git rev-parse --short HEAD)
APP_VERSION=$(shell cat $(VERSION_FILE))-$(REV)

.PHONY: clean

all: compile

clean:
	go clean ./...

compile:
	CGO_ENABLED=0 go build -ldflags "-X main.appVersion=$(APP_VERSION)" main.go
