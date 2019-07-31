SHELL := /bin/bash

GCPROJECT = evmdsfa
GOBUILD   = go build -v
OS        = linux
ARCH      = amd64
APPNAME   = gronos
IMAGE     = evmd-gronos
VERSION   = 1.0.1

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

docker-auth:
	@gcloud auth configure-docker

docker-clean:
	@docker rmi -f $$(docker images $(IMAGE) --format "{{.ID}}" | sort --unique)

docker-build:
	@docker build \
		-t $(IMAGE):$(VERSION) \
		-t $(IMAGE):latest . 

docker-tags:
	@docker tag $(IMAGE) gcr.io/$(GCPROJECT)/$(IMAGE):$(VERSION)
	@docker tag $(IMAGE) gcr.io/$(GCPROJECT)/$(IMAGE):latest
	@gcloud container images untag gcr.io/$(GCPROJECT)/$(IMAGE):latest --quiet

docker-publish:
	@docker push gcr.io/$(GCPROJECT)/$(IMAGE):$(VERSION)
	@docker push gcr.io/$(GCPROJECT)/$(IMAGE):latest

docker-unpublish:
	@gcloud container images delete gcr.io/$(GCPROJECT)/$(IMAGE):$(VERSION) --force-delete-tags

docker-up:
	@docker-compose up -d

docker-down:
	@docker-compose down

tag:
	@git tag v$(VERSION) && git push --tags || :

default: build