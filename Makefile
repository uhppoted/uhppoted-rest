VERSION = v0.5.1x
LDFLAGS = -ldflags "-X uhppote.VERSION=$(VERSION)" 
CLI     = ./bin/uhppote-cli

SERIALNO ?= 405419896
CARD     ?= 65538
DOOR     ?= 3
DEVICEIP ?= 192.168.1.125
DATETIME  = $(shell date "+%Y-%m-%d %H:%M:%S")
LISTEN   ?= 192.168.1.100:60001
DEBUG    ?= --debug

all: test      \
	 benchmark \
     coverage

clean:
	go clean
	rm -rf bin

format: 
	go fmt ./...

build: format
	mkdir -p bin
	go build -o bin ./...

test: build
	go test ./...

vet: build
	go vet ./...

lint: build
	golint ./...

benchmark: build
	go test -bench ./...

coverage: build
	go test -cover ./...

release: test vet
	mkdir -p dist/$(DIST)/windows
	mkdir -p dist/$(DIST)/darwin
	mkdir -p dist/$(DIST)/linux
	mkdir -p dist/$(DIST)/arm7
	env GOOS=linux   GOARCH=amd64       go build -o dist/$(DIST)/linux   ./...
	env GOOS=linux   GOARCH=arm GOARM=7 go build -o dist/$(DIST)/arm7    ./...
	env GOOS=darwin  GOARCH=amd64       go build -o dist/$(DIST)/darwin  ./..
	env GOOS=windows GOARCH=amd64       go build -o dist/$(DIST)/windows ./...

release-tar: release
	tar --directory=dist --exclude=".DS_Store" -cvzf dist/$(DIST).tar.gz $(DIST)

debug: build
	go test ./...

version: build
	./bin/uhppoted-rest version

help: build
	./bin/uhppoted-rest help
	./bin/uhppoted-rest help commands
	./bin/uhppoted-rest help version
	./bin/uhppoted-rest help help

daemonize: build
	sudo ./bin/uhppoted-rest daemonize

undaemonize: build
	sudo ./bin/uhppoted-rest undaemonize



