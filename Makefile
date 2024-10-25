PWD := $(shell pwd)
OUTPUT_DIR := $(PWD)/web
SRC_DIR := $(PWD)/src
ARTICLE_DIR := $(PWD)/articles

export OUTPUT_DIR
export SRC_DIR
export ARTICLE_DIR

.PHONY: all
all: gopath
	mkdir -p $(OUTPUT_DIR)
	mkdir -p $(OUTPUT_DIR)/css
	cp -r $(SRC_DIR)/fonts $(OUTPUT_DIR)/.
	cp -r $(SRC_DIR)/images $(OUTPUT_DIR)/.
	cp -r $(SRC_DIR)/scripts $(OUTPUT_DIR)/.
	cp -r $(SRC_DIR)/pdfs $(OUTPUT_DIR)/.
	cp $(SRC_DIR)/robot.txt $(OUTPUT_DIR)/.
	$(GOPATH)/bin/templ generate -path $(SRC_DIR)
	go run $(SRC_DIR)/*.go

.PHONY: clean
clean:
	rm -rf $(SRC_DIR)/*templ.go
	rm -rf $(OUTPUT_DIR)/*

.PHONY: init
init: gopath
	ln -s $(PWD)/hooks/pre-commit $(PWD)/.git/hooks/pre-commit
	go install github.com/a-h/templ/cmd/templ@latest

.PHONY: deploy
deploy: clean init all

.PHONY: format
format: gopath
	gofmt -w $(SRC_DIR)/*.go
	$(GOPATH)/bin/templ fmt .

# Deploying via cloudflare pages
# GOPATH is not set by default, had to set it myself
# Gotta make sure it is able to find templ's binary
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
