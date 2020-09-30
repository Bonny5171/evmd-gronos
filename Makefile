SHELL := /bin/bash

GCPROJECT = evmdsfa
GOBUILD   = go build -v
OS        = linux
ARCH      = amd64
APPNAME   = gronos
IMAGE     = evmd-gronos
VERSION   = $$(git tag --sort=-version:refname | head -n 1)

build:
	@GOOS=$(OS) GOARCH=$(ARCH) $(GOBUILD) -o tmp/$(APPNAME) *.go

clean:
	@go clean -i -x ./...
	@rm -rf tmp

run:
	@source .env && go run -ldflags "-X main.version=$(VERSION)" main.go

gae-deploy-dev:
	@gcloud config set account $(shell gcloud config list account --format "value(core.account)") && gcloud config set project evmdsfa && gcloud config list
	@gcloud app deploy app.stg.dev.yaml --version=$(subst .,-,$(shell git tag --sort=-version:refname | head -n 1))

gae-deploy-qa:
	@gcloud config set account $(shell gcloud config list account --format "value(core.account)") && gcloud config set project evmdsfa && gcloud config list
	@gcloud app deploy app.stg.qa.yaml --version=$(subst .,-,$(shell git tag --sort=-version:refname | head -n 1))

gae-deploy-snd:
	@gcloud config set account $(shell gcloud config list account --format "value(core.account)") && gcloud config set project evmdsfa-snd && gcloud config list
	@gcloud app deploy app.snd.yaml --version=$(subst .,-,$(shell git tag --sort=-version:refname | head -n 1))

gae-deploy-grendene-dev:
	@gcloud config set account $(shell gcloud config list account --format "value(core.account)") && gcloud config set project crmgrendene && gcloud config list
	@gcloud app deploy app.grendene.dev.yaml --version=$(subst .,-,$(shell git tag --sort=-version:refname | head -n 1))

gae-deploy-grendene-qa:
	@gcloud config set account $(shell gcloud config list account --format "value(core.account)") && gcloud config set project crmgrendene && gcloud config list
	@gcloud app deploy app.grendene.qa.yaml --version=$(subst .,-,$(shell git tag --sort=-version:refname | head -n 1))

docker-auth:
	@gcloud auth configure-docker

docker-build:
	@docker build \
		-t $(IMAGE):$(VERSION) \
		-t $(IMAGE):latest . 

docker-clean:
	@docker rmi -f $$(docker images $(IMAGE) --format "{{.ID}}" | sort --unique)
	@docker rmi $$(docker images --filter "dangling=true" -q)

docker-delete:
	@gcloud container images delete gcr.io/$(shell gcloud config list project --format "value(core.project)")/$(IMAGE):$(VERSION) --force-delete-tags

docker-push:
	@docker push gcr.io/$(shell gcloud config list project --format "value(core.project)")/$(IMAGE):$(VERSION)
	@docker push gcr.io/$(shell gcloud config list project --format "value(core.project)")/$(IMAGE):latest

docker-publish: docker-build docker-tags docker-push docker-clean

docker-tags:
	@docker tag $(IMAGE) gcr.io/$(shell gcloud config list project --format "value(core.project)")/$(IMAGE):$(VERSION)
	@docker tag $(IMAGE) gcr.io/$(shell gcloud config list project --format "value(core.project)")/$(IMAGE):latest

docker-up:
	@docker-compose up -d

docker-down:
	@docker-compose down

git-merge-publish:
	@git checkout master
	@git merge develop
	@git tag -a v$(VERSION) -m "Release v$(VERSION)"
	@git checkout develop
	@git push --all && git push --tags

git-tag:
	@git tag -a v$(VERSION) -m "Release v$(VERSION)"

default: build