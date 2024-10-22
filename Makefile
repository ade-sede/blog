PWD := $(shell pwd)
OUTPUT_DIR := $(PWD)/web
TEMPL_DIR := $(PWD)/templ

export OUTPUT_DIR

# Deploying via cloudflare pages
# GOPATH is not set by default, had to set it myself
# Gotta make sure it is able to find templ's binary
.PHONY: all
all:
	$(GOPATH)/bin/templ generate -path $(TEMPL_DIR)
	go run $(TEMPL_DIR)/*.go

.PHONY: clean
clean:
	rm -f $(OUTPUT_DIR)/*.html
	rm -f $(TEMPL_DIR)/*templ.go

.PHONY: deps
deps:
	go install github.com/a-h/templ/cmd/templ@latest

.PHONY: deploy
deploy: clean deps all

.PHONY: re
re: clean all
