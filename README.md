# MSGGW (Message Gateway)

[![CI](https://github.com/leonkaihao/msggw/workflows/CI/badge.svg)](https://github.com/leonkaihao/msggw/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Version](https://img.shields.io/badge/version-1.1.0-blue.svg)](https://github.com/leonkaihao/msggw)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/leonkaihao/msggw)](https://goreportcard.com/report/github.com/leonkaihao/msggw)
[![codecov](https://codecov.io/gh/leonkaihao/msggw/branch/main/graph/badge.svg)](https://codecov.io/gh/leonkaihao/msggw)

A high-performance message gateway written in Go for filtering and transforming messages between different message brokers. MSGGW enables seamless message routing and transformation between multiple NATS brokers with powerful rule-based processing.

## Features

- **Multi-Broker Support**: Connect and route messages between any number of NATS brokers
- **Flexible Filtering**: Advanced filtering with ternary expressions (IS, NOT, MATCH operators)
- **Message Transformation**: Transform message topics, metadata, and payload formats
- **Branch Processing**: Conditional message routing with multiple processing branches
- **Built-in Functions**: Extensible function system for dynamic message transformations
- **Environment Variables**: Support for configuration properties and environment substitution
- **Multiple Payload Formats**: Support for msgbus, edgebus, and raw payload formats
- **Cloud Native**: Kubernetes-ready with Docker containerization support

## Quick Start

### Prerequisites

- Go 1.19 or higher
- NATS server(s)
- Docker (optional, for containerized deployment)
- Kubernetes (optional, for K8s deployment)

### Installation

#### Build from Source

```bash
# Clone the repository
git clone https://github.com/leonkaihao/msggw.git
cd msggw

# Build binary
make build

# Binary will be available in ./bin/msggw
```

#### Build Docker Image

```bash
make container
```

### Running MSGGW

#### Run Binary

```bash
./bin/msggw ./configs/server_side.json
```

#### Run with Docker

```bash
docker run --network host leonkaihao/msggw:1.1.0
```

#### Deploy to Kubernetes

```bash
kubectl apply -f ./deployment/msggw-configmap.yaml
kubectl apply -f ./deployment/msggw-deployment.yaml
```

## Configuration

### Configuration Structure

```json
{
  "logLevel": "info",
  "props": {
    "key": "value"
  },
  "brokers": [...],
  "flows": [...]
}
```

### Key Components

#### 1. Log Level
Set the logging verbosity: `info` or `debug`

#### 2. Properties
Define key-value pairs that can be referenced in flows using `{prop.key}` syntax.

#### 3. Brokers
Define multiple message brokers with unique names. You can configure any number of NATS brokers:

```json
[
  {
    "name": "broker-1",
    "type": "nats",
    "url": "nats://localhost:4222"
  },
  {
    "name": "broker-2",
    "type": "nats",
    "url": "nats://another-server:4222"
  }
]
```

**Note**: Broker names must be unique and are referenced by flows for message routing.

#### 4. Flows
Define message processing pipelines:

```json
{
  "name": "flow-name",
  "source": "broker-name",
  "subscribes": ["topic.>"],
  "payload": "msgbus",
  "branches": [...]
}
```

### Branch Configuration

Each branch contains filtering, transformation, and routing logic:

```json
{
  "name": "branch-name",
  "filters": [
    "{metadata.name} IS notification",
    "{metadata.pub_para} NOT NULL"
  ],
  "transforms": [
    {"{topic}": "/.notification.{func.PubParaToTopic}"}
  ],
  "sendTo": {
    "dest": "target-broker",
    "payload": "raw"
  }
}
```

### Filter Expressions

Filters use ternary expressions: `lval OPERATOR rval`

**Operators:**
- `IS` - Exact match
- `NOT` - Negation
- `MATCH` - Regular expression match

**Value Types:**
- `raw` - Raw string literals
- `{topic}` - Message topic
- `{metadata.field}` - Message metadata field
- `{func.FunctionName}` - Built-in function
- `{prop.key}` - Configuration property
- `NULL` - Empty/null value keyword
- `{/.mixed.{metadata.field}}` - Mixed expressions with sub-symbols

### Transformations

Transform message attributes before routing:

```json
{
  "{topic}": "new.topic.{metadata.field}",
  "{metadata.custom}": "{func.Timestamp}"
}
```

### Built-in Functions

- `PubParaToTopic` - Convert pub_para metadata to topic format
- `PubAddrToTopic` - Convert pub_addr metadata to topic format
- `Timestamp` - Generate timestamp

## Project Structure

```
.
├── Makefile                    # Build automation
├── README.md                   # This file
├── app
│   └── msggw                  # Application entry point
├── bin                        # Compiled binaries
├── build
│   └── docker                 # Dockerfile
├── component.json             # Version information
├── configs                    # Configuration examples
├── deployment                 # Kubernetes deployment files
├── go.mod                     # Go module definition
├── go.sum                     # Go dependencies checksum
└── pkg                        # Go packages
    ├── config                 # Configuration loader
    ├── funcs                  # Built-in functions
    ├── model                  # Core interfaces and constants
    ├── operator               # Expression operators (IS, NOT, MATCH)
    ├── parser                 # Expression parser
    ├── service                # Business logic processor
    └── symbol                 # Symbol definitions
```

## Example Configuration

See [configs/server_side.json](configs/server_side.json) for a complete example demonstrating:

- Bidirectional message routing between multiple NATS brokers
- Metadata-based filtering and routing
- Dynamic topic transformation
- Multiple payload format support

**Note**: The example uses `nats-internal` and `nats-public` as broker names, but you can configure any number of brokers with custom names to suit your architecture.

## Development

### Run Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# View coverage report
go tool cover -html=coverage.out
```

### Linting

```bash
# Install golangci-lint (if not already installed)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linters
golangci-lint run
```

### Build Commands

```bash
make build      # Build binary
make container  # Build Docker image
make clean      # Clean build artifacts
```

### Continuous Integration

This project uses GitHub Actions for CI/CD:

- **CI Workflow**: Runs on every push and pull request
  - Linting with golangci-lint
  - Unit tests with race detection
  - Code coverage reporting
  - Binary build
  - Docker image build

- **Release Workflow**: Triggers on version tags (e.g., `v1.1.0`)
  - Builds multi-platform binaries (Linux, macOS, Windows)
  - Creates GitHub release with changelog
  - Publishes Docker image to GitHub Container Registry

To create a new release:
```bash
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin v1.2.0
```

## Dependencies

- [msgbus](https://github.com/leonkaihao/msgbus) - Message bus abstraction library
- [nats.go](https://github.com/nats-io/nats.go) - NATS client for Go
- [logrus](https://github.com/sirupsen/logrus) - Structured logging

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

### Quick Contribution Guide

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure:
- All tests pass
- Code is properly formatted
- Linters pass without errors
- Documentation is updated if needed

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

**Leon KHao** - [leonkaihao](https://github.com/leonkaihao)

## Support

For issues and questions, please open an issue on the [GitHub repository](https://github.com/leonkaihao/msggw/issues).
