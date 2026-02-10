FROM debian:bullseye-slim

# Avoid interactive prompts during installation
ENV DEBIAN_FRONTEND=noninteractive

ARG GO_VERSION=1.24.6

ENV GOROOT=/usr/local/go
ENV GOPATH=/app/go-modules
ENV GOCACHE=/app/.go-build-cache
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

# Install Go and LaTeX dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    wget \
    ca-certificates \
    texlive-latex-extra \
    texlive-fonts-recommended \
    texlive-pictures \
    texlive-fonts-extra \
    lmodern \
    && wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz \
    && rm go${GO_VERSION}.linux-amd64.tar.gz \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

ENTRYPOINT ["go", "run", "cmd/resume/main.go"]
