DOCKER_REPO=test/cmeter
VERSION_FILE=VERSION
REV=$(shell git rev-parse --short HEAD)
ifeq ($(BUILD_VERSION),)
	BUILD_VERSION=$(shell cat $(VERSION_FILE))-$(REV)
endif

.PHONY: clean image

all: compile

clean:
	go clean ./...

compile:
	go build -ldflags "-X main.appVersion=$(BUILD_VERSION)" .

dist:
	GOOS=linux go build -ldflags "-X main.appVersion=$(BUILD_VERSION)" -o dist .

image:
	docker build -t $(DOCKER_REPO):BUILD_VERSION .
