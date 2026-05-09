# Contributing

## Getting started

```bash
git clone https://github.com/lucasassuncao/devcontainerwizard.git
cd devcontainerwizard
go build ./...
go test ./...
```

Requires Go 1.21 or later.

## Workflow

1. Fork the repository
2. Create a branch: `git checkout -b feat/my-feature`
3. Commit your changes
4. Open a pull request against `main`

Write tests for new behaviour. Keep commits focused and well-described.

## Makefile targets

```bash
make build    # Build binary
make test     # Run tests
make lint     # Run linter
make install  # Install globally
make clean    # Remove build artifacts
```

## Generating demo GIFs

Demo GIFs are created with [VHS](https://github.com/charmbracelet/vhs) running inside Docker.

**Prerequisites:** Docker and a built `devcontainerwizard` binary.

```bash
# Build the binary
make build

# Build the Docker image
docker build -t devcontwiz:latest .

# Run VHS to generate the GIF
docker run --rm -v ${PWD}:/vhs devcontwiz:latest demo.tape
```

Edit `demo.tape` to change the recorded commands. VHS syntax reference: https://github.com/charmbracelet/vhs#vhs-command-reference
