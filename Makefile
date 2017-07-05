.PHONY: build

ImageName=malcolm	
PACKAGES = $(shell go list ./... | grep -v /vendor/)

ifneq ($(shell uname), Darwin)
	EXTLDFLAGS = -extldflags "-static" $(null)
else
	EXTLDFLAGS =
endif

BUILD_NUMBER=$(shell git rev-parse --short HEAD)

all: build_static

test:
	go test -cover $(PACKAGES)

build: build_static build_cross

build_static:
	mkdir -p make/release
	go build -o  make/release/malcolm -ldflags '${EXTLDFLAGS}-X github.com/u2takey/malcolm/version.VersionDev=build.$(BUILD_NUMBER)' github.com/u2takey/malcolm/cmd

build_cross:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-X github.com/u2takey/malcolm/version.VersionDev=build.$(BUILD_NUMBER)' -o make/release/linux/amd64/malcolm   github.com/u2takey/malcolm/cmd
	# GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-X github.com/u2takey/malcolm/version.VersionDev=build.$(malcolm_BUILD_NUMBER)' -o make/release/darwin/amd64/malcolm  github.com/u2takey/malcolm/cmd

build_tar:
	tar -cvzf make/release/linux/amd64/malcolm.tar.gz   -C make/release/linux/amd64/malcolm
	tar -cvzf make/release/darwin/amd64/malcolm.tar.gz  -C make/release/darwin/amd64/malcolm

build_compose:
	cd make && docker-compose build && cd - 

run_compose:
	cd make && docker-compose up && cd - 

