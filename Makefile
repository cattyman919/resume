.PHONY: watch clean all

# Directory for output files
OUT_DIR = out
# Quote the name to handle spaces safely in variables
NAME = Seno Pamungkas Rahman - CV

# 1. Define your variants here
CV_TYPES = default frontend

# 2. 'all' triggers the specific targets below.
# This replaces the shell loop. Make will run 'default' and 'frontend'.
all: $(CV_TYPES)

# Ensure output directory exists
$(OUT_DIR):
    # On Windows, mkdir might error if it exists, so we suppress error with -
	-mkdir $(OUT_DIR)

# --- Build Targets ---

# 3. Use a "Pattern Rule" to build ANY type automatically.
# Usage: make default, make frontend, make devops
$(CV_TYPES): $(OUT_DIR)
	@echo "Building $@ CV..."
	typst compile src/main.typ "$(OUT_DIR)/$(NAME) ($@).pdf" --root . --input type=$@

# --- Utilities ---

watch-dir:
	-mkdir watch

# Watch mode. Usage: make watch TYPE=frontend
TYPE ?= default
watch: watch-dir
	typst watch src/main.typ "watch/$(NAME) ($(TYPE)).pdf" --root . --input type=$(TYPE)

clean:
    # Use a cross-platform delete if possible, or ignore errors
	rm -rf $(OUT_DIR) || rmdir /s /q $(OUT_DIR)
