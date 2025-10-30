# CLAUDE.md

## Project Context

This is **Forgejo Classroom**, a Go-based educational assignment management system that integrates with Forgejo. The complete technical design, architecture, and implementation strategy is documented in `design.md` in the project root.

**Critical**: Always read and reference `design.md` before making any architectural or design decisions. It contains:
- Complete API specifications (OpenAPI)
- Database schemas
- Security requirements
- Error handling standards
- Testing strategies
- Performance targets

## Version Control: Jujutsu (jj)

This project uses **Jujutsu (jj)** instead of Git. Key differences and commands:

### Essential jj Commands

```bash
# Check status
jj status

# Create new change (like git commit, but creates change immediately)
jj new

# Describe current change (add commit message)
jj describe -m "Your message here"

# Show diff
jj diff

# Show log
jj log

# Rebase/move changes
jj rebase -d <destination>

# Squash changes together
jj squash

# Split a change
jj split

# Undo last operation
jj undo

# Push to remote
jj git push

# Pull from remote
jj git fetch
jj rebase -d main@origin

# Abandon a change (like git reset)
jj abandon
```

### jj Workflow for This Project

1. **Before starting work**: `jj new -m "Brief description of what you're about to do"`
2. **As you work**: Files are automatically tracked in current change
3. **Test and verify**: Run tests before describing the change
4. **If tests pass**: `jj describe -m "Detailed commit message"`
5. **If tests fail**: Continue working in same change OR `jj undo` to revert
6. **Ready to push**: `jj git push`

### Key jj Concepts

- **Changes are created upfront**: Unlike git, you create a change before making edits
- **Automatic tracking**: Modified files are automatically included
- **Changes are immutable**: Each operation creates a new change; old ones are preserved
- **No staging area**: All modified files are in the current change
- **Easy to undo**: `jj undo` reverts the last operation
- **Conflicts are explicit**: jj materializes conflicts in the working copy

Reference: https://jj-vcs.github.io/jj/latest/

## Development Workflow

### 1. Before Writing Code

**ALWAYS:**
1. Read the relevant section of `design.md`
2. Understand the architecture layer you're working in (API → Service → Repository → Database)
3. Check error code taxonomy in `design.md` Section 6.2
4. Review the OpenAPI spec in `design.md` Section 14 if working on API

### 2. Creating a Change

```bash
# Start new work
jj new -m "impl: Add classroom service layer with caching"

# Or for bug fixes
jj new -m "fix: Correct deadline validation in assignment service"

# Or for tests
jj new -m "test: Add unit tests for roster bulk operations"
```

### 3. Development Cycle

**For EVERY compile-ready change:**

```bash
# 1. Write code following design.md architecture

# 2. Run tests IMMEDIATELY when code compiles
make test

# 3. If tests PASS:
jj describe -m "Detailed commit message explaining what was done"
# Update CHANGELOG.md (see below)

# 4. If tests FAIL:
# - Fix the code in the same change
# - Keep working until tests pass
# - Do NOT create a new change for fixes

# 5. Only when tests pass AND changelog is updated:
jj new -m "Next task description"
```

## Testing Requirements

### Test Execution

**MANDATORY**: Run tests before describing any change.

```bash
# Run all tests
make test

# Run specific test package
go test ./internal/service/... -v

# Run with coverage
go test ./... -cover -coverprofile=coverage.out

# Run integration tests (requires docker-compose)
make test-integration

# Run contract tests
make test-contract
```

### Docker Compose Test Environment

The project uses Docker Compose for integration testing:

```bash
# Start test environment
docker-compose -f docker-compose.test.yml up -d

# Run integration tests
make test-integration

# View logs
docker-compose -f docker-compose.test.yml logs -f

# Clean up
docker-compose -f docker-compose.test.yml down -v
```

**docker-compose.test.yml** should include:
- PostgreSQL (test database)
- Redis (test cache)
- Mock Forgejo server (if available)

### Test Structure Requirements

Reference `design.md` Section 10 for complete testing strategy.

#### Unit Tests (60% of tests)

```go
// internal/service/assignment_test.go
package service

import (
    "testing"
    "github.com/stretchr/testify/suite"
)

type AssignmentServiceTestSuite struct {
    suite.Suite
    service     *AssignmentService
    mockRepo    *MockAssignmentRepo
    mockForgejo *MockForgejoClient
    mockCache   *MockCache
}

func (s *AssignmentServiceTestSuite) SetupTest() {
    // Create mocks
    s.mockRepo = new(MockAssignmentRepo)
    s.mockForgejo = new(MockForgejoClient)
    s.mockCache = new(MockCache)

    s.service = &AssignmentService{
        repo:          s.mockRepo,
        forgejoClient: s.mockForgejo,
        cache:         s.mockCache,
    }
}

func (s *AssignmentServiceTestSuite) TestCreate_Success() {
    // Setup mocks
    s.mockRepo.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
    s.mockForgejo.On("GetRepository", mock.Anything, "org/repo").Return(&Repository{ID: 100}, nil)

    // Execute
    result, err := s.service.Create(ctx, req)

    // Assert
    s.NoError(err)
    s.NotNil(result)

    // Verify mocks called correctly
    s.mockRepo.AssertExpectations(s.T())
}

func TestAssignmentServiceSuite(t *testing.T) {
    suite.Run(t, new(AssignmentServiceTestSuite))
}
```

**Requirements:**
- Use testify/suite for organization
- Mock all external dependencies
- Test both success and error cases
- Verify mock expectations
- Target 80%+ coverage for business logic

#### Integration Tests (30% of tests)

```go
// test/api/assignment_test.go
package api_test

import (
    "testing"
    "github.com/stretchr/testify/suite"
)

type AssignmentAPITestSuite struct {
    suite.Suite
    server *httptest.Server
    db     *sql.DB
    cache  *redis.Client
    client *client.Client
}

func (s *AssignmentAPITestSuite) SetupSuite() {
    s.db = setupTestDB()
    s.cache = setupTestRedis()
    s.server = startTestServer(s.db, s.cache)
    s.client = client.New(s.server.URL, testToken)
}

func (s *AssignmentAPITestSuite) SetupTest() {
    cleanDatabase(s.db)
    s.cache.FlushAll(context.Background())
}
```

**Requirements:**
- Test against real database (in Docker)
- Clean state before each test
- Verify database state after operations
- Test caching behavior

#### Contract Tests (10% of tests)

Verify API responses match OpenAPI specification (see `design.md` Section 14).

## Code Quality Standards

### Adherence to design.md

- **API endpoints**: Must match OpenAPI spec exactly (Section 14)
- **Error codes**: Must use taxonomy from Section 6.2
- **Resource models**: Must match schemas in Section 5
- **Architecture layers**: Must follow API → Service → Repository pattern (Section 3)
- **Caching**: Must follow TTL strategy in Section 8.2
- **Security**: Must implement permission checks from Section 11.1

### Code Style

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Check for common issues
go vet ./...
```

### Transaction Pattern

Always use the transaction wrapper from `design.md` Section 8.6:

```go
err := s.repo.WithTransaction(ctx, func(tx *sql.Tx) error {
    // All database operations here
    return nil
})
```

### Error Handling

Use standardized error codes from `design.md` Section 6.2:

```go
// Good
return nil, fmt.Errorf("BUSINESS_DEADLINE_PASSED: %w", err)

// Bad
return nil, errors.New("deadline passed")
```

### Context Cancellation

Always check context cancellation in loops:

```go
for _, item := range items {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }

    // Process item
}
```

## CHANGELOG.md Format

Maintain an **atomic changelog** that documents ALL activities, not just commit-worthy changes. This is your development log.

### Structure

```markdown
# Changelog

## [Unreleased]

### [2025-10-29 14:30] - Classroom Service Implementation
**Change**: `<jj change id>`
**Status**: ✅ Success / ❌ Failed / 🔄 In Progress

#### What I Did
- Implemented `ClassroomService.Create()` method
- Added cache invalidation logic
- Created unit tests for Create operation

#### Tests
- ✅ TestClassroomService_Create_Success (PASS)
- ✅ TestClassroomService_Create_DuplicateOrg (PASS)
- ❌ TestClassroomService_Create_ForgejoError (FAIL - need to fix mock setup)

#### Issues Encountered
- Mock setup for Forgejo client was incorrect
- Fixed by properly setting up mock expectations
- Cache invalidation pattern needed adjustment per design.md Section 8.2

#### Files Changed
- `internal/service/classroom.go`
- `internal/service/classroom_test.go`

#### Next Steps
- Fix failing test for Forgejo error case
- Add integration test for classroom creation

---

### [2025-10-29 12:15] - Database Schema Migration
**Change**: `<jj change id>`
**Status**: ✅ Success

#### What I Did
- Created migration 000001_create_classrooms
- Added indexes per design.md recommendations
- Verified migration up/down

#### Tests
- ✅ Migration applies successfully
- ✅ Migration rollback works correctly
- ✅ Indexes created as expected

#### Files Changed
- `migrations/000001_create_classrooms.up.sql`
- `migrations/000001_create_classrooms.down.sql`

---
```

### Changelog Entry Template

After EVERY change where tests pass:

```markdown
### [YYYY-MM-DD HH:MM] - Brief Description
**Change**: `<jj change id from 'jj log'>`
**Status**: ✅ Success / ❌ Failed / 🔄 In Progress

#### What I Did
- Bullet points of actual work done
- Both successful implementations and failed attempts
- Design decisions and why they were made

#### Tests
- ✅/❌ Test names with pass/fail status
- Coverage percentage if relevant
- Any tests skipped and why

#### Issues Encountered
- Problems faced and how they were solved
- Bugs discovered and fixed
- Performance issues and optimizations

#### Files Changed
- List of modified files
- New files created
- Files deleted

#### References
- design.md sections consulted
- OpenAPI endpoints implemented
- Related changes

#### Next Steps (if incomplete)
- What needs to be done next
- Open questions
- Blockers
```

## Change Atomicity Rules

### What Makes a Change Atomic

A change is atomic when:
1. ✅ Code compiles without errors
2. ✅ All tests pass
3. ✅ Implements a single logical unit of work
4. ✅ Can be safely reverted without breaking other changes
5. ✅ Changelog is updated

### Change Granularity

**Good atomic changes:**
- "Add classroom service Create method with unit tests"
- "Implement cache invalidation for assignment updates"
- "Add validation for deadline field in assignment creation"

**Too large (split into multiple changes):**
- "Implement entire classroom management feature"
- "Add all API endpoints for assignments"

**Too small (combine):**
- "Add import statement"
- "Fix typo in comment"

## Common Development Patterns

### Starting a New Feature

```bash
# 1. Read design.md relevant section
cat design.md | grep -A 50 "Assignment Distribution"

# 2. Create change for feature
jj new -m "feat: Implement assignment creation API endpoint"

# 3. Start with tests (TDD)
# Create test file first: internal/service/assignment_test.go

# 4. Implement code to make tests pass

# 5. Run tests
make test

# 6. If pass, describe and update changelog
jj describe -m "feat: Implement assignment creation with validation

- Add AssignmentService.Create() method
- Implement template repository validation
- Add deadline enforcement logic
- Cache invalidation for classroom assignments list
- Tests: 8 passing, 80% coverage

Refs: design.md Section 9.4"

# Update CHANGELOG.md with entry

# 7. Move to next task
jj new -m "feat: Add assignment creation API handler"
```

### Fixing a Test Failure

```bash
# DO NOT create a new change for fixes
# Stay in current change and fix

# 1. Run specific failing test
go test ./internal/service -run TestAssignmentService_Create_InvalidDeadline -v

# 2. Fix the code

# 3. Run test again
go test ./internal/service -run TestAssignmentService_Create_InvalidDeadline -v

# 4. Once passing, run full test suite
make test

# 5. Only then describe
jj describe -m "feat: Assignment creation with deadline validation

- Implemented Create method
- Fixed deadline validation logic (was allowing past dates)
- All tests passing

Tests:
- TestAssignmentService_Create_Success ✅
- TestAssignmentService_Create_InvalidDeadline ✅ (fixed)
- TestAssignmentService_Create_TemplateNotFound ✅"
```

### Refactoring

```bash
jj new -m "refactor: Extract cache key generation to separate package"

# Make changes
# Run ALL tests to ensure no regressions
make test

# Verify no behavioral changes
make test-integration

jj describe -m "refactor: Extract cache key generation

- Created internal/cache/keys.go
- Centralized all cache key patterns
- No behavioral changes
- All tests passing (100% coverage maintained)"
```

## Integration with CI/CD

The CI pipeline will:
1. Run `make test` on every push
2. Check test coverage (must be >80% for service layer)
3. Run linters (`golangci-lint`, `go vet`)
4. Verify contract tests against OpenAPI spec
5. Run integration tests in Docker environment

Ensure your local tests pass before pushing:

```bash
# Pre-push checklist
make test              # All tests pass
make test-integration  # Integration tests pass
make lint             # No lint errors
jj log -r @          # Review your change
jj describe          # Ensure good commit message

# Push
jj git push
```

## Emergency: Undo Changes

```bash
# Undo last operation (safe, reversible)
jj undo

# Abandon current change (discard all work)
jj abandon

# Restore abandoned change
jj undo

# View operation log
jj op log

# Restore to specific operation
jj op restore <operation-id>
```

## Key Reminders

### ✅ DO
- Read design.md before implementing features
- Run tests after every compile-ready change
- Use jj commands properly (not git)
- Update CHANGELOG.md with every change
- Mock external dependencies in unit tests
- Follow error code taxonomy
- Use transaction wrapper pattern
- Check context cancellation in loops
- Write tests first (TDD encouraged)

### ❌ DON'T
- Create changes without running tests
- Use git commands (use jj instead)
- Skip changelog updates
- Implement features not in design.md without discussion
- Create custom error codes (use taxonomy)
- Access database directly (use repository layer)
- Ignore test failures (fix in same change)
- Create changes for minor fixes (stay in current change)
- Push without running full test suite locally

## Performance Targets

Reference `design.md` Section 9 for complete targets. Key metrics:

- API endpoints: <200ms p95 latency
- Repository creation: <2s p95 latency
- List operations: <500ms for 50 students
- Test suite: <30s for unit tests, <2min for integration
- Cache hit rate: >80%

Run benchmarks:

```bash
go test -bench=. -benchmem ./internal/service/...
```

## Getting Help

1. **Design questions**: Check `design.md` first
2. **API contracts**: See `design.md` Section 14 (OpenAPI spec)
3. **Error codes**: See `design.md` Section 6.2
4. **Testing patterns**: See `design.md` Section 10
5. **jj commands**: https://jj-vcs.github.io/jj/latest/

## Quick Reference Card

```bash
# Start work
jj new -m "Brief description"

# Check status
jj status

# Run tests (ALWAYS before describing)
make test

# Tests pass? Describe change
jj describe -m "Detailed message"

# Update CHANGELOG.md
# (Add entry with change ID, status, details)

# Tests fail? Keep working in same change
# Don't create new change until tests pass

# Ready for next task
jj new -m "Next task"

# Push (after tests pass)
jj git push

# Made a mistake?
jj undo
```

---

**Remember**: Tests → Describe → Changelog → New Change. Every change must be atomic and tested.
