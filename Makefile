OS = linux
RM = rm -rf
TARGET := bin/cv_builder
IMAGE_NAME := resume-generator

.PHONY: all run run_debug docker clean build

all: run

build:
	mkdir -p bin
	@echo "Building $(TARGET)..."
	CGO_ENABLED=0 GOOS=$(OS) go build -a -installsuffix cgo -o ./$(TARGET) ./cmd/resume/main.go

docker:
	@# It checks if the output of 'docker images -q' is empty (-z).
	@if [ -z "$(shell docker images -q $(IMAGE_NAME))" ]; then \
		echo "Image '$(IMAGE_NAME)' not found. Building..."; \
		docker build -t $(IMAGE_NAME) .; \
	else \
		echo "Image '$(IMAGE_NAME)' found. Skipping build."; \
	fi
	@echo "Running Docker container..."
	docker run --rm \
		-u "$(shell id -u):$(shell id -g)" \
		-v "./:/app" \
		$(IMAGE_NAME)

# if you have pdflatex installed, you can use this target to generate a PDF
# This process is much faster than using the Docker container
# but requires a manually installed TexLive on your system
run:
	@odin run . --collection:deps=deps

clean:
	@echo "Cleaning up..."
	@chmod -R u+w go-modules .go-build-cache || true
	$(RM) out
	$(RM) cv
	$(RM) bin
	$(RM) go-modules
	$(RM) .go-build-cache
