GO ?= go

.PHONY: all
all:
	$(GO) mod tidy
	$(GO) build ./proxy
