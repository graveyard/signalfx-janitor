include golang.mk

.PHONY: all test build run
SHELL := /bin/bash
PKGS = $(shell go list ./... | grep -v "vendor/" )
$(eval $(call golang-version-check,1.8))

all: test build

test: $(PKGS)
$(PKGS): golang-test-all-deps
	$(call golang-test-all,$@)

build:
	go build

run: build
	./signalfx-janitor
