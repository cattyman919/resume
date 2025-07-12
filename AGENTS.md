## AGENTS.md

This document provides instructions for AI agents working in this repository.

### Build Commands

- To build both the main and black-and-white CVs:
  ```bash
  make all
  ```
- To generate only the LaTeX sections from `cv_data.json`:
  ```bash
  make generate
  ```
- To build only the main CV:
  ```bash
  make main_cv
  ```
- To build only the black-and-white CV:
  ```bash
  make bw_cv
  ```
- To clean the output directory:
  ```bash
  make clean
  ```

### Code Style

- **Python**: The `generate_sections.py` script is used to generate LaTeX files from `cv_data.json`. It uses standard Python libraries and f-strings for templating. Ensure any changes are compatible with Python 3.
- **LaTeX**: The main documents are `main_cv/main.tex` and `bw_cv/main.tex`. They use `subfile` to include sections from the `sections` directory. Custom commands and environments are defined in `packages/preamble.sty`.
- **Data**: All CV data is stored in `cv_data.json`. To update the CV, edit this file and then run `make`.

### Linting and Testing

There are no specific linting or testing frameworks configured in this repository. Ensure that the `generate_sections.py` script runs without errors and that the LaTeX files compile successfully using `pdflatex`.
