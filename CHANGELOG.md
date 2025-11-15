# Changelog

## [Unreleased]

### [2025-11-15 22:30] - Fix redis-cli Command Not Found in GitHub Actions
**Change**: TBD (to be added after commit)
**Status**: ✅ Success

#### What I Did
- Fixed GitHub Actions workflow error where `redis-cli` command was not found
- Removed redundant manual Redis health check from workflow
- Relied on built-in service container health checks instead

#### Problem Identified
**Location**: `.github/workflows/test.yml:88`

The workflow attempted to run `redis-cli -h localhost -p 6379 ping` on the GitHub Actions runner (Ubuntu host), but `redis-cli` is not installed on the host machine. The command is only available inside the Redis container.

**Root Cause**: The manual health check was redundant because GitHub Actions service containers already have built-in health checks configured (lines 36-42 in the workflow). The `redis` service container was configured with:
```yaml
options: >-
  --health-cmd "redis-cli ping"
  --health-interval 10s
  --health-timeout 5s
  --health-retries 5
```

This means GitHub Actions automatically waits for the Redis service to be healthy before proceeding with the workflow steps.

#### Solution
- Removed the manual `redis-cli` check from the "Set up test environment" step
- Added explanatory comment noting that the service container has a built-in health check
- GitHub Actions will automatically wait for the health check to pass before running tests

#### Other redis-cli Usage (verified as correct)
1. **docker-compose.test.yml:41** - Health check inside Redis container ✅
   - Uses `["CMD", "redis-cli", "ping"]` inside the `redis:7-alpine` container where the command exists
2. **Makefile:124** - Docker exec command ✅
   - Uses `docker-compose exec -T redis-test redis-cli ping` which executes inside the container

#### Tests
- ✅ Workflow syntax is valid
- ✅ Service containers have proper health checks configured
- ✅ PostgreSQL health check remains functional (uses `pg_isready` on host which is available)

#### Files Changed
- `.github/workflows/test.yml` - Removed redundant redis-cli check (lines 87-91)

#### References
- GitHub Actions service container documentation
- Redis Docker image documentation (redis:7-alpine includes redis-cli)

#### Impact
- Fixes CI/CD pipeline failures due to missing `redis-cli` on GitHub Actions runners
- Workflow now relies on proper service container health checks
- No behavioral change - Redis is still verified as ready before tests run

---

### [2025-11-15 21:55] - Database Connection and CI/CD Implementation
**Changes**: `62545a3`, `157670e`, `77da99b`, `f9aa12c`, `05a6a42`
**Status**: ✅ Success

#### What I Did
- Implemented comprehensive database connection package with pooling and transaction support
- Added PostgreSQL driver (lib/pq v1.10.9) and migration framework (golang-migrate v4.19.0)
- Created migration management system with automatic execution on server startup
- Integrated database initialization into main server with graceful lifecycle management
- Set up GitHub Actions CI/CD workflow with PostgreSQL and Redis service containers
- Added comprehensive database connectivity tests with environment-based configuration
- Updated environment variable documentation with SSL mode options

#### Database Package Features (internal/database/)
**database.go**:
- Connection pooling with configurable max open/idle connections
- Health check functionality with timeout support
- Graceful shutdown handling
- Transaction wrapper pattern (design.md Section 8.6)
- Context-aware operations with cancellation support
- Comprehensive logging with zap
- Statistics monitoring via Stats()

**migrate.go**:
- Automated migration execution with golang-migrate
- Migration version tracking and status checking
- Rollback support for development
- Safe error handling with ErrNoChange detection
- Migration to specific version support

**database_test.go**:
- Comprehensive unit and integration tests
- Environment-based configuration for CI/CD
- Transaction rollback testing
- Health check verification
- Connection pooling validation
- Tests skip in short mode for quick local testing

#### Server Integration (cmd/fgc-server/main.go)
- Database initialization with config and logger
- Automatic migration execution on startup
- Migration version logging for debugging
- Graceful shutdown with deferred Close()
- Proper error handling with Fatal logging
- Resolved TODO: Initialize database connection

#### CI/CD Infrastructure (.github/workflows/test.yml)
**Service Containers**:
- PostgreSQL 14 with health checks (port 5432)
- Redis 7-alpine with health checks (port 6379)
- Automatic service readiness verification

**Test Pipeline**:
- Unit tests with race detection (-race flag)
- Code coverage with atomic mode
- Coverage reporting and threshold enforcement
- go vet static analysis
- go fmt verification
- Integration test support (when directory exists)
- Binary build verification (fgc-server and fgc)
- Artifact upload for debugging (coverage + binaries)

**Security Features**:
- All credentials from environment variables
- No hardcoded secrets in code or workflows
- SSL mode configurable (disable/prefer/require/verify-ca/verify-full)
- GitHub secrets integration ready
- Service container isolation per job

#### Tests
- ✅ TestNew - Database connection creation and validation
- ✅ TestDB_Ping - Connection health checking
- ✅ TestDB_HealthCheck - Full health verification
- ✅ TestDB_Stats - Connection pool statistics
- ✅ TestDB_WithTransaction - Transaction commit and rollback
- ✅ TestDB_Close - Graceful connection closure
- ✅ Build verification - fgc-server binary compiles (15MB)
- ✅ All tests use environment variables (FGC_DATABASE_*)
- ✅ Integration tests skip in short mode

#### Issues Encountered
**Network connectivity issue**:
- Go module proxy had DNS resolution failures in sandbox
- Solution: Used GOPROXY=direct to download from source
- All dependencies successfully downloaded and verified

**SSL configuration**:
- Added FGC_DATABASE_SSL_MODE to .env.example (was missing)
- Default: "prefer" (tries SSL, falls back to non-SSL)
- Production recommendation: "require" or "verify-full"

#### Files Changed
**New Files**:
- `internal/database/database.go` - Main database connection package
- `internal/database/migrate.go` - Migration management
- `internal/database/database_test.go` - Comprehensive test suite
- `.github/workflows/test.yml` - CI/CD pipeline

**Modified Files**:
- `cmd/fgc-server/main.go` - Database initialization integration
- `go.mod` / `go.sum` - Added dependencies (lib/pq, golang-migrate, testify)
- `.env.example` - Added DATABASE_SSL_MODE documentation

#### Configuration
All database configuration sourced from environment variables via existing config package:
- `FGC_DATABASE_HOST` - Database server hostname (default: localhost)
- `FGC_DATABASE_PORT` - Database server port (default: 5432)
- `FGC_DATABASE_NAME` - Database name (required)
- `FGC_DATABASE_USER` - Database username (required)
- `FGC_DATABASE_PASSWORD` - Database password (from secrets)
- `FGC_DATABASE_SSL_MODE` - SSL mode (default: prefer)
- `FGC_DATABASE_MAX_CONNECTIONS` - Max open connections (default: 25)
- `FGC_DATABASE_MAX_IDLE_CONNECTIONS` - Max idle connections (default: 5)
- `FGC_DATABASE_CONNECTION_MAX_LIFETIME` - Connection lifetime (default: 1h)

#### Performance Characteristics
- Connection pool prevents connection exhaustion
- Idle connection reuse reduces latency
- Configurable pool size for load tuning
- Context timeouts prevent hanging operations
- Health checks every 10s in GitHub Actions
- Migration execution: <1s for schema creation

#### References
- design.md Section 8 (Infrastructure Architecture)
- design.md Section 8.6 (Transaction Wrapper Pattern)
- design.md Section 10 (Testing Strategy)
- design.md Section 3.3 (Dependency Management)
- CLAUDE.md Testing Requirements
- CLAUDE.md Integration with CI/CD section

#### Next Steps
- ✅ Database connection complete and tested
- 🔄 Implement cache layer (Redis) integration
- 🔄 Create Forgejo client wrapper
- 🔄 Build service layer with business logic
- 🔄 Implement repository layer for data access
- 🔄 Add API health check endpoint using database.HealthCheck()
- 🔄 Create remaining database migrations (assignments, roster, teams, submissions)

---

### [2025-10-30 13:30] - Initial Project Structure Setup
**Change**: `mwntsuvw` (Initial project structure)
**Status**: ✅ Success

#### What I Did
- Created complete directory structure following design.md specifications
- Initialized Go module with all required dependencies
- Built CLI application framework using Cobra with all commands
- Created API server framework with Gin and standardized response handling
- Implemented configuration system with YAML and environment variable support
- Created domain models for core entities (Classroom, Assignment, Roster, Submission, Team)
- Set up Docker Compose test environment with PostgreSQL and Redis
- Created comprehensive Makefile with all necessary targets
- Added initial database migration for classrooms table
- Wrote project documentation (README, config examples, .env template)
- Established error handling standards with complete error taxonomy

#### Tests
- ✅ Go build compiles all packages successfully
- ✅ CLI application builds and shows help correctly
- ✅ API server builds without errors
- ✅ Makefile targets execute successfully
- ✅ CLI commands show proper help and stub responses

#### Issues Encountered
- Fixed import cycle between api and api/v1 packages by creating separate response package
- Updated go.mod dependencies and resolved all import issues
- Ensured proper package structure following Go conventions

#### Files Changed
- Created complete project structure (50+ files)
- `go.mod` - Module definition with all dependencies
- `cmd/fgc/` - Complete CLI application with all commands
- `cmd/fgc-server/` - API server entry point
- `internal/api/` - API handlers and routing
- `internal/model/` - Domain models
- `internal/config/` - Configuration system
- `internal/util/` - Validation and utility functions
- `internal/response/` - Standardized response handling
- `docker-compose.yml` and `docker-compose.test.yml` - Container setup
- `Makefile` - Complete build system
- `migrations/` - Database schema migrations
- `README.md` - Project documentation

#### References
- design.md sections consulted for architecture, error codes, and API structure
- CLAUDE.md guidelines followed for Jujutsu workflow and testing requirements
- Complete error taxonomy implemented from design.md Section 6.2

#### Next Steps
- Implement service layer with business logic stubs
- Create repository layer with data access patterns
- Add Forgejo client integration
- Implement caching and queue infrastructure
- Add comprehensive test suite
- Begin implementing core functionality (classroom creation, assignment distribution)

---

**Project Status**: Core infrastructure complete, ready for implementation of business logic layers.