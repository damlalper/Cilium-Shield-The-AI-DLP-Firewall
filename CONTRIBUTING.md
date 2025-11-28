# Contributing to Cilium-Shield

Thank you for your interest in contributing to Cilium-Shield! This document provides guidelines for contributing to the project.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-username/Cilium-Shield.git`
3. Create a branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Test your changes
6. Commit your changes: `git commit -m "Add: your feature description"`
7. Push to your fork: `git push origin feature/your-feature-name`
8. Open a Pull Request

## Development Setup

```bash
# Install dependencies
make install

# Start development environment
make dev
```

## Code Style

### Go Code (Wasm Filter & Control Plane)

- Follow standard Go formatting: `gofmt -s -w .`
- Use meaningful variable names
- Add comments for complex logic
- Write unit tests for new features

### JavaScript/TypeScript (Backend & Frontend)

- Use 2 spaces for indentation
- Use semicolons
- Add JSDoc comments for functions
- Follow React best practices for components

## Testing

Before submitting a PR, ensure all tests pass:

```bash
# Run Go tests
make test-go

# Run integration tests
make test

# Test local dev environment
make dev
# Then manually test the dashboard
```

## Commit Messages

Use clear, descriptive commit messages:

- `Add: New feature description`
- `Fix: Bug description`
- `Update: Change description`
- `Refactor: Code improvement description`
- `Docs: Documentation update description`

## Pull Request Guidelines

1. **Title**: Clear, concise description of changes
2. **Description**:
   - What changes were made?
   - Why were these changes needed?
   - How were they tested?
3. **Tests**: Include tests for new features
4. **Documentation**: Update relevant docs
5. **Small PRs**: Keep PRs focused on a single feature or fix

## Areas for Contribution

### High Priority

- [ ] Additional regex patterns for PII detection
- [ ] Performance benchmarks and optimizations
- [ ] Integration with other observability tools
- [ ] Support for additional protocols (gRPC, WebSocket)
- [ ] Production-grade configuration management

### Medium Priority

- [ ] Advanced filtering rules (allow/deny lists)
- [ ] Custom redaction policies per namespace
- [ ] Metrics and Prometheus integration
- [ ] Helm chart for easy deployment
- [ ] E2E tests

### Documentation

- [ ] More examples and use cases
- [ ] Video tutorials
- [ ] Blog posts and articles
- [ ] Translation to other languages

## Questions?

- Open an issue for bugs or feature requests
- Check existing issues before creating new ones
- Join discussions in existing issues

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Help others learn and grow
- Focus on what's best for the community

## License

By contributing to Cilium-Shield, you agree that your contributions will be licensed under the same license as the project (see LICENSE file).

## Recognition

Contributors will be recognized in:
- README.md (Contributors section)
- Release notes
- Project documentation

Thank you for contributing to Cilium-Shield! 🚀
