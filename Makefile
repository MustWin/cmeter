VERSION_FILE=VERSION
REV=$(shell git rev-parse --short HEAD)
ifeq ($(BUILD_VERSION),)
	BUILD_VERSION=$(shell cat $(VERSION_FILE))-$(REV)
endif

.PHONY: clean

all: compile

clean:
	go clean ./...

compile:
	go build -ldflags "-X main.appVersion=$(BUILD_VERSION)" .
