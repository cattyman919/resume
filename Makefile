OS = linux
RM = rm -rf
TARGET := bin/cv_builder

.PHONY: all run run_debug docker clean

all:
	cargo run --bin autocv-cli

# Better for development (Local Build)
run:
	cargo run --bin autocv-cli

run_debug:
	cargo run --bin autocv-cli -- --debug

run_benchmark:
	cargo run --bin autocv-cli -- --benchmark

clean:
	@echo "Cleaning up..."
	$(RM) out
	$(RM) cv
	$(RM) bin
