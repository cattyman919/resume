## AGENTS.md

This document provides instructions for AI agents working in this repository.

### Build Commands

- To build both the main and black-and-white CVs:
  ```bash
  make all
  ```
- To generate only the LaTeX sections from `cv_data.yaml`:
  ```bash
  go run main.go generate_sections.go
  ```
- To build only the main CV:
  ```bash
  make main-cv
  ```
- To build only the devops CV:
  ```bash
  make devops-cv
  ```
- To build only the fullstack CV:
  ```bash
  make fullstack-cv
  ```
- To clean the output directory:
  ```bash
  make clean
  ```

### Code Style

- **Go**: The `main.go` and `generate_sections.go` scripts are used to generate LaTeX files from `cv_data.yaml`. Ensure any changes are compatible with Go.
- **LaTeX**: The main documents are `main_cv/main.tex` and `bw_cv/main.tex`. They use `subfile` to include sections from the `sections` directory. Custom commands and environments are defined in `packages/preamble.sty`.
- **Data**: All CV data is stored in `cv_data.yaml`. To update the CV, edit this file and then run `make`.

### Linting and Testing

There are no specific linting or testing frameworks configured in this repository. Ensure that the Go scripts run without errors and that the LaTeX files compile successfully using `pdflatex`.
