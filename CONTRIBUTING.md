# Contributing to MSGGW

Thank you for your interest in contributing to MSGGW! This document provides guidelines and instructions for contributing.

## Development Setup

### Prerequisites

- Go 1.19 or higher
- Git
- Make
- Docker (optional)

### Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/msggw.git
   cd msggw
   ```

3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/leonkaihao/msggw.git
   ```

4. Install dependencies:
   ```bash
   go mod download
   ```

## Development Workflow

### Creating a Branch

```bash
git checkout -b feature/your-feature-name
```

Branch naming conventions:
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation updates
- `refactor/` - Code refactoring
- `test/` - Test additions or modifications

### Making Changes

1. Make your changes in your feature branch
2. Write or update tests as needed
3. Ensure your code follows Go conventions
4. Run tests locally:
   ```bash
   go test ./...
   ```

5. Run linters:
   ```bash
   golangci-lint run
   ```

### Code Style

- Follow standard Go conventions and idioms
- Use `gofmt` to format your code
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions small and focused

### Commit Messages

Write clear and meaningful commit messages:

```
<type>: <subject>

<body>

<footer>
```

Types:
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, missing semi colons, etc)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Example:
```
feat: add support for MQTT broker type

Implement MQTT broker support to enable message routing
between NATS and MQTT brokers.

Closes #123
```

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# View coverage report
go tool cover -html=coverage.out
```

### Writing Tests

- Write table-driven tests when possible
- Test edge cases and error conditions
- Aim for high code coverage (>80%)
- Use descriptive test names that explain what is being tested

## Submitting Changes

### Before Submitting

- [ ] Tests pass locally
- [ ] Code is properly formatted
- [ ] Linters pass without errors
- [ ] Documentation is updated if needed
- [ ] Commit messages follow the convention

### Creating a Pull Request

1. Push your changes to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

2. Go to the [MSGGW repository](https://github.com/leonkaihao/msggw) on GitHub
3. Click "New Pull Request"
4. Select your feature branch
5. Fill in the PR template with:
   - Clear description of the changes
   - Related issue numbers
   - Testing performed
   - Screenshots (if applicable)

6. Submit the pull request

### Pull Request Review

- Maintainers will review your PR
- Address any feedback or requested changes
- Once approved, your PR will be merged

## Code Review Guidelines

When reviewing code:
- Be respectful and constructive
- Focus on the code, not the person
- Explain the reasoning behind suggestions
- Approve when the code meets quality standards

## Issue Reporting

### Bug Reports

When reporting bugs, include:
- Clear description of the issue
- Steps to reproduce
- Expected vs actual behavior
- Environment details (OS, Go version, etc.)
- Relevant logs or error messages

### Feature Requests

When requesting features, include:
- Clear description of the feature
- Use case and motivation
- Proposed implementation (if any)
- Examples of similar features in other projects

## Questions?

If you have questions about contributing:
- Open an issue on GitHub
- Check existing issues and discussions
- Review the project documentation

Thank you for contributing to MSGGW!
