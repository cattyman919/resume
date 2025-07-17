# Stage 1: Final image with a stable and reliable TeX Live distribution
FROM debian:bullseye-slim

# Avoid interactive prompts during installation
ENV DEBIAN_FRONTEND=noninteractive

# Install make and the necessary TeX Live packages.
RUN apt-get update && apt-get install -y --no-install-recommends \
    make \
    texlive-latex-extra \
    texlive-fonts-recommended \
    texlive-pictures \
    texlive-fonts-extra \
    lmodern \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the pre-compiled Go binary
COPY bin/cv_builder /app/bin/cv_builder

# Copy the project files needed to generate the CVs
COPY Makefile /app/
COPY template_cv/ /app/template_cv/
COPY images/ /app/images/

# Set the default command to run the resume generation
CMD ["make", "all"]
