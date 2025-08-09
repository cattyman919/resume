FROM texlive/texlive:latest

ARG GO_VERSION=1.24.6

ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

# Install make and the necessary TeX Live packages.
RUN apt-get update && apt-get install -y --no-install-recommends \
    wget \
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

CMD ["go", "run", "cmd/resume/main.go"]
