include golang.mk
include sfncli.mk

.PHONY: all test build run
SHELL := /bin/bash
PKG := github.com/Clever/signalfx-janitor
PKGS = $(shell go list ./... | grep -v "vendor/"  | grep -v /vendor)
EXECUTABLE = $(shell basename $(PKG))
SFNCLI_VERSION := latest
$(eval $(call golang-version-check,1.16))

all: test build

test: $(PKGS)
$(PKGS): golang-test-all-deps
	$(call golang-test-all,$@)

build: ./bin/sfncli
	go build -o bin/$(EXECUTABLE) $(PKG)

run: build
	bin/sfncli --activityname $(_DEPLOY_ENV)--$(_APP_NAME) \
	  --region us-west-2 \
		--cloudwatchregion us-west-1 \
	  --workername `hostname` \
	  --cmd ./bin/$(EXECUTABLE)


install_deps:
	go mod vendor
