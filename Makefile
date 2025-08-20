OS = linux
RM = rm -rf
TARGET := Resume
IMAGE_NAME := resume-generator

.PHONY: all run run_debug docker clean build

all:
	@echo "Running application..."
	@ninja --quiet -C build
	@./build/bin/$(TARGET)

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

clean:
	@echo "Cleaning up..."
	$(RM) out
	$(RM) cv
	$(RM) bin
