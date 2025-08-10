# envvars-cli

A command-line interface tool for managing environment variables.

## Development Setup

This project uses a devcontainer for consistent development environments.

### Prerequisites

- [Docker](https://www.docker.com/products/docker-desktop/)
- [VS Code](https://code.visualstudio.com/) with the [Dev Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)

### Getting Started

1. Clone this repository
2. Open the project in VS Code
3. When prompted, click "Reopen in Container" or use the command palette (`Cmd/Ctrl + Shift + P`) and select "Dev Containers: Reopen in Container"
4. Wait for the container to build and start

### Available Commands

The project includes a Makefile with common development commands:

```bash
make help      # Show all available commands
make build     # Build the application
make run       # Run the application
make test      # Run tests
make lint      # Run linter
make format    # Format code
make deps      # Download dependencies
make tools     # Install development tools
```

### Go Tools Included

The devcontainer comes with the following Go development tools pre-installed:

- `goimports` - Code formatting and import organization
- `golint` - Code linting
- `golangci-lint` - Fast linter runner
- `dlv` - Go debugger
- `gomodifytags` - Struct tag manipulation
- `impl` - Interface implementation generator
- `gotests` - Test generation
- `gopkgs` - Package listing
- `gocov` - Code coverage

### Project Structure

```
envvars-cli/
├── .devcontainer/          # Devcontainer configuration
│   ├── devcontainer.json   # Main devcontainer config
│   ├── docker-compose.yml  # Docker services
│   └── Dockerfile         # Custom container image
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── Makefile               # Build and development commands
├── README.md              # This file
└── .gitignore            # Git ignore patterns
```

### Building and Running

```bash
# Build the application
go build -o bin/envvars-cli .

# Run the application
go run .

# Or use make commands
make build
make run
```

### Testing

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Or use make
make test
```

### Code Quality

```bash
# Format code
make format

# Run linter
make lint

# Clean build artifacts
make clean
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Submit a pull request

## License

This project is licensed under the terms specified in the LICENSE file.