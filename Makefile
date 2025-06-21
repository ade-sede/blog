PWD := $(shell pwd)
OUTPUT_DIR := $(PWD)/web
SRC_DIR := $(PWD)/src
ARTICLE_DIR := $(PWD)/articles
QUICK_NOTE_DIR := $(PWD)/quick-notes

export OUTPUT_DIR
export SRC_DIR
export ARTICLE_DIR
export QUICK_NOTE_DIR

.PHONY: all
all: prepare
	go run $(SRC_DIR)/*.go

.PHONY: prepare
prepare: gopath
	mkdir -p $(OUTPUT_DIR)
	mkdir -p $(OUTPUT_DIR)/css
	mkdir -p $(OUTPUT_DIR)/libs
	cp -r $(SRC_DIR)/fonts $(OUTPUT_DIR)/.
	cp -r $(SRC_DIR)/webfonts $(OUTPUT_DIR)/.
	cp -r $(SRC_DIR)/images $(OUTPUT_DIR)/.
	cp -r $(SRC_DIR)/scripts $(OUTPUT_DIR)/.
	cp -r $(SRC_DIR)/pdfs $(OUTPUT_DIR)/.
	cp -r $(SRC_DIR)/libs/katex $(OUTPUT_DIR)/libs/.
	cp -r $(SRC_DIR)/libs/color-thief $(OUTPUT_DIR)/libs/.
	find $(QUICK_NOTE_DIR) -name "*.css" -type f -exec cp {} $(OUTPUT_DIR)/css/. \; 2>/dev/null || true
	find $(QUICK_NOTE_DIR) -name "*.js" -type f -exec cp {} $(OUTPUT_DIR)/scripts/. \; 2>/dev/null || true
	find $(ARTICLE_DIR) -name "*.css" -type f -exec cp {} $(OUTPUT_DIR)/css/. \; 2>/dev/null || true
	find $(ARTICLE_DIR) -name "*.js" -type f -exec cp {} $(OUTPUT_DIR)/scripts/. \; 2>/dev/null || true
	cp $(SRC_DIR)/robots.txt $(OUTPUT_DIR)/.
	$(GOPATH)/bin/templ generate -path $(SRC_DIR)

.PHONY: clean
clean:
	rm -rf $(SRC_DIR)/*templ.go
	rm -rf $(OUTPUT_DIR)/*

.PHONY: init
init: gopath
	go install github.com/a-h/templ/cmd/templ@latest

.PHONY: init-dev
init-dev: init
	pre-commit install

.PHONY: deploy
deploy: clean init all


.PHONY: pdf
pdf: prepare
	go install github.com/chromedp/chromedp
	go run $(SRC_DIR)/*.go -pdf

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

.PHONY: serve
serve: all
	python3 -m http.server 8080 --directory $(OUTPUT_DIR)
