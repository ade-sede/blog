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
all: build
	go run $(SRC_DIR)/*.go

.PHONY: build
build: gopath
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
	if command -v templ >/dev/null 2>&1; then \
		templ generate -path $(SRC_DIR); \
	else \
		go run github.com/a-h/templ/cmd/templ@latest generate -path $(SRC_DIR); \
	fi

.PHONY: clean
clean:
	rm -rf $(SRC_DIR)/*templ.go
	find $(OUTPUT_DIR) -mindepth 1 -not -name "*.png" -exec rm -rf {} +

.PHONY: init
init:
	pre-commit install

.PHONY: deploy
deploy: clean all


.PHONY: pdf
pdf: build
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
	@if ! command -v entr > /dev/null; then \
		echo "Error: entr is not installed. Please install it (e.g., 'brew install entr' or 'sudo apt-get install entr')" && exit 1; \
	fi
	@echo "Starting server and watching for file changes..." >&2
	@echo "Serving on http://localhost:8080. Press Ctrl+C to exit." >&2
	@python3 -m http.server 8080 --directory $(OUTPUT_DIR) & \
	SERVER_PID=$$!; \
	trap "echo '\\nStopping server...'; kill $$SERVER_PID" INT TERM EXIT; \
	find $(SRC_DIR) $(ARTICLE_DIR) $(QUICK_NOTE_DIR) -type f | entr -c make all
