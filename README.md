# AutoCV

This repository hosts AutoCV, a powerful CV generation tool that leverages Go and LaTeX to produce multiple, tailored resumes from a single YAML data source. It streamlines the process of creating and maintaining professional CVs, making it easy to customize them for different job applications.

## Key Features

- **Single Source of Truth**: Manage all your CV data in a `config/` directory with `general.yaml`, `experiences.yaml`, and `projects.yaml` files, ensuring consistency and simplifying updates.
- **Dynamic CV Versions**: Generate different CV versions (e.g., `main`, `fullstack`, `devops`) by tagging experiences and projects.
- **Concurrent Builds**: Utilizes Go's concurrency features to build multiple CVs in parallel, significantly reducing generation time.
- **LaTeX Templating**: Leverages Go's `text/template` package to create clean and maintainable LaTeX code.
- **Dockerized Environment**: Comes with a Dockerfile for a consistent and reproducible build process.

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) (1.18 or higher)
- [TeX Live](https://www.tug.org/texlive/) (or another LaTeX distribution)
- [Make](https://www.gnu.org/software/make/)
- [Docker](https://www.docker.com/get-started) (Optional)

### Local Build

1.  **Clone the repository**:

    ```bash
    git clone https://github.com/senohebat/resume.git
    cd resume
    ```

2.  **Install Go dependencies**:

    ```bash
    go mod tidy
    ```

3.  **Customize your CV data**:
    Edit the YAML files in the `config/` directory (`general.yaml`, `experiences.yaml`, and `projects.yaml`) to add your personal information, work experience, projects, and other details.

4.  **Build the CVs**:
    To build all CV versions, run:
    ```bash
    make all
    ```
    The generated PDFs will be available in the `out` directory.

### Docker Build

For a containerized build without needing a local LaTeX installation (Go is still required):

1.  **Clone and set up the repository** as shown above.

2.  **Install Go dependencies**:

    ```bash
    go mod tidy
    ```

3.  **Build with Docker**:
    ```bash
    make docker
    ```
    The PDFs will be generated in the `out` directory.

## Customization

- **Add a new CV type**: To create a new CV version (e.g., `data-science`), simply add the new tag to the relevant items in your YAML files.
- **Modify LaTeX templates**: The LaTeX templates are in `template_cv/main_cv` and `template_cv/bw_cv`. You can edit the `.tex` files and the `preamble.sty` to change the styling.
- **Extend the data model**: To add new sections, update the Go structs in `internals/types.go`, create a corresponding generator function in `internals/generate_sections.go`, and include the new section in the `main.tex` files.

## Contributing

Contributions are welcome! If you have ideas for improvements, please open an issue or submit a pull request.
