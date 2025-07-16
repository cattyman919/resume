# Utils
RM = rm -fr

.PHONY: all clean

all:
	go run cmd/resume/main.go

clean:
	$(RM) -rf out/*
	$(RM) -rf cv/*
