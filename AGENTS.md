## AGENTS.md

This document provides instructions for AI agents working in this repository.

### Build Commands

- To build all CVs:
  ```bash
  make all
  ```
- To clean the output directory:
  ```bash
  make clean
  ```

### Code Style

- **Go**: The `cmd/resume/main.go` script is used to generate and compile LaTeX files from `cv_data.yaml`. Ensure any changes are compatible with Go.
- **LaTeX**: The main documents are `template_cv/main_cv/main.tex` and `template_cv/bw_cv/main.tex`. They use `subfile` to include sections from the `cv` directory. Custom commands and environments are defined in `packages/preamble.sty`.
- **Data**: All CV data is stored in `cv_data.yaml`. To update the CV, edit this file and then run `make`.

### Linting and Testing

There are no specific linting or testing frameworks configured in this repository. Ensure that the Go scripts run without errors and that the LaTeX files compile successfully using `pdflatex`.
