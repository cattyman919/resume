# LaTeX Resume Generator

This repository contains a LaTeX-based resume generator that automatically builds multiple CV versions from a single `cv_data.yaml` file. It's designed to be easily customizable and extensible, allowing you to generate different CVs for various job applications without duplicating data.

## Features

- **Single Source of Truth**: All CV data is managed in a single `cv_data.yaml` file, making it easy to update and maintain.
- **Multiple CV Versions**: The generator can create different CV versions (e.g., `main`, `fullstack`, `devops`) by filtering experiences and projects based on tags.
- **Automated Build Process**: A `Makefile` automates the entire process of generating LaTeX files and compiling them into PDFs.
- **Customizable Sections**: The CV is structured into sections that can be easily modified or extended in the `generate_sections.go` file.
- **Modern LaTeX Template**: The resume template is clean, modern, and easy to read, with both a main version and a black-and-white version.

## Requirements

- [Go](https://golang.org/dl/) (1.18 or higher)
- [TeX Live](https://www.tug.org/texlive/) (or another LaTeX distribution)
- [Make](https://www.gnu.org/software/make/)

## Usage

To get started, you'll need to have Go and a LaTeX distribution (like TeX Live) installed on your system.

1.  **Clone the repository**:

    ```bash
    git clone https://github.com/cattyman919/resume.git
    cd resume
    ```

2.  **Install Go dependencies**:

    ```bash
    go mod tidy
    ```

3.  **Update your CV data**:
    Edit the `cv_data.yaml` file to include your personal information, experience, education, and other relevant details.

4.  **Build the CVs**:
    To build all CV versions, run the following command:
    ```bash
    make all
    ```
    This will generate the PDFs in the `out` directory.

## Customization

You can customize the CV template and generation process in the following ways:

- **Add a new CV type**: To add a new CV type (e.g., `data-science`), add it to the `TYPES` list in the `Makefile`. You can then tag experiences and projects in `cv_data.yaml` with the new type.
- **Modify the LaTeX template**: The LaTeX template is located in the `main_cv` and `bw_cv` directories. You can modify the `main.tex` files and the `preamble.sty` package to change the layout, fonts, and colors.
- **Add new sections**: To add a new section to the CV, you'll need to:
  1.  Add the corresponding data structure to `convert_latex.go`.
  2.  Create a new generator function in `generate_sections.go`.
  3.  Add the new section to the `main.tex` files.

## Contributing

Contributions are welcome! If you have any suggestions or improvements, please open an issue or submit a pull request.
