PWD := $(shell pwd)
OUTPUT_DIR := $(PWD)/web
TEMPL_DIR := $(PWD)/templ

export OUTPUT_DIR

# Deploying via cloudflare pages
# GOPATH is not set by default, had to set it myself
# Gotta make sure it is able to find templ's binary
.PHONY: all
all: gopath
	$(GOPATH)/bin/templ generate -path $(TEMPL_DIR)
	go run $(TEMPL_DIR)/*.go

.PHONY: clean
clean:
	rm -f $(OUTPUT_DIR)/*.html
	rm -f $(TEMPL_DIR)/*templ.go

.PHONY: init
init: gopath
	ln -s $(PWD)/hooks/pre-commit $(PWD)/.git/hooks/pre-commit
	go install github.com/a-h/templ/cmd/templ@latest

.PHONY: deploy
deploy: clean init all

.PHONY: format
format: gopath
	gofmt -w $(TEMPL_DIR)/*.go
	$(GOPATH)/bin/templ fmt .

.PHONY: gopath
gopath:
ifeq ($(GOPATH),)
	@echo "Error: GOPATH is not set" >&2
	@exit 1
else
	@echo "GOPATH is $(GOPATH)"
endif

.PHONY: re
re: clean all
