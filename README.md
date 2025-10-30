# Forgejo Classroom

An educational assignment management system that integrates with Forgejo to provide GitHub Classroom-like functionality for self-hosted Git platforms.

## Project Status

🚧 **This project is in early development** - Most functionality is not yet implemented. This repository contains the initial project structure and API framework.

## Features (Planned)

- **Classroom Management**: Create and manage coding classrooms
- **Assignment Distribution**: Create assignments from template repositories
- **Student Management**: Manage student rosters and account linking
- **Team Support**: Support for both individual and team assignments
- **Submission Tracking**: Automatic tracking of student submissions
- **Deadline Enforcement**: Automatic deadline management with git tags
- **CLI & API**: Complete CLI tool and REST API

## Architecture

This project follows a clean architecture with clear separation of concerns:

- **CLI Application** (`cmd/fgc/`) - Command-line interface for all operations
- **API Server** (`cmd/fgc-server/`) - REST API for web integrations
- **Service Layer** (`internal/service/`) - Business logic
- **Repository Layer** (`internal/repository/`) - Data access
- **Models** (`internal/model/`) - Domain objects
- **Forgejo Integration** (`internal/forgejo/`) - Forgejo API client

## Requirements

- Go 1.21+
- PostgreSQL 14+
- Redis 7+
- Docker & Docker Compose (for development)
- Forgejo instance with API access

## Quick Start

### 1. Clone and Setup

```bash
git clone <repository-url>
cd forgejo-classroom

# Copy configuration examples
cp config.yaml.example config.yaml
cp .env.example .env

# Edit configuration files with your settings
```

### 2. Development Environment

```bash
# Install dependencies
make deps

# Start development services (PostgreSQL + Redis)
make docker-dev-up

# Build the applications
make build

# Run tests
make test
```

### 3. CLI Usage

```bash
# Build and test CLI
make build-cli

# Show help
./bin/fgc --help

# Example commands (not yet implemented):
./bin/fgc classroom create "CS 101" --org="university-cs"
./bin/fgc assignment create "Homework 1" --classroom="cs-101" --template="https://your-forgejo.com/templates/hw1"
```

### 4. API Server

```bash
# Build and run API server
make build-server
./bin/fgc-server

# Or with Docker
make docker-build
docker-compose up
```

## Development

### Project Structure

```
forgejo-classroom/
├── cmd/                    # Application entry points
│   ├── fgc/               # CLI application
│   └── fgc-server/        # API server
├── internal/              # Private application code
│   ├── api/               # API handlers and routing
│   ├── service/           # Business logic
│   ├── repository/        # Data access layer
│   ├── model/             # Domain models
│   ├── forgejo/           # Forgejo integration
│   ├── cache/             # Caching layer
│   ├── config/            # Configuration
│   └── util/              # Utilities
├── pkg/                   # Public libraries
├── migrations/            # Database migrations
├── test/                  # Integration tests
└── docs/                  # Documentation
```

### Development Workflow

This project uses **Jujutsu (jj)** for version control. See `CLAUDE.md` for complete development guidelines.

```bash
# Create new change
jj new -m "Add feature description"

# Run tests before committing
make test

# Describe your change
jj describe -m "Detailed commit message"

# Push to remote
jj git push
```

### Available Make Targets

- `make build` - Build CLI and server
- `make test` - Run all tests
- `make test-integration` - Run integration tests
- `make lint` - Run linters
- `make docker-test-up` - Start test environment
- `make dev-setup` - Set up development environment
- `make clean` - Clean build artifacts

## Configuration

### Environment Variables

Key environment variables (see `.env.example`):

- `FGC_FORGEJO_BASE_URL` - Your Forgejo instance URL
- `FGC_FORGEJO_TOKEN` - Forgejo API token
- `FGC_DATABASE_*` - Database connection settings
- `FGC_REDIS_*` - Redis connection settings

### Configuration File

See `config.yaml.example` for full configuration options.

## API Documentation

API documentation is available at `/api/v1` when running the server. The complete OpenAPI specification is documented in `design.md`.

## Testing

```bash
# Unit tests
make test-unit

# Integration tests (requires Docker)
make test-integration

# All tests
make test

# Coverage report
make coverage
```

## Contributing

Please read `CONTRIBUTING.md` and `CLAUDE.md` for development guidelines.

1. Understand the architecture in `design.md`
2. Follow the development workflow in `CLAUDE.md`
3. Write tests for new functionality
4. Ensure all tests pass before submitting changes

## License

[License TBD]

## Roadmap

See `design.md` for detailed implementation phases and roadmap.

### Phase 1: Core MVP (In Progress)
- ✅ Project structure and build system
- ✅ CLI framework
- ✅ API framework
- 🔄 Database layer
- 🔄 Forgejo integration
- ⏳ Basic classroom management
- ⏳ Assignment creation

### Phase 2: Assignment Distribution
- ⏳ Template repository handling
- ⏳ Student repository creation
- ⏳ Assignment acceptance flow

### Phase 3: Advanced Features
- ⏳ Team support
- ⏳ Deadline enforcement
- ⏳ Submission tracking
- ⏳ Statistics and reporting

## Support

For questions and support:
- Check the documentation in `docs/`
- Review the technical design in `design.md`
- Create an issue for bugs or feature requests