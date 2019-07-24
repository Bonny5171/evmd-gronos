SHELL := /bin/bash

GOBUILD   = go build -v
OS        = linux
ARCH      = amd64
APPNAME   = gronos
IMAGE     = everymind/gronos
VERSION   = 1.3.6

clean:
	@go clean -i -x ./...
	@rm -rf tmp

get:
	@go mod vendor

run:
	@source .env
	@go run app/*.go

build:
	@GOOS=$(OS) GOARCH=$(ARCH) $(GOBUILD) -o tmp/$(APPNAME) app/*.go

docker-login:
	@docker login -u everyci -pS0eusei@01

docker-logout:
	@docker logout

docker-build:
	@docker build \
		-t $(IMAGE):$(VERSION) \
		-t $(IMAGE):latest . 

docker-publish:
	@docker push $(IMAGE):$(VERSION)
	@docker push $(IMAGE):latest

docker-up:
	@docker-compose up -d

docker-down:
	@docker-compose down

tag:
	@git tag v$(VERSION) && git push --tags || :

default: build
