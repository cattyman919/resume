OS = linux
RM = rm -rf
TARGET := bin/cv_builder

.PHONY: all run clean

all: $(TARGET)
	@echo "Running application..."
	./$(TARGET)

$(TARGET): bin
	mkdir -p bin
	@echo "Building $(TARGET)..."
	CGO_ENABLED=0 GOOS=$(OS) go build -a -installsuffix cgo -o ./$(TARGET) ./cmd/resume/main.go

run:
	go run cmd/resume/main.go

clean:
	@echo "Cleaning up..."
	$(RM) out/*
	$(RM) cv/*
	$(RM) bin/*
