# Contributing

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## Reporting Bugs

- Check existing [issues](https://github.com/Educentr/go-iproto/issues) to avoid duplicates.
- Open a new issue with a clear title and description.
- Include Go version, OS, and a minimal reproducing example.

## Suggesting Improvements

- Open an issue describing the improvement and its motivation.
- For significant changes, discuss the approach before submitting a PR.

## Development Setup

1. **Go 1.24+** — [https://go.dev/dl/](https://go.dev/dl/)
2. **golangci-lint** — [https://golangci-lint.run/welcome/install/](https://golangci-lint.run/welcome/install/)
3. **make** — for running targets

```bash
git clone https://github.com/Educentr/go-iproto.git
cd go-iproto
make test
make lint
```

## Code Style

- Format with `gofmt` and `goimports`.
- Run `make lint` before committing.
- Follow existing code conventions in the project.

## Pull Request Process

1. Fork the repository and create a feature branch.
2. Make your changes with clear commit messages.
3. Ensure `make test` and `make lint` pass.
4. Submit a PR against the `main` branch.
5. Describe what the PR does and link any related issues.

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](https://github.com/Educentr/go-iproto/blob/main/LICENSE).
