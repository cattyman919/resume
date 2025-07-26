OS = linux
RM = rm -rf
TARGET := bin/cv_builder

.PHONY: all run run_debug docker clean

all: $(TARGET)
	mkdir -p out
	@echo "Running application..."
	./$(TARGET)

$(TARGET):
	mkdir -p bin
	@echo "Building $(TARGET)..."
	CGO_ENABLED=0 GOOS=$(OS) go build -a -installsuffix cgo -o ./$(TARGET) ./cmd/resume/main.go

docker: $(TARGET)
	@echo "Building Docker image..."
	docker build -t resume-generator .
	@echo "Running Docker container..."
	docker run --rm \
		-u "$(shell id -u):$(shell id -g)" \
		-v "./out:/app/out" \
		-v "./$(TARGET):/app/$(TARGET)" \
		-v "./cv_data.yaml:/app/cv_data.yaml" \
		resume-generator

# Better for development (Local Build)
run:
	cargo run

run_debug:
	cargo run -- --debug

run_benchmark:
	cargo run -- --benchmark

clean:
	@echo "Cleaning up..."
	$(RM) out
	$(RM) cv
	$(RM) bin
