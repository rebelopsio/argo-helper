# Argo Helper

[![CI](https://github.com/rebelopsio/argo-helper/actions/workflows/ci.yml/badge.svg)](https://github.com/rebelopsio/argo-helper/actions/workflows/ci.yml)
[![Release](https://github.com/rebelopsio/argo-helper/actions/workflows/goreleaser.yml/badge.svg)](https://github.com/rebelopsio/argo-helper/actions/workflows/goreleaser.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/rebelopsio/argo-helper)](https://goreportcard.com/report/github.com/rebelopsio/argo-helper)
[![License](https://img.shields.io/github/license/rebelopsio/argo-helper)](LICENSE)

A CLI tool to help bootstrap and manage opinionated ArgoCD repository structures with best practices built in.

## Features

- Interactive TUI interface for easier navigation
- CLI commands for scripting and automation
- Initialize new ArgoCD repository structures
- Generate ApplicationSet manifests
- Example templates and best practices

## Installation

### Homebrew

```bash
brew install rebelopsio/tap/argo-helper
```

### Go Install

```bash
go install github.com/rebelopsio/argo-helper@latest
```

### Docker

```bash
docker run --rm -it rebelopsio/argo-helper
```

### Binary Downloads

Download pre-built binaries from the [Releases](https://github.com/rebelopsio/argo-helper/releases) page.

## Usage

### Interactive Mode

Simply run the command without arguments to start the interactive TUI:

```bash
argo-helper
```

### CLI Commands

#### Initialize a Repository

Create a new ArgoCD repository structure with an opinionated layout:

```bash
argo-helper init --project myproject [path]
```

Options:
- `--project, -p`: Name of the ArgoCD project (required)
- `--examples, -e`: Include example applications and ApplicationSet
- `--dry-run`: Preview the changes without making them

#### Create a New Resource

Create a new ArgoCD resource:

```bash
argo-helper new applicationset my-apps [--output path]
```

Options:
- `--output, -o`: Output path (default is templates/apps/)
- `--dry-run`: Preview the resource without creating it

## Directory Structure

When you initialize a repository, the following structure is created:

```
.
├── Chart.yaml                  # Helm chart metadata
├── README.md                   # Documentation
├── custom-resources/           # Custom Resource Definitions
├── templates/
│   ├── _helpers.tpl            # Common template helpers
│   ├── apps/                   # Application templates
│   └── projects/               # Project templates
├── values/                     # Environment-specific values
│   ├── dev/                    # Development values
│   └── prod/                   # Production values
└── values.yaml                 # Default values
```

## Development

### Requirements

- Go 1.24+
- [Task](https://taskfile.dev/) (optional, for using the Taskfile)

### Setup

Clone the repository:

```bash
git clone https://github.com/rebelopsio/argo-helper.git
cd argo-helper
```

Install dependencies:

```bash
go mod download
```

### Common Tasks

The project includes a Taskfile for common operations:

```bash
# Build the binary
task build

# Run tests
task test

# Format code
task fmt

# Run linters
task lint

# Run the application
task run

# Create a snapshot release
task release:snapshot
```

## Documentation

- [ApplicationSet Guide](docs/applicationset-guide.md) - Comprehensive guide to ApplicationSets
- Example templates are available in the `docs/examples/` directory

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.