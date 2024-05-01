DIST   ?= development
DEBUG  ?= --debug
CMD     = ./bin/uhppoted-rest
DOCKER ?= ghcr.io/uhppoted/restd:latest

.DEFAULT_GOAL := test
.PHONY: update
.PHONY: update-release
.PHONY: open-api

all: test      \
	 benchmark \
     coverage

clean:
	go clean
	rm -rf bin

update:
	go get -u github.com/uhppoted/uhppote-core@master
	go get -u github.com/uhppoted/uhppoted-lib@main
	go get -u golang.org/x/sys

update-release:
	go get -u github.com/uhppoted/uhppote-core
	go get -u github.com/uhppoted/uhppoted-lib
	go get -u golang.org/x/sys

format: 
	go fmt ./...

build: format
	mkdir -p bin
	go build -trimpath -o bin ./...

test: 
	go test ./...

benchmark: build
	go test -bench ./...

coverage: build
	go test -cover ./...

vet: 
	go vet ./...

lint: 
	env GOOS=darwin  GOARCH=amd64 staticcheck ./...
	env GOOS=linux   GOARCH=amd64 staticcheck ./...
	env GOOS=windows GOARCH=amd64 staticcheck ./...

vuln:
	govulncheck ./...

open-api:
	swagger-cli bundle documentation/openapi/uhppoted-api.yaml --outfile generated.yaml --type yaml

build-all: build test vet lint
	mkdir -p dist/$(DIST)/windows
	mkdir -p dist/$(DIST)/darwin
	mkdir -p dist/$(DIST)/linux
	mkdir -p dist/$(DIST)/arm
	mkdir -p dist/$(DIST)/arm7
	env GOOS=linux   GOARCH=amd64         GOWORK=off go build -trimpath -o dist/$(DIST)/linux   ./...
	env GOOS=linux   GOARCH=arm64         GOWORK=off go build -trimpath -o dist/$(DIST)/arm     ./...
	env GOOS=linux   GOARCH=arm   GOARM=7 GOWORK=off go build -trimpath -o dist/$(DIST)/arm7    ./...
	env GOOS=darwin  GOARCH=amd64         GOWORK=off go build -trimpath -o dist/$(DIST)/darwin  ./...
	env GOOS=windows GOARCH=amd64         GOWORK=off go build -trimpath -o dist/$(DIST)/windows ./...

release: update-release build-all
	find . -name ".DS_Store" -delete
	tar --directory=dist --exclude=".DS_Store" -cvzf dist/$(DIST).tar.gz $(DIST)
	cd dist; zip --recurse-paths $(DIST).zip $(DIST)

publish: release
	echo "Releasing version $(VERSION)"
	rm -f dist/development.tar.gz
	gh release create "$(VERSION)" "./dist/uhppoted-rest_$(VERSION).tar.gz" "./dist/uhppoted-rest_$(VERSION).zip" --draft --prerelease --title "$(VERSION)-beta" --notes-file release-notes.md

debug: build
	$(CMD) run --console

godoc:
	godoc -http=:80	-index_interval=60s

version: build
	$(CMD) version

help: build
	$(CMD) help
	$(CMD) help commands
	$(CMD) help version
	$(CMD) help config
	$(CMD) help help
	$(CMD) help run
	$(CMD) help daemonize
	$(CMD) help undaemonize

config: build
	$(CMD) config

daemonize: build
	sudo $(CMD) daemonize

undaemonize: build
	sudo $(CMD) undaemonize

run: build
	$(CMD) --log-level debug --console

swagger: 
	docker run --detach --publish 80:8080 --rm swaggerapi/swagger-editor 
	open http://127.0.0.1:80

docker: docker-dev docker-ghcr
	cd docker && find . -name .DS_Store -delete && rm -f compose.zip && zip --recurse-paths compose.zip compose

docker-dev: build
	rm -rf dist/docker/dev/*
	mkdir -p dist/docker/dev
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o dist/docker/dev ./...
	cp docker/dev/Dockerfile    dist/docker/dev
	cp docker/dev/ca.cert       dist/docker/dev
	cp docker/dev/uhppoted.cert dist/docker/dev
	cp docker/dev/uhppoted.conf dist/docker/dev
	cp docker/dev/uhppoted.key  dist/docker/dev
	cd dist/docker/dev && docker build --no-cache -f Dockerfile -t uhppoted/uhppoted-rest-dev .

docker-ghcr: build
	rm -rf dist/docker/ghcr/*
	mkdir -p dist/docker/ghcr
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o dist/docker/ghcr ./...
	cp docker/ghcr/Dockerfile    dist/docker/ghcr
	cp docker/ghcr/uhppoted.conf dist/docker/ghcr
	cd dist/docker/ghcr && docker build --no-cache -f Dockerfile -t $(DOCKER) .

docker-run-dev:
	docker run --publish 8080:8080 --name restd --rm uhppoted/uhppoted-rest-dev
	sleep 1

docker-run-ghcr:
	docker run --publish 8080:8080 --publish 8443:8443 --name restd --mount source=uhppoted-rest,target=/usr/local/etc/uhppoted --rm ghcr.io/uhppoted/restd
	sleep 1

docker-compose:
	cd docker/compose && docker compose up

docker-clean:
	docker image     prune -f
	docker container prune -f

get-controllers:
	curl -X 'GET' 'http://127.0.0.1:8080/uhppote/device' -H 'accept: application/json' | jq .
