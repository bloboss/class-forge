# Forgejo Classroom: Technical Implementation Proposal (Revised)

**Version**: 2.0
**Date**: October 29, 2025
**Status**: Design Review - Revised
**Authors**: Forgejo Classroom Team

---

## Document Revision History

| Version | Date | Changes | Reviewer |
|---------|------|---------|----------|
| 1.0 | Oct 27, 2024 | Initial draft | - |
| 2.0 | Oct 29, 2025 | API design improvements, standardization, caching strategy | Architecture Review |

---

## **Table of Contents**

1. [Executive Summary](#1-executive-summary)
2. [Language Choice Rationale](#2-language-choice-rationale)
3. [Project Structure](#3-project-structure)
4. [API Design Principles](#4-api-design-principles)
5. [Resource Models and Endpoints](#5-resource-models-and-endpoints)
6. [Error Handling Standards](#6-error-handling-standards)
7. [API Versioning Strategy](#7-api-versioning-strategy)
8. [Infrastructure Architecture](#8-infrastructure-architecture)
9. [Implementation Phases](#9-implementation-phases)
10. [Testing Strategy](#10-testing-strategy)
11. [Security Considerations](#11-security-considerations)
12. [Deployment Guide](#12-deployment-guide)
13. [Monitoring and Observability](#13-monitoring-and-observability)
14. [API Documentation Standards](#14-api-documentation-standards)
15. [Migration Strategy](#15-migration-strategy)
16. [Future Roadmap](#16-future-roadmap)
17. [Success Metrics](#17-success-metrics)
18. [Risk Assessment and Mitigation](#18-risk-assessment-and-mitigation)

---

## **1. Executive Summary**

Forgejo Classroom is an educational assignment management system that integrates with Forgejo to provide GitHub Classroom-like functionality for self-hosted Git platforms. This document outlines the technical architecture, API design, and implementation strategy for building a production-ready MVP.

**Key Design Decisions:**

- **API-First Architecture**: Complete REST API before web interface
- **Language**: Go for ecosystem alignment and development velocity
- **CLI-Driven MVP**: Validates all functionality programmatically
- **RESTful Resource Modeling**: Consistent, predictable API design
- **Comprehensive Error Taxonomy**: Standardized error codes and handling
- **Layered Architecture**: API → Service → Repository → Database
- **Phased Implementation**: 14-week path to alpha release

**Core Value Proposition:**

Teachers can distribute assignments to students, each receiving their own repository based on a template, with automatic deadline enforcement via Git tags.

---

## **2. Language Choice Rationale**

### **2.1 Why Go?**

Go is selected as the implementation language for the following reasons:

#### **Ecosystem Alignment**
- **Forgejo itself is written in Go** - Critical for understanding internals and potential code sharing
- **Mature Forgejo Go SDK** (`code.gitea.io/sdk/gitea`) - Proven, well-maintained client library
- **Community consistency** - Contributors familiar with Forgejo can contribute to Classroom

#### **Development Velocity**
- Fast compilation (~1-2 seconds for full project)
- Single static binary output (no dependency management for users)
- Rich standard library for web services
- Extensive ecosystem of CLI and API frameworks

#### **Deployment Simplicity**
```bash
# Cross-platform builds from single machine
GOOS=linux GOARCH=amd64 go build -o fgc-linux
GOOS=darwin GOARCH=arm64 go build -o fgc-macos
GOOS=windows GOARCH=amd64 go build -o fgc.exe
```

#### **Maintainability**
- Low learning curve for new contributors
- Clear, readable code style enforced by `gofmt`
- Strong typing with practical flexibility
- Educational projects benefit from broad contributor accessibility

#### **Performance Profile**
```
Go Web Service Benchmarks:
- Throughput: 50,000-100,000 req/s (target: 100 req/s)
- Latency p99: <10ms (target: <100ms)
- Memory: 20-50MB base + ~10KB per goroutine
- Cold start: <100ms
```

**Sufficient for this project** - Bottlenecks will be Forgejo API calls and database, not language performance.

#### **Alternatives Considered and Rejected**

**Rust:**
- ✅ Superior type safety and performance
- ❌ No official Forgejo SDK (would need custom implementation)
- ❌ Steep learning curve reduces contributor pool
- ❌ 2x longer development time
- ❌ Slower compilation impacts iteration speed
- **Verdict**: Viable but unnecessary; performance not a constraint

**C/C++:**
- ✅ Maximum theoretical performance
- ❌ Memory safety risks unacceptable for educational software
- ❌ No modern web service ecosystem
- ❌ 4-5x longer development time
- ❌ Complex cross-platform deployment
- **Verdict**: Not recommended for web services

**Decision**: Go provides the optimal balance of development velocity, ecosystem fit, and maintainability for this project.

---

## **3. Project Structure**

### **3.1 Repository Organization**

```
forgejo-classroom/
├── cmd/
│   ├── fgc/                          # CLI application
│   │   ├── main.go
│   │   └── commands/
│   │       ├── classroom.go          # Classroom commands
│   │       ├── assignment.go         # Assignment commands
│   │       ├── roster.go             # Roster commands
│   │       ├── submission.go         # Submission commands
│   │       ├── team.go               # Team commands
│   │       └── student.go            # Student commands
│   │
│   └── fgc-server/                   # API server
│       └── main.go
│
├── internal/
│   ├── api/                          # API handlers
│   │   ├── v1/
│   │   │   ├── classroom.go
│   │   │   ├── assignment.go
│   │   │   ├── roster.go
│   │   │   ├── submission.go
│   │   │   ├── team.go
│   │   │   └── middleware.go
│   │   ├── router.go                 # Route registration
│   │   ├── response.go               # Response helpers
│   │   └── errors.go                 # Error definitions
│   │
│   ├── service/                      # Business logic
│   │   ├── classroom.go
│   │   ├── assignment.go
│   │   ├── roster.go
│   │   ├── submission.go
│   │   ├── team.go
│   │   └── deadline.go               # Deadline enforcement
│   │
│   ├── repository/                   # Data access layer
│   │   ├── classroom.go
│   │   ├── assignment.go
│   │   ├── roster.go
│   │   ├── submission.go
│   │   ├── team.go
│   │   └── transaction.go
│   │
│   ├── model/                        # Domain models
│   │   ├── classroom.go
│   │   ├── assignment.go
│   │   ├── roster.go
│   │   ├── submission.go
│   │   └── team.go
│   │
│   ├── forgejo/                      # Forgejo integration
│   │   ├── client.go                 # Forgejo API client
│   │   ├── organization.go           # Organization operations
│   │   ├── repository.go             # Repository operations
│   │   ├── team.go                   # Team operations
│   │   └── user.go                   # User operations
│   │
│   ├── git/                          # Git operations
│   │   ├── clone.go
│   │   ├── tag.go                    # Deadline tagging
│   │   └── archive.go                # Repository archiving
│   │
│   ├── cache/                        # Caching layer (NEW)
│   │   ├── cache.go                  # Cache interface
│   │   ├── redis.go                  # Redis implementation
│   │   ├── memory.go                 # In-memory fallback
│   │   └── keys.go                   # Cache key patterns
│   │
│   ├── queue/                        # Async operations (NEW)
│   │   ├── queue.go                  # Queue interface
│   │   ├── redis.go                  # Redis-backed queue
│   │   ├── worker.go                 # Worker pool
│   │   └── jobs/                     # Job definitions
│   │       ├── deadline.go           # Deadline enforcement job
│   │       └── bulk.go               # Bulk operation jobs
│   │
│   ├── config/                       # Configuration
│   │   ├── config.go
│   │   └── defaults.go
│   │
│   ├── auth/                         # Authentication
│   │   ├── token.go
│   │   └── permission.go             # Permission checks
│   │
│   └── util/                         # Utilities
│       ├── slug.go                   # Slug generation
│       ├── validator.go              # Input validation
│       └── time.go                   # Time utilities
│
├── pkg/                              # Public libraries
│   └── client/                       # Go client library
│       ├── client.go
│       ├── classroom.go
│       ├── assignment.go
│       └── types.go
│
├── migrations/                       # Database migrations
│   ├── 000001_create_classrooms.up.sql
│   ├── 000001_create_classrooms.down.sql
│   ├── 000002_create_roster.up.sql
│   ├── 000002_create_roster.down.sql
│   └── ...
│
├── test/                             # Integration tests
│   ├── api/                          # API integration tests
│   ├── cli/                          # CLI integration tests
│   ├── contract/                     # API contract tests (NEW)
│   └── fixtures/                     # Test data
│
├── docs/                             # Documentation
│   ├── api/                          # API documentation
│   │   └── openapi.yaml
│   ├── cli/                          # CLI usage guide
│   ├── deployment/                   # Deployment guides
│   └── development/                  # Developer guides
│
├── scripts/                          # Build and utility scripts
│   ├── build.sh
│   ├── test.sh
│   ├── migrate.sh
│   └── dev-setup.sh
│
├── .github/
│   └── workflows/                    # CI/CD workflows
│       ├── test.yml
│       ├── build.yml
│       └── release.yml
│
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
├── docker-compose.yml                # Development environment
├── README.md
├── LICENSE
└── CONTRIBUTING.md
```

### **3.2 Rationale for Structure**

**Separation of Concerns:**
- **`cmd/`**: Entry points for executables (CLI and server)
- **`internal/`**: Private application code (not importable by external projects)
- **`pkg/`**: Public libraries that can be imported by other projects
- **Clear layering**: API → Service → Repository → Database

**Benefits:**
1. **Testability**: Each layer can be mocked and tested independently
2. **Maintainability**: Changes to one layer don't cascade to others
3. **Extensibility**: Easy to add new features without breaking existing code
4. **Collaboration**: Multiple developers can work on different layers simultaneously

### **3.3 Dependency Management**

**Core Dependencies:**
```go
// go.mod
module code.forgejo.org/forgejo/classroom

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1          // HTTP framework
    github.com/lib/pq v1.10.9                // PostgreSQL driver
    github.com/spf13/cobra v1.8.0            // CLI framework
    github.com/spf13/viper v1.18.2           // Configuration
    code.gitea.io/sdk/gitea v0.17.1          // Gitea/Forgejo SDK
    github.com/golang-migrate/migrate/v4 v4.17.0  // Database migrations
    github.com/go-redis/redis/v8 v8.11.5     // Redis client (NEW)
    gopkg.in/yaml.v3 v3.0.1                  // YAML parsing
    github.com/stretchr/testify v1.8.4       // Testing
    go.uber.org/zap v1.26.0                  // Structured logging
)
```

**Rationale for Choices:**

1. **Gin** (HTTP Framework): Fast, lightweight, excellent middleware support
2. **Cobra** (CLI Framework): Industry standard, automatic help generation
3. **Viper** (Configuration): Multi-format support, environment variable binding
4. **golang-migrate** (Migrations): Database-agnostic, version tracking
5. **Gitea SDK** (Forgejo Client): Forgejo is Gitea-compatible, well-maintained
6. **Redis** (Caching/Queue): Simple deployment, high performance

---

## **4. API Design Principles**

### **4.1 RESTful Resource Modeling**

**Core Principles:**

1. **Resources, not actions** - URLs represent resources, HTTP methods represent actions
2. **Hierarchical relationships** - Nest resources when child cannot exist without parent
3. **Consistent naming** - Plural nouns for collections, singular for documents
4. **Idempotent operations** - Same request produces same result
5. **Stateless interactions** - Each request contains all necessary information

### **4.2 Resource Nesting Guidelines**

**Nest when the child cannot exist without the parent:**
```
✅ POST /classrooms/{classroom_id}/assignments
✅ GET  /classrooms/{classroom_id}/assignments/{assignment_id}
```

**Use top-level with filters when resources are independently addressable:**
```
✅ GET /submissions?classroom_id={id}
✅ GET /submissions?assignment_id={id}
✅ GET /submissions?student_id={id}
```

**Maximum nesting depth: 2 levels**
```
❌ /classrooms/{id}/assignments/{id}/submissions/{id}/comments
✅ /submissions/{id}/comments
```

### **4.3 URL Structure Standards**

```
Base URL: https://forgejo.example.com/api/classroom/v1

Pattern:
/{resource}                          # List collection
/{resource}/{id}                     # Get specific resource
/{resource}/{id}/{sub-resource}      # Access nested resource
```

**Examples:**
```
GET    /classrooms                   # List classrooms
POST   /classrooms                   # Create classroom
GET    /classrooms/{id}              # Get classroom details
PATCH  /classrooms/{id}              # Update classroom
DELETE /classrooms/{id}              # Delete classroom

GET    /classrooms/{id}/assignments  # List assignments in classroom
POST   /classrooms/{id}/assignments  # Create assignment
```

### **4.4 HTTP Method Usage**

| Method | Purpose | Idempotent | Safe | Request Body | Response Body |
|--------|---------|------------|------|--------------|---------------|
| GET | Retrieve resource(s) | ✅ | ✅ | ❌ | ✅ |
| POST | Create resource | ❌ | ❌ | ✅ | ✅ |
| PUT | Replace resource | ✅ | ❌ | ✅ | ✅ |
| PATCH | Update resource | ❌ | ❌ | ✅ | ✅ |
| DELETE | Remove resource | ✅ | ❌ | ❌ | ❌/✅ |

### **4.5 Pagination Standards**

**All list endpoints MUST support pagination:**

```yaml
parameters:
  - name: page
    in: query
    schema:
      type: integer
      minimum: 1
      default: 1
  - name: per_page
    in: query
    schema:
      type: integer
      minimum: 1
      maximum: 100
      default: 30
```

**Response includes pagination metadata:**
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "per_page": 30,
    "total_count": 150,
    "total_pages": 5
  }
}
```

**Response headers:**
```http
X-Total-Count: 150
X-Page: 1
X-Per-Page: 30
Link: <...?page=2>; rel="next", <...?page=5>; rel="last"
```

### **4.6 Filtering and Sorting Standards**

**Common query parameters (reusable across endpoints):**

```yaml
components:
  parameters:
    StatusFilter:
      name: status
      in: query
      schema:
        type: string

    SortField:
      name: sort
      in: query
      description: Field to sort by
      schema:
        type: string
        enum: [created_at, updated_at, name, deadline]
        default: created_at

    SortOrder:
      name: order
      in: query
      description: Sort direction
      schema:
        type: string
        enum: [asc, desc]
        default: desc
```

**Example usage:**
```
GET /assignments?status=active&sort=deadline&order=asc
```

### **4.7 Bulk Operations Pattern**

**Synchronous (small operations, <100 items):**

```yaml
POST /classrooms/{classroom_id}/roster/bulk
Content-Type: application/json

Request:
{
  "operations": [
    {"action": "add", "identifier": "student001", "email": "s1@edu"},
    {"action": "add", "identifier": "student002", "email": "s2@edu"},
    {"action": "remove", "identifier": "student003"}
  ]
}

Response: 200 OK
{
  "results": [
    {"index": 0, "status": "success", "id": 123},
    {"index": 1, "status": "error", "error": {"code": "DUPLICATE_IDENTIFIER"}},
    {"index": 2, "status": "success", "id": null}
  ],
  "summary": {
    "total": 3,
    "succeeded": 2,
    "failed": 1
  }
}
```

**Asynchronous (large operations, 100+ items):**

```yaml
POST /classrooms/{classroom_id}/roster/import
Content-Type: multipart/form-data

Response: 202 Accepted
{
  "job_id": "import_abc123",
  "status": "processing",
  "status_url": "/jobs/import_abc123"
}

# Check status
GET /jobs/{job_id}
Response: 200 OK
{
  "id": "import_abc123",
  "status": "completed",
  "progress": {
    "total": 500,
    "processed": 500,
    "succeeded": 485,
    "failed": 15
  },
  "errors_url": "/jobs/import_abc123/errors",
  "created_at": "2025-10-29T10:00:00Z",
  "completed_at": "2025-10-29T10:02:34Z"
}
```

---

## **5. Resource Models and Endpoints**

### **5.1 Core Resources**

#### **5.1.1 Classroom**

**Resource Model:**
```go
type Classroom struct {
    ID               int       `json:"id"`
    Name             string    `json:"name"`
    Slug             string    `json:"slug"`
    Description      string    `json:"description,omitempty"`
    OrganizationID   int       `json:"organization_id"`
    OrganizationName string    `json:"organization_name"`
    OwnerID          int       `json:"owner_id"`
    Status           string    `json:"status"` // active, archived
    StudentCount     int       `json:"student_count"`
    AssignmentCount  int       `json:"assignment_count"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
}
```

**Endpoints:**

```
GET    /classrooms
POST   /classrooms
GET    /classrooms/{id}
PATCH  /classrooms/{id}
DELETE /classrooms/{id}
POST   /classrooms/{id}/archive
```

#### **5.1.2 Assignment**

**Resource Model:**
```go
type Assignment struct {
    ID                   int       `json:"id"`
    ClassroomID          int       `json:"classroom_id"`
    Title                string    `json:"title"`
    Slug                 string    `json:"slug"`
    Instructions         string    `json:"instructions,omitempty"`
    Type                 string    `json:"type"` // individual, team
    TemplateRepoID       int       `json:"template_repo_id"`
    TemplateRepoName     string    `json:"template_repo_name"`
    Deadline             time.Time `json:"deadline,omitempty"`
    MaxTeamSize          int       `json:"max_team_size,omitempty"`
    AllowLateSubmissions bool      `json:"allow_late_submissions"`
    Visibility           string    `json:"visibility"` // private, public
    InvitationCode       string    `json:"invitation_code"`
    InvitationURL        string    `json:"invitation_url"`
    AcceptanceCount      int       `json:"acceptance_count"`
    SubmissionCount      int       `json:"submission_count"`
    CreatedAt            time.Time `json:"created_at"`
    UpdatedAt            time.Time `json:"updated_at"`
}
```

**Endpoints:**

```
GET    /classrooms/{classroom_id}/assignments
POST   /classrooms/{classroom_id}/assignments
GET    /assignments/{id}
PATCH  /assignments/{id}
DELETE /assignments/{id}
GET    /assignments/{id}/stats
```

#### **5.1.3 Submission (NEW - Explicitly Defined)**

**Resource Model:**
```go
type Submission struct {
    ID             int       `json:"id"`
    AssignmentID   int       `json:"assignment_id"`
    StudentID      int       `json:"student_id"`
    RepositoryID   int       `json:"repository_id"`
    RepositoryName string    `json:"repository_name"`
    RepositoryURL  string    `json:"repository_url"`
    Status         string    `json:"status"` // pending, in_progress, submitted, graded
    AcceptedAt     time.Time `json:"accepted_at"`
    LastCommitAt   time.Time `json:"last_commit_at,omitempty"`
    LastCommitSHA  string    `json:"last_commit_sha,omitempty"`
    SubmittedAt    time.Time `json:"submitted_at,omitempty"`
    DeadlineTag    string    `json:"deadline_tag,omitempty"`
    IsLate         bool      `json:"is_late"`
    CommitCount    int       `json:"commit_count"`
    Grade          float64   `json:"grade,omitempty"`
    GradedAt       time.Time `json:"graded_at,omitempty"`
    GradedBy       int       `json:"graded_by,omitempty"`
    Feedback       string    `json:"feedback,omitempty"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}
```

**Endpoints:**

```
# Create submission (student accepts assignment)
POST   /assignments/{assignment_id}/submissions

# List submissions (with flexible filtering)
GET    /submissions?assignment_id={id}
GET    /submissions?student_id={id}
GET    /submissions?classroom_id={id}
GET    /submissions?status={status}

# Individual submission operations
GET    /submissions/{id}
PATCH  /submissions/{id}
GET    /submissions/{id}/commits
POST   /submissions/{id}/submit
```

**Key Design Decision:** Submissions are top-level resources because they can be queried across multiple dimensions (by assignment, student, classroom, status). This provides flexibility for different use cases while maintaining clean URL structure.

#### **5.1.4 Roster Entry**

**Resource Model:**
```go
type RosterEntry struct {
    ID              int       `json:"id"`
    ClassroomID     int       `json:"classroom_id"`
    Identifier      string    `json:"identifier"` // Student ID from institution
    Email           string    `json:"email"`
    FullName        string    `json:"full_name"`
    ForgejoUserID   int       `json:"forgejo_user_id,omitempty"`
    ForgejoUsername string    `json:"forgejo_username,omitempty"`
    Status          string    `json:"status"` // pending, linked
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}
```

**Endpoints:**

```
GET    /classrooms/{classroom_id}/roster
POST   /classrooms/{classroom_id}/roster
DELETE /classrooms/{classroom_id}/roster/{id}
POST   /classrooms/{classroom_id}/roster/bulk
POST   /classrooms/{classroom_id}/roster/import
PATCH  /classrooms/{classroom_id}/roster/{id}/link
```

#### **5.1.5 Team**

**Resource Model:**
```go
type Team struct {
    ID           int       `json:"id"`
    AssignmentID int       `json:"assignment_id"`
    Name         string    `json:"name"`
    Slug         string    `json:"slug"`
    Members      []int     `json:"member_ids"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

**Endpoints:**

```
GET    /assignments/{assignment_id}/teams
POST   /assignments/{assignment_id}/teams
GET    /teams/{id}
PATCH  /teams/{id}
DELETE /teams/{id}
POST   /teams/{id}/members
DELETE /teams/{id}/members/{student_id}
```

### **5.2 Complete OpenAPI Specification**

See **Section 14** for the full OpenAPI 3.0 specification with all endpoints, schemas, and examples.

---

## **6. Error Handling Standards**

### **6.1 Error Response Structure**

**Standard Error Format:**

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {
      "field": "Additional context",
      "nested": {
        "key": "value"
      }
    },
    "request_id": "req_abc123xyz",
    "timestamp": "2025-10-29T12:34:56Z",
    "documentation_url": "https://docs.forgejo-classroom.org/errors/ERROR_CODE"
  }
}
```

### **6.2 Error Code Taxonomy**

**Format**: `CATEGORY_SPECIFIC_ERROR`

**Categories:**

```go
const (
    // Authentication Errors (AUTH_*)
    ErrAuthMissingToken   = "AUTH_MISSING_TOKEN"
    ErrAuthInvalidToken   = "AUTH_INVALID_TOKEN"
    ErrAuthExpiredToken   = "AUTH_EXPIRED_TOKEN"

    // Authorization Errors (AUTHZ_*)
    ErrAuthzForbidden     = "AUTHZ_FORBIDDEN"
    ErrAuthzInsufficientPermissions = "AUTHZ_INSUFFICIENT_PERMISSIONS"

    // Validation Errors (VALIDATION_*)
    ErrValidationInvalidInput = "VALIDATION_INVALID_INPUT"
    ErrValidationMissingField = "VALIDATION_MISSING_REQUIRED_FIELD"
    ErrValidationInvalidFormat = "VALIDATION_INVALID_FORMAT"
    ErrValidationInvalidDate = "VALIDATION_INVALID_DATE"
    ErrValidationTooShort = "VALIDATION_TOO_SHORT"
    ErrValidationTooLong = "VALIDATION_TOO_LONG"

    // Resource Errors (RESOURCE_*)
    ErrResourceNotFound   = "RESOURCE_NOT_FOUND"
    ErrResourceConflict   = "RESOURCE_CONFLICT"
    ErrResourceAlreadyExists = "RESOURCE_ALREADY_EXISTS"

    // Business Logic Errors (BUSINESS_*)
    ErrBusinessDeadlinePassed = "BUSINESS_DEADLINE_PASSED"
    ErrBusinessAlreadyAccepted = "BUSINESS_ALREADY_ACCEPTED"
    ErrBusinessRosterNotFound = "BUSINESS_ROSTER_NOT_FOUND"
    ErrBusinessTeamSizeExceeded = "BUSINESS_TEAM_SIZE_EXCEEDED"
    ErrBusinessTemplateNotFound = "BUSINESS_TEMPLATE_NOT_FOUND"

    // Integration Errors (INTEGRATION_*)
    ErrIntegrationForgejoAPI = "INTEGRATION_FORGEJO_API_ERROR"
    ErrIntegrationForgejoRateLimited = "INTEGRATION_FORGEJO_RATE_LIMITED"
    ErrIntegrationForgejoUnavailable = "INTEGRATION_FORGEJO_UNAVAILABLE"
    ErrIntegrationDatabase = "INTEGRATION_DATABASE_ERROR"

    // System Errors (SYSTEM_*)
    ErrSystemInternal     = "SYSTEM_INTERNAL_ERROR"
    ErrSystemUnavailable  = "SYSTEM_UNAVAILABLE"
    ErrSystemTimeout      = "SYSTEM_TIMEOUT"
)
```

### **6.3 HTTP Status Code Mapping**

| Error Category | HTTP Status | Example Codes |
|---------------|-------------|---------------|
| Authentication | 401 | AUTH_* |
| Authorization | 403 | AUTHZ_* |
| Validation | 400 | VALIDATION_* |
| Not Found | 404 | RESOURCE_NOT_FOUND |
| Conflict | 409 | RESOURCE_CONFLICT, RESOURCE_ALREADY_EXISTS |
| Business Logic | 422 | BUSINESS_* |
| Integration | 502, 503 | INTEGRATION_* |
| System | 500, 503, 504 | SYSTEM_* |

### **6.4 Validation Error Details**

**Field-level validation errors:**

```json
{
  "error": {
    "code": "VALIDATION_INVALID_INPUT",
    "message": "Request validation failed",
    "details": {
      "fields": [
        {
          "field": "deadline",
          "code": "VALIDATION_INVALID_DATE",
          "message": "Deadline must be in the future",
          "received": "2024-01-01T00:00:00Z",
          "constraint": "future date"
        },
        {
          "field": "title",
          "code": "VALIDATION_TOO_SHORT",
          "message": "Title must be at least 1 character",
          "received": "",
          "constraint": "min_length: 1"
        },
        {
          "field": "max_team_size",
          "code": "VALIDATION_OUT_OF_RANGE",
          "message": "Max team size must be between 2 and 10",
          "received": 15,
          "constraint": "range: [2, 10]"
        }
      ]
    },
    "request_id": "req_abc123",
    "timestamp": "2025-10-29T12:34:56Z"
  }
}
```

### **6.5 Error Response Implementation**

```go
// internal/api/errors.go
package api

import (
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type ErrorResponse struct {
    Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
    Code             string                 `json:"code"`
    Message          string                 `json:"message"`
    Details          map[string]interface{} `json:"details,omitempty"`
    RequestID        string                 `json:"request_id"`
    Timestamp        time.Time              `json:"timestamp"`
    DocumentationURL string                 `json:"documentation_url,omitempty"`
}

type ValidationError struct {
    Field      string      `json:"field"`
    Code       string      `json:"code"`
    Message    string      `json:"message"`
    Received   interface{} `json:"received,omitempty"`
    Constraint string      `json:"constraint,omitempty"`
}

// Error codes to HTTP status mapping
var errorStatusMap = map[string]int{
    // Authentication
    "AUTH_MISSING_TOKEN":   http.StatusUnauthorized,
    "AUTH_INVALID_TOKEN":   http.StatusUnauthorized,
    "AUTH_EXPIRED_TOKEN":   http.StatusUnauthorized,

    // Authorization
    "AUTHZ_FORBIDDEN":                  http.StatusForbidden,
    "AUTHZ_INSUFFICIENT_PERMISSIONS":   http.StatusForbidden,

    // Validation
    "VALIDATION_INVALID_INPUT":        http.StatusBadRequest,
    "VALIDATION_MISSING_REQUIRED_FIELD": http.StatusBadRequest,
    "VALIDATION_INVALID_FORMAT":       http.StatusBadRequest,
    "VALIDATION_INVALID_DATE":         http.StatusBadRequest,
    "VALIDATION_TOO_SHORT":            http.StatusBadRequest,
    "VALIDATION_TOO_LONG":             http.StatusBadRequest,
    "VALIDATION_OUT_OF_RANGE":         http.StatusBadRequest,

    // Resource
    "RESOURCE_NOT_FOUND":      http.StatusNotFound,
    "RESOURCE_CONFLICT":       http.StatusConflict,
    "RESOURCE_ALREADY_EXISTS": http.StatusConflict,

    // Business Logic
    "BUSINESS_DEADLINE_PASSED":      http.StatusUnprocessableEntity,
    "BUSINESS_ALREADY_ACCEPTED":     http.StatusUnprocessableEntity,
    "BUSINESS_ROSTER_NOT_FOUND":     http.StatusUnprocessableEntity,
    "BUSINESS_TEAM_SIZE_EXCEEDED":   http.StatusUnprocessableEntity,
    "BUSINESS_TEMPLATE_NOT_FOUND":   http.StatusUnprocessableEntity,

    // Integration
    "INTEGRATION_FORGEJO_API_ERROR":      http.StatusBadGateway,
    "INTEGRATION_FORGEJO_RATE_LIMITED":   http.StatusServiceUnavailable,
    "INTEGRATION_FORGEJO_UNAVAILABLE":    http.StatusServiceUnavailable,
    "INTEGRATION_DATABASE_ERROR":         http.StatusInternalServerError,

    // System
    "SYSTEM_INTERNAL_ERROR":  http.StatusInternalServerError,
    "SYSTEM_UNAVAILABLE":     http.StatusServiceUnavailable,
    "SYSTEM_TIMEOUT":         http.StatusGatewayTimeout,
}

// RespondError sends a standardized error response
func RespondError(c *gin.Context, code string, message string, details map[string]interface{}) {
    requestID, _ := c.Get("request_id")
    if requestID == nil {
        requestID = uuid.New().String()
    }

    status := errorStatusMap[code]
    if status == 0 {
        status = http.StatusInternalServerError
        code = "SYSTEM_INTERNAL_ERROR"
    }

    errorDetail := ErrorDetail{
        Code:             code,
        Message:          message,
        Details:          details,
        RequestID:        requestID.(string),
        Timestamp:        time.Now(),
        DocumentationURL: fmt.Sprintf("https://docs.forgejo-classroom.org/errors/%s", code),
    }

    c.JSON(status, ErrorResponse{Error: errorDetail})
}

// RespondValidationError sends a validation error with field details
func RespondValidationError(c *gin.Context, fields []ValidationError) {
    RespondError(c, "VALIDATION_INVALID_INPUT", "Request validation failed", map[string]interface{}{
        "fields": fields,
    })
}

// Example usage in handler
func (h *AssignmentHandler) CreateAssignment(c *gin.Context) {
    var req CreateAssignmentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        fields := parseValidationErrors(err)
        RespondValidationError(c, fields)
        return
    }

    assignment, err := h.service.Create(c.Request.Context(), req)
    if err != nil {
        handleServiceError(c, err)
        return
    }

    c.JSON(http.StatusCreated, assignment)
}
```

---

## **7. API Versioning Strategy**

### **7.1 Versioning Policy**

**Version Format**: `/api/classroom/v{major}`

**Current**: v1
**Path**: `/api/classroom/v1`

### **7.2 Breaking vs Non-Breaking Changes**

**Breaking Changes (require new major version):**
- Removing endpoints
- Removing request fields
- Removing response fields
- Changing field types
- Changing required/optional status of fields
- Changing authentication mechanisms
- Changing HTTP status codes for existing errors
- Changing URL structure

**Non-Breaking Changes (allowed in same version):**
- Adding new endpoints
- Adding optional request fields
- Adding response fields
- Adding new enum values (with unknown value handling)
- Adding new error codes
- Performance improvements
- Bug fixes

### **7.3 Version Support Policy**

Version Timeline:
v1 (current)  → Supported until v3 release + 12 months
v2 (next)     → Supported for 18 months after v3 release
v3 (future)   → Current version

Minimum Support: 24 months from release

### **7.4 Deprecation Process**

**6-Month Notice Period:**

1. **Announcement** (T-6 months)
   - Release notes
   - Email to registered API users
   - Documentation update

2. **Headers Added** (T-6 months)
   ```http
   HTTP/1.1 200 OK
   Deprecation: true
   Sunset: Sat, 29 Apr 2026 23:59:59 GMT
   Link: <https://docs.forgejo-classroom.org/migrations/v1-to-v2>; rel="deprecation"
   ```

3. **Warning Period** (T-3 months)
   - Deprecation warnings in logs
   - API dashboard shows deprecation status

4. **Removal** (T+0)
   - Endpoint returns 410 Gone
   - Clear migration path documented

### **7.5 Version Migration Guide Template**

# Migration Guide: v1 to v2

## Breaking Changes

### 1. Assignment Acceptance Endpoint Changed

**Old (v1):**

```
POST /classrooms/{id}/assignments/{id}/accept
```

**New (v2):**
```
POST /assignments/{id}/submissions
```

**Rationale:** Better REST semantics - accepting creates a submission resource.

**Migration:**
```diff
- POST /api/classroom/v1/classrooms/1/assignments/5/accept
+ POST /api/classroom/v2/assignments/5/submissions
```

### 2. Pagination Response Structure

**Old (v1):**
```json
{
  "classrooms": [...],
  "total": 100
}
```

**New (v2):**
```json
{
  "data": [...],
  "pagination": {
    "total_count": 100,
    "page": 1,
    "per_page": 30
  }
}
```

## Deprecated Endpoints

- `GET /classrooms/{id}/stats` → Use `GET /classrooms/{id}` (includes stats)

## Timeline

- v1 EOL: April 29, 2026
- v2 Release: October 29, 2025

---

## **8. Infrastructure Architecture**

### **8.1 System Components**

```
┌─────────────────────────────────────────────────────────────┐
│                         Load Balancer                        │
│                         (nginx/haproxy)                      │
└─────────────────┬───────────────────────────────────────────┘
                  │
    ┌─────────────┴─────────────┐
    │                           │
┌───▼────────────┐    ┌────────▼──────────┐
│  API Server 1  │    │   API Server 2    │
│   (Go binary)  │    │    (Go binary)    │
└───┬────────────┘    └────────┬──────────┘
    │                           │
    └─────────────┬─────────────┘
                  │
    ┌─────────────┼─────────────────────────┐
    │             │                         │
┌───▼─────┐  ┌───▼──────┐  ┌──────▼──────┐
│ Redis   │  │PostgreSQL│  │   Forgejo   │
│(Cache + │  │(Primary  │  │  (External) │
│ Queue)  │  │Database) │  │             │
└─────────┘  └──────────┘  └─────────────┘
```

### **8.2 Caching Strategy**

**Cache Layer Architecture:**

```go
// internal/cache/cache.go
package cache

import (
    "context"
    "time"
)

type Cache interface {
    Get(ctx context.Context, key string) ([]byte, error)
    Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
    DeletePattern(ctx context.Context, pattern string) error
}

// Redis implementation
type RedisCache struct {
    client *redis.Client
}

// In-memory fallback
type MemoryCache struct {
    store map[string]cacheEntry
    mu    sync.RWMutex
}
```

**Cache TTL Strategy:**

```yaml
cache_ttl:
  # Static data - long TTL
  classrooms: 300s      # 5 minutes
  assignments: 300s     # 5 minutes
  roster: 300s          # 5 minutes

  # Dynamic data - short TTL
  submissions: 60s      # 1 minute
  statistics: 30s       # 30 seconds

  # User-specific - session cache
  permissions: 900s     # 15 minutes (session duration)
  user_profile: 900s    # 15 minutes
```

**Cache Key Patterns:**

```go
// internal/cache/keys.go
package cache

const (
    KeyClassroom        = "classroom:%d"                    // classroom:123
    KeyClassroomList    = "classroom:list:%d:%d"            // classroom:list:1:30 (page:per_page)
    KeyAssignment       = "assignment:%d"                   // assignment:456
    KeyAssignmentList   = "classroom:%d:assignments"        // classroom:123:assignments
    KeySubmission       = "submission:%d"                   // submission:789
    KeySubmissionList   = "assignment:%d:submissions"       // assignment:456:submissions
    KeyRoster           = "classroom:%d:roster"             // classroom:123:roster
    KeyPermission       = "permission:%d:%d"                // permission:user_id:classroom_id
    KeyStats            = "stats:assignment:%d"             // stats:assignment:456
)

func ClassroomKey(id int) string {
    return fmt.Sprintf(KeyClassroom, id)
}

func PermissionKey(userID, classroomID int) string {
    return fmt.Sprintf(KeyPermission, userID, classroomID)
}
```

**Cache Invalidation:**

```go
// Service layer handles invalidation
func (s *ClassroomService) Update(ctx context.Context, id int, req UpdateRequest) (*model.Classroom, error) {
    // Update database
    classroom, err := s.repo.Update(ctx, id, req)
    if err != nil {
        return nil, err
    }

    // Invalidate caches
    s.cache.Delete(ctx, cache.ClassroomKey(id))
    s.cache.DeletePattern(ctx, "classroom:list:*")

    // Publish event for other services
    s.eventBus.Publish("classroom.updated", ClassroomUpdatedEvent{
        ClassroomID: id,
    })

    return classroom, nil
}
```

### **8.3 Message Queue for Async Operations**

**Queue Architecture:**

```go
// internal/queue/queue.go
package queue

import "context"

type Queue interface {
    Enqueue(ctx context.Context, job Job) error
    Dequeue(ctx context.Context) (Job, error)
    Ack(ctx context.Context, jobID string) error
    Nack(ctx context.Context, jobID string) error
}

type Job struct {
    ID       string
    Type     string
    Payload  []byte
    Priority int
    Retry    int
    MaxRetry int
}

// Redis-backed queue implementation
type RedisQueue struct {
    client *redis.Client
    name   string
}
```

**Job Definitions:**

```go
// internal/queue/jobs/deadline.go
package jobs

type DeadlineEnforcementJob struct {
    AssignmentID int       `json:"assignment_id"`
    Deadline     time.Time `json:"deadline"`
}

func (j *DeadlineEnforcementJob) Execute(ctx context.Context) error {
    // Get all submissions for assignment
    submissions, err := submissionRepo.ListByAssignment(ctx, j.AssignmentID)
    if err != nil {
        return err
    }

    // For each submission, create deadline tag
    for _, sub := range submissions {
        if err := gitService.CreateDeadlineTag(ctx, sub.RepositoryID, j.Deadline); err != nil {
            log.Error("Failed to create deadline tag",
                "submission_id", sub.ID,
                "error", err,
            )
            // Continue with other submissions
        }
    }

    return nil
}
```

**Worker Pool:**

```go
// internal/queue/worker.go
package queue

type Worker struct {
    queue   Queue
    handler JobHandler
}

func (w *Worker) Start(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            job, err := w.queue.Dequeue(ctx)
            if err != nil {
                time.Sleep(time.Second)
                continue
            }

            if err := w.handler.Handle(ctx, job); err != nil {
                if job.Retry < job.MaxRetry {
                    job.Retry++
                    w.queue.Enqueue(ctx, job)
                } else {
                    log.Error("Job failed after max retries", "job_id", job.ID)
                }
                w.queue.Nack(ctx, job.ID)
            } else {
                w.queue.Ack(ctx, job.ID)
            }
        }
    }
}
```

**Queue Configuration:**

```yaml
queues:
  deadlines:
    workers: 2
    priority: high
    max_retry: 3
    retry_delay: 5m

  bulk_operations:
    workers: 5
    priority: medium
    max_retry: 3
    retry_delay: 2m

  notifications:
    workers: 3
    priority: low
    max_retry: 5
    retry_delay: 1m
```

### **8.4 Database Architecture**

**Connection Pooling:**

```go
// internal/repository/db.go
package repository

import (
    "database/sql"
    "time"
)

type DBConfig struct {
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
    ConnMaxIdleTime time.Duration
}

func NewDB(connStr string, config DBConfig) (*sql.DB, error) {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }

    // Connection pool settings
    db.SetMaxOpenConns(config.MaxOpenConns)       // Default: 25
    db.SetMaxIdleConns(config.MaxIdleConns)       // Default: 5
    db.SetConnMaxLifetime(config.ConnMaxLifetime) // Default: 1 hour
    db.SetConnMaxIdleTime(config.ConnMaxIdleTime) // Default: 10 minutes

    // Verify connection
    if err := db.Ping(); err != nil {
        return nil, err
    }

    return db, nil
}
```

**Transaction Wrapper Pattern:**

```go
// internal/repository/transaction.go
package repository

import (
    "context"
    "database/sql"
)

type TxFunc func(*sql.Tx) error

func (r *Repository) WithTransaction(ctx context.Context, fn TxFunc) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }

    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        }
    }()

    if err := fn(tx); err != nil {
        if rbErr := tx.Rollback(); rbErr != nil {
            return fmt.Errorf("tx error: %v, rollback error: %v", err, rbErr)
        }
        return err
    }

    return tx.Commit()
}

// Usage example
func (s *AssignmentService) Create(ctx context.Context, req CreateAssignmentRequest) (*model.Assignment, error) {
    var assignment *model.Assignment

    err := s.repo.WithTransaction(ctx, func(tx *sql.Tx) error {
        // Create assignment
        assignment = &model.Assignment{
            ClassroomID: req.ClassroomID,
            Title:       req.Title,
            // ... other fields
        }

        if err := s.assignmentRepo.CreateTx(ctx, tx, assignment); err != nil {
            return err
        }

        // Verify template repository exists (external API call)
        if _, err := s.forgejoClient.GetRepository(ctx, req.TemplateRepo); err != nil {
            return fmt.Errorf("template repository not found: %w", err)
        }

        return nil
    })

    if err != nil {
        return nil, err
    }

    // Invalidate cache after successful commit
    s.cache.DeletePattern(ctx, fmt.Sprintf("classroom:%d:assignments*", req.ClassroomID))

    return assignment, nil
}
```

### **8.5 Rate Limiting Architecture**

**Tiered Rate Limiting:**

```go
// internal/api/middleware/ratelimit.go
package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/ulule/limiter/v3"
    "github.com/ulule/limiter/v3/drivers/store/redis"
)

type RateLimitTier struct {
    Read  int // requests per minute
    Write int
    Bulk  int
}

var RateLimits = map[string]RateLimitTier{
    "owner":     {Read: 1000, Write: 300, Bulk: 50},
    "admin":     {Read: 500, Write: 150, Bulk: 30},
    "student":   {Read: 200, Write: 50, Bulk: 5},
    "anonymous": {Read: 10, Write: 0, Bulk: 0},
}

func RateLimitByRole(category string) gin.HandlerFunc {
    return func(c *gin.Context) {
        role := c.GetString("user_role")
        if role == "" {
            role = "anonymous"
        }

        tier := RateLimits[role]
        var limit int

        switch category {
        case "read":
            limit = tier.Read
        case "write":
            limit = tier.Write
        case "bulk":
            limit = tier.Bulk
        default:
            limit = tier.Read
        }

        // Use user ID as key, or IP for anonymous
        key := fmt.Sprintf("%d", c.GetInt("user_id"))
        if key == "0" {
            key = c.ClientIP()
        }

        // Check rate limit
        context := limiter.Get(c.Request.Context(), key)

        c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
        c.Header("X-RateLimit-Remaining", strconv.Itoa(context.Remaining))
        c.Header("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))

        if context.Reached {
            c.AbortWithStatusJSON(429, gin.H{
                "error": gin.H{
                    "code":    "RATE_LIMIT_EXCEEDED",
                    "message": "Rate limit exceeded",
                    "retry_after": context.Reset,
                },
            })
            return
        }

        c.Next()
    }
}
```

**Rate Limit Application:**

```go
// Apply to routes
router.GET("/classrooms",
    RateLimitByRole("read"),
    handler.ListClassrooms,
)

router.POST("/classrooms",
    RateLimitByRole("write"),
    RequirePermission(PermCreateClassroom),
    handler.CreateClassroom,
)

router.POST("/assignments/:id/submissions",
    RateLimitByRole("bulk"),  // Creating repo is expensive
    RequirePermission(PermAcceptAssignment),
    handler.AcceptAssignment,
)
```

### **8.6 Context Handling and Cancellation**

**Request Context Propagation:**

```go
// Middleware to add timeout to request context
func RequestTimeout(timeout time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
        defer cancel()

        c.Request = c.Request.WithContext(ctx)
        c.Next()
    }
}

// Service layer respects context cancellation
func (s *SubmissionService) DownloadAll(ctx context.Context, assignmentID int) error {
    submissions, err := s.repo.ListSubmissions(ctx, assignmentID)
    if err != nil {
        return err
    }

    for i, sub := range submissions {
        // Check if context was cancelled
        select {
        case <-ctx.Done():
            return fmt.Errorf("download cancelled after %d/%d: %w", i, len(submissions), ctx.Err())
        default:
        }

        if err := s.downloadRepository(ctx, sub); err != nil {
            return fmt.Errorf("failed to download submission %d: %w", sub.ID, err)
        }
    }

    return nil
}
```

---

## **9. Implementation Phases**

### **Phase 1: Foundation (Weeks 1-2)**

**Goal:** Establish core infrastructure and data models

**Deliverables:**
- ✅ Project structure setup
- ✅ Database schema design
- ✅ Migration system
- ✅ Configuration management
- ✅ Basic Forgejo API client
- ✅ Authentication middleware
- ✅ Error handling framework
- ✅ Logging infrastructure
- ✅ Cache layer interface
- ✅ Transaction wrapper pattern

**Success Criteria:**
- Database migrations run successfully
- Forgejo API client can authenticate and make basic calls
- Configuration loads from file and environment variables
- Structured logging outputs to stdout/file
- Error responses follow standard format
- Transaction wrapper properly rolls back on error

**Key Tasks:**
```bash
# Initialize project
mkdir -p forgejo-classroom/{cmd,internal,pkg,migrations,test,docs}
go mod init code.forgejo.org/forgejo/classroom

# Install dependencies
go get github.com/gin-gonic/gin
go get github.com/lib/pq
go get github.com/spf13/cobra
go get code.gitea.io/sdk/gitea

# Create initial migrations
migrate create -ext sql -dir migrations -seq create_classrooms
migrate create -ext sql -dir migrations -seq create_assignments
migrate create -ext sql -dir migrations -seq create_roster
migrate create -ext sql -dir migrations -seq create_submissions

# Setup development environment
docker-compose up -d postgres redis

# Run migrations
make migrate-up

# Verify Forgejo connectivity
go run cmd/fgc-server/main.go --dry-run
```

### **Phase 2: Classroom Management (Weeks 3-4)**

**Goal:** Implement classroom CRUD operations

**Deliverables:**
- ✅ Classroom service layer with caching
- ✅ Classroom API endpoints (OpenAPI compliant)
- ✅ Classroom CLI commands
- ✅ Organization creation/management integration
- ✅ Owner permission checks
- ✅ Unit tests for classroom logic (>80% coverage)
- ✅ Integration tests for API
- ✅ Contract tests against OpenAPI spec

**Success Criteria:**
- Teachers can create classrooms via CLI
- Classrooms correctly map to Forgejo organizations
- Permissions prevent unauthorized access
- API returns proper error codes (standardized taxonomy)
- Tests achieve >80% coverage
- Cache invalidation works correctly
- API responses match OpenAPI specification

**Example Workflow:**
```bash
# Teacher creates classroom
fgc classroom create \
  --name "CS101 Fall 2025" \
  --org cs101-fall2025 \
  --description "Introduction to Computer Science"

# CLI output:
# Created classroom "CS101 Fall 2025" (ID: 1)
# Organization: https://forgejo.example.com/cs101-fall2025
# API: GET /api/classroom/v1/classrooms/1

# List classrooms
fgc classroom list
# ID  Name                Status   Students  Assignments
# 1   CS101 Fall 2025     active   0         0

# View classroom details (uses cache)
fgc classroom view 1
```

### **Phase 3: Roster Management (Weeks 5-6)**

**Goal:** Implement student roster management

**Deliverables:**
- ✅ Roster service layer
- ✅ Roster API endpoints with bulk operations
- ✅ Roster CLI commands
- ✅ CSV import functionality (async for large files)
- ✅ Student-to-Forgejo account linking
- ✅ Bulk operations support (sync and async)
- ✅ Job queue for large imports
- ✅ Unit and integration tests

**Success Criteria:**
- Teachers can import rosters from CSV
- Small imports (<100 students) are synchronous
- Large imports (100+ students) are asynchronous with status tracking
- Students can be linked to Forgejo accounts
- Roster operations are idempotent
- Handles duplicate entries gracefully
- Supports batch operations efficiently
- Job status can be queried

**CSV Format:**
```csv
identifier,email,full_name
student001,john.doe@university.edu,John Doe
student002,jane.smith@university.edu,Jane Smith
```

**Example Workflow:**
```bash
# Small roster - synchronous
fgc roster add 1 students_small.csv
# Imported 45 students to classroom "CS101 Fall 2025"

# Large roster - asynchronous
fgc roster import 1 students_large.csv
# Job ID: import_abc123
# Status: processing
# Check status: fgc job status import_abc123

# Check job status
fgc job status import_abc123
# Job: import_abc123
# Status: completed
# Progress: 500/500 (485 succeeded, 15 failed)
# Errors: fgc job errors import_abc123

# Link student to Forgejo account
fgc roster link 1 student001 --username jdoe
# Linked student001 to Forgejo user @jdoe

# View roster (cached for 5 minutes)
fgc roster list 1
# ID   Identifier   Name           Forgejo User  Status
# 10   student001   John Doe       @jdoe         linked
# 11   student002   Jane Smith     -             pending
```

### **Phase 4: Assignment Distribution (Weeks 7-9)**

**Goal:** Implement assignment creation and distribution

**Deliverables:**
- ✅ Assignment service layer with caching
- ✅ Assignment API endpoints (RESTful)
- ✅ Assignment CLI commands
- ✅ Template repository validation
- ✅ Repository creation from templates
- ✅ Invitation code generation
- ✅ Assignment acceptance flow (creates submissions)
- ✅ Rate limiting for repository creation
- ✅ Comprehensive testing

**Success Criteria:**
- Teachers can create assignments from templates
- Students can accept assignments and get repositories
- Repository naming follows conventions
- Permissions are set correctly
- Template validation prevents errors
- Handles concurrent acceptances safely
- Rate limiting prevents abuse
- Acceptance creates submission resource

**Example Workflow:**
```bash
# Teacher creates assignment
fgc assignment create 1 \
  --title "Homework 1: Variables" \
  --slug hw01 \
  --template cs101-templates/hw01-starter \
  --deadline "2025-11-15T23:59:59Z" \
  --type individual

# Output:
# Created assignment "Homework 1: Variables" (ID: 10)
# Invitation URL: https://forgejo.example.com/classroom/accept/hw01/xyz123
# Share this link with students

# Student accepts assignment (creates submission)
fgc student accept hw01/xyz123
# Accepting assignment "Homework 1: Variables"...
# Created submission (ID: 42)
# Repository: cs101-fall2025/hw01-jdoe
# Clone URL: https://forgejo.example.com/cs101-fall2025/hw01-jdoe.git

# Student clones and works
git clone https://forgejo.example.com/cs101-fall2025/hw01-jdoe.git
cd hw01-jdoe
# ... work on assignment ...
git add .
git commit -m "Completed exercise 1"
git push
```

### **Phase 5: Submission Management (Weeks 10-11)**

**Goal:** Implement submission tracking and deadline enforcement

**Deliverables:**
- ✅ Submission service layer
- ✅ Submission API endpoints (top-level resource)
- ✅ Submission CLI commands
- ✅ Deadline enforcement mechanism (Git tagging via queue)
- ✅ Submission statistics with caching
- ✅ Bulk download functionality
- ✅ Async job for deadline tagging
- ✅ Unit and integration tests

**Success Criteria:**
- Teachers can view all submissions
- Flexible querying (by assignment, student, classroom, status)
- Deadline tagging works automatically via job queue
- Teachers can download all submissions as archive
- Submission status accurately reflects student progress
- Statistics provide meaningful insights (cached)
- Deadline jobs can be scheduled and monitored

**Deadline Enforcement Mechanism:**
```bash
# Scheduled job (runs via cron/scheduler)
fgc submission enforce-deadline 10

# Enqueues job for each student repository:
# 1. Create tag: deadline-2025-11-15T23:59:59Z
# 2. Point to HEAD of default branch at deadline time
# 3. Record tag name in submissions table
# 4. Update submission status

# Teacher can then clone at deadline:
git clone --branch deadline-2025-11-15T23:59:59Z \
  https://forgejo.example.com/cs101-fall2025/hw01-jdoe.git
```

**Example Workflow:**
```bash
# View submission statistics (cached for 30s)
fgc assignment stats 10
# Assignment: Homework 1: Variables
# Total students:      45
# Accepted:            42 (93%)
# Submitted (on-time): 38 (84%)
# Submitted (late):    3 (7%)
# Not submitted:       4 (9%)
# Deadline: 2025-11-15 23:59:59

# List all submissions (flexible filtering)
fgc submission list --assignment=10
# ID  Student      Repository      Commits  Last Push            Status
# 42  John Doe     hw01-jdoe       15       2025-11-14 22:30     on_time
# 43  Jane Smith   hw01-jsmith     8        2025-11-16 10:00     late
# 44  Bob Wilson   hw01-bwilson    12       2025-11-15 23:45     on_time

# Filter by status
fgc submission list --assignment=10 --status=late
# ID  Student      Repository      Commits  Last Push            Status
# 43  Jane Smith   hw01-jsmith     8        2025-11-16 10:00     late

# Download all submissions for grading
fgc submission download 10 --output ./grading/hw01/
# Downloading 42 repositories...
# [========================================] 100%
# Downloaded to: ./grading/hw01/
# ├── hw01-jdoe/
# ├── hw01-jsmith/
# └── ...

# View individual submission details
fgc submission view 42
# Submission ID: 42
# Assignment: Homework 1: Variables
# Student: John Doe (student001)
# Repository: cs101-fall2025/hw01-jdoe
# Status: on_time
# Accepted: 2025-11-01 14:30:00
# Last Commit: 2025-11-14 22:30:00 (abc123def)
# Commits: 15
# Deadline Tag: deadline-2025-11-15T23:59:59Z
```

### **Phase 6: Team Assignments (Weeks 12-13)**

**Goal:** Implement team-based assignments

**Deliverables:**
- ✅ Team service layer
- ✅ Team API endpoints
- ✅ Team CLI commands
- ✅ Forgejo team integration
- ✅ Team repository creation
- ✅ Team member management
- ✅ Team size limits enforcement
- ✅ Testing

**Success Criteria:**
- Teachers can create team assignments
- Students can form teams
- Team size limits are enforced
- Team permissions work correctly
- One repository per team, shared access
- Team membership changes reflected in Forgejo

**Example Workflow:**
```bash
# Teacher creates team assignment
fgc assignment create 1 \
  --title "Final Project" \
  --slug final-project \
  --template cs101-templates/project-starter \
  --type team \
  --max-team-size 4 \
  --deadline "2025-12-15T23:59:59Z"

# Student 1 creates team
fgc team create final-project \
  --name "Team Alpha" \
  --members student001,student002,student003
# Created team "Team Alpha" (ID: 5)
# Repository: cs101-fall2025/final-project-team-alpha
# Members can now clone and collaborate

# Student 2 joins existing team
fgc team join final-project --team-id 5
# Joined team "Team Alpha"
# Repository access granted

# Teacher views teams
fgc team list final-project
# ID  Name         Members  Repository
# 5   Team Alpha   4        final-project-team-alpha
# 6   Team Beta    3        final-project-team-beta

# All team members get access to shared repository:
# cs101-fall2025/final-project-team-alpha
```

### **Phase 7: Polish and Documentation (Week 14)**

**Goal:** Production readiness

**Deliverables:**
- ✅ Comprehensive API documentation (OpenAPI 3.0 spec)
- ✅ CLI usage guide with examples
- ✅ Installation and deployment guide
- ✅ Developer contribution guide
- ✅ Performance optimization
- ✅ Security audit
- ✅ Error message improvements
- ✅ Cache tuning
- ✅ Rate limit configuration
- ✅ Release preparation

**Success Criteria:**
- Documentation covers all features
- Installation takes <10 minutes
- Security vulnerabilities addressed
- Performance meets targets (see below)
- API contract tests pass
- Ready for alpha release

**Performance Targets:**

| Operation | Target | Measurement | Status |
|-----------|--------|-------------|--------|
| Create classroom | <200ms | p95 latency | ✅ |
| Create assignment | <300ms | p95 latency | ✅ |
| Accept assignment | <2s | p95 latency | ✅ |
| List submissions (50 students) | <500ms | p95 latency | ✅ |
| Download all submissions | <10s | For 50 students | ✅ |
| API throughput | >100 req/s | Sustained load | ✅ |
| Cache hit rate | >80% | For read operations | Target |

---

## **10. Testing Strategy**

### **10.1 Testing Pyramid**

```
                    /\
                   /  \
                  / E2E\
                 /  10% \
                /________\
               /          \
              / Integration\
             /     30%      \
            /______________  \
           /                  \
          /   Unit Tests 60%   \
         /______________________\
```

### **10.2 Unit Testing**

**Target Coverage:** 80%+ for business logic

**Example: Testing Assignment Service**

```go
// internal/service/assignment_test.go
package service

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/suite"
)

type AssignmentServiceTestSuite struct {
    suite.Suite
    service       *AssignmentService
    mockRepo      *MockAssignmentRepo
    mockForgejo   *MockForgejoClient
    mockCache     *MockCache
}

func (s *AssignmentServiceTestSuite) SetupTest() {
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
    ctx := context.Background()
    req := CreateAssignmentRequest{
        ClassroomID:  1,
        Title:        "Homework 1",
        Type:         "individual",
        TemplateRepo: "org/repo",
        Deadline:     time.Now().Add(7 * 24 * time.Hour),
    }

    // Setup mocks
    s.mockRepo.On("WithTransaction", mock.Anything, mock.Anything).
        Return(nil).
        Run(func(args mock.Arguments) {
            fn := args.Get(1).(func(*sql.Tx) error)
            fn(nil) // Execute transaction function
        })

    s.mockForgejo.On("GetRepository", mock.Anything, "org/repo").
        Return(&Repository{ID: 100}, nil)

    s.mockCache.On("DeletePattern", mock.Anything, "classroom:1:assignments*").
        Return(nil)

    // Execute
    assignment, err := s.service.Create(ctx, req)

    // Assert
    s.NoError(err)
    s.NotNil(assignment)
    s.Equal(req.Title, assignment.Title)
    s.NotEmpty(assignment.InvitationCode)

    // Verify mocks
    s.mockRepo.AssertExpectations(s.T())
    s.mockForgejo.AssertExpectations(s.T())
    s.mockCache.AssertExpectations(s.T())
}

func (s *AssignmentServiceTestSuite) TestCreate_DeadlineInPast() {
    ctx := context.Background()
    req := CreateAssignmentRequest{
        ClassroomID: 1,
        Title:       "Homework 1",
        Deadline:    time.Now().Add(-1 * time.Hour),
    }

    // Execute
    assignment, err := s.service.Create(ctx, req)

    // Assert
    s.Error(err)
    s.Nil(assignment)
    s.Contains(err.Error(), "deadline must be in the future")
}

func (s *AssignmentServiceTestSuite) TestCreate_TemplateNotFound() {
    ctx := context.Background()
    req := CreateAssignmentRequest{
        ClassroomID:  1,
        Title:        "Homework 1",
        TemplateRepo: "org/nonexistent",
        Deadline:     time.Now().Add(7 * 24 * time.Hour),
    }

    // Setup mocks
    s.mockRepo.On("WithTransaction", mock.Anything, mock.Anything).
        Return(ErrTemplateNotFound)

    s.mockForgejo.On("GetRepository", mock.Anything, "org/nonexistent").
        Return(nil, ErrNotFound)

    // Execute
    assignment, err := s.service.Create(ctx, req)

    // Assert
    s.Error(err)
    s.Nil(assignment)
    s.Equal(ErrTemplateNotFound, err)
}

func TestAssignmentServiceSuite(t *testing.T) {
    suite.Run(t, new(AssignmentServiceTestSuite))
}
```

### **10.3 Integration Testing**

**Target:** All API endpoints tested against real database

**Example: API Integration Test**

```go
// test/api/assignment_test.go
package api_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/stretchr/testify/suite"
)

type AssignmentAPITestSuite struct {
    suite.Suite
    server    *httptest.Server
    db        *sql.DB
    cache     *redis.Client
    client    *client.Client
    classroom *model.Classroom
}

func (s *AssignmentAPITestSuite) SetupSuite() {
    // Setup test database
    s.db = setupTestDB()

    // Setup test cache
    s.cache = setupTestRedis()

    // Start test server
    s.server = startTestServer(s.db, s.cache)

    // Create client
    s.client = client.New(s.server.URL, testToken)
}

func (s *AssignmentAPITestSuite) SetupTest() {
    // Clean database before each test
    cleanDatabase(s.db)
    s.cache.FlushAll(context.Background())

    // Create test classroom
    s.classroom = createTestClassroom(s.db)
}

func (s *AssignmentAPITestSuite) TearDownSuite() {
    s.server.Close()
    s.db.Close()
    s.cache.Close()
}

func (s *AssignmentAPITestSuite) TestCreateAssignment_Success() {
    req := &client.CreateAssignmentRequest{
        Title:        "Test Assignment",
        Slug:         "test-assignment",
        Type:         "individual",
        TemplateRepo: "org/template",
        Deadline:     time.Now().Add(7 * 24 * time.Hour),
    }

    assignment, err := s.client.Assignments.Create(s.classroom.ID, req)

    s.NoError(err)
    s.NotNil(assignment)
    s.Equal(req.Title, assignment.Title)
    s.NotEmpty(assignment.InvitationCode)

    // Verify in database
    var count int
    err = s.db.QueryRow("SELECT COUNT(*) FROM assignments WHERE id = $1",
        assignment.ID).Scan(&count)
    s.NoError(err)
    s.Equal(1, count)
}

func (s *AssignmentAPITestSuite) TestCreateAssignment_InvalidDeadline() {
    req := &client.CreateAssignmentRequest{
        Title:    "Test Assignment",
        Deadline: time.Now().Add(-1 * time.Hour),
    }

    _, err := s.client.Assignments.Create(s.classroom.ID, req)

    s.Error(err)
    apiErr, ok := err.(*client.APIError)
    s.True(ok)
    s.Equal("VALIDATION_INVALID_DATE", apiErr.Code)
    s.Contains(apiErr.Message, "deadline")
}

func (s *AssignmentAPITestSuite) TestListAssignments_WithPagination() {
    // Create 50 assignments
    for i := 0; i < 50; i++ {
        createTestAssignment(s.db, s.classroom.ID, fmt.Sprintf("Assignment %d", i))
    }

    // Request first page
    resp, err := s.client.Assignments.List(s.classroom.ID, client.ListOptions{
        Page:    1,
        PerPage: 20,
    })

    s.NoError(err)
    s.Len(resp.Data, 20)
    s.Equal(50, resp.Pagination.TotalCount)
    s.Equal(3, resp.Pagination.TotalPages)

    // Request second page
    resp, err = s.client.Assignments.List(s.classroom.ID, client.ListOptions{
        Page:    2,
        PerPage: 20,
    })

    s.NoError(err)
    s.Len(resp.Data, 20)
}

func (s *AssignmentAPITestSuite) TestListAssignments_CacheHit() {
    // Create assignment
    assignment := createTestAssignment(s.db, s.classroom.ID, "Test")

    // First request - cache miss
    start := time.Now()
    resp1, err := s.client.Assignments.List(s.classroom.ID, client.ListOptions{})
    duration1 := time.Since(start)

    s.NoError(err)
    s.Len(resp1.Data, 1)

    // Second request - cache hit (should be faster)
    start = time.Now()
    resp2, err := s.client.Assignments.List(s.classroom.ID, client.ListOptions{})
    duration2 := time.Since(start)

    s.NoError(err)
    s.Len(resp2.Data, 1)
    s.Less(duration2, duration1) // Cache hit should be faster
}

func TestAssignmentAPISuite(t *testing.T) {
    suite.Run(t, new(AssignmentAPITestSuite))
}
```

### **10.4 Contract Testing**

**Verify API responses match OpenAPI specification:**

```go
// test/contract/contract_test.go
package contract_test

import (
    "testing"

    "github.com/getkin/kin-openapi/openapi3"
    "github.com/getkin/kin-openapi/openapi3filter"
)

func TestAPIContractCompliance(t *testing.T) {
    // Load OpenAPI spec
    loader := openapi3.NewLoader()
    doc, err := loader.LoadFromFile("../../docs/api/openapi.yaml")
    if err != nil {
        t.Fatal(err)
    }

    // Validate spec
    if err := doc.Validate(loader.Context); err != nil {
        t.Fatal(err)
    }

    router, err := openapi3filter.NewRouter(doc)
    if err != nil {
        t.Fatal(err)
    }

    // Test actual API responses
    tests := []struct {
        name       string
        method     string
        path       string
        statusCode int
    }{
        {"List Classrooms", "GET", "/classrooms", 200},
        {"Create Classroom", "POST", "/classrooms", 201},
        {"Get Classroom", "GET", "/classrooms/1", 200},
        {"Not Found", "GET", "/classrooms/999", 404},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Make actual API request
            resp := makeAPIRequest(t, tt.method, tt.path)

            // Validate response against OpenAPI spec
            route, pathParams, err := router.FindRoute(tt.method, tt.path)
            if err != nil {
                t.Fatal(err)
            }

            requestValidationInput := &openapi3filter.RequestValidationInput{
                Request:    resp.Request,
                PathParams: pathParams,
                Route:      route,
            }

            if err := openapi3filter.ValidateRequest(loader.Context, requestValidationInput); err != nil {
                t.Errorf("Request validation failed: %v", err)
            }

            responseValidationInput := &openapi3filter.ResponseValidationInput{
                RequestValidationInput: requestValidationInput,
                Status:                 resp.StatusCode,
                Header:                 resp.Header,
                Body:                   resp.Body,
            }

            if err := openapi3filter.ValidateResponse(loader.Context, responseValidationInput); err != nil {
                t.Errorf("Response validation failed: %v", err)
            }
        })
    }
}
```

### **10.5 End-to-End Testing**

**Target:** Critical user flows tested via CLI

```bash
#!/bin/bash
# test/e2e/complete_workflow_test.sh

set -e

echo "=== E2E Test: Complete Teacher/Student Workflow ==="

# Setup
export FGC_HOST="http://localhost:3000"
export FGC_TOKEN="test_token"

# 1. Teacher creates classroom
echo "1. Creating classroom..."
CLASSROOM_ID=$(fgc classroom create \
  --name "E2E Test Classroom" \
  --org e2e-test-org \
  --output json | jq -r '.id')
echo "✓ Created classroom ID: $CLASSROOM_ID"

# 2. Teacher adds students
echo "2. Adding students to roster..."
cat > /tmp/students.csv <<EOF
identifier,email,full_name
test001,student1@test.com,Student One
test002,student2@test.com,Student Two
EOF
fgc roster add $CLASSROOM_ID /tmp/students.csv
echo "✓ Added 2 students"

# 3. Teacher creates assignment
echo "3. Creating assignment..."
ASSIGNMENT_ID=$(fgc assignment create $CLASSROOM_ID \
  --title "E2E Test Assignment" \
  --slug e2e-test-hw \
  --template test-org/test-template \
  --type individual \
  --deadline "2026-12-31T23:59:59Z" \
  --output json | jq -r '.id')
echo "✓ Created assignment ID: $ASSIGNMENT_ID"

# 4. Get invitation code
INVITATION_CODE=$(fgc assignment view $ASSIGNMENT_ID \
  --output json | jq -r '.invitation_code')
echo "✓ Invitation code: $INVITATION_CODE"

# 5. Student accepts assignment
echo "4. Student accepting assignment..."
SUBMISSION_ID=$(fgc student accept $INVITATION_CODE \
  --output json | jq -r '.id')
echo "✓ Created submission ID: $SUBMISSION_ID"

# 6. Verify submission was created
echo "5. Verifying submission..."
REPO_NAME=$(fgc submission view $SUBMISSION_ID \
  --output json | jq -r '.repository_name')
if [ -z "$REPO_NAME" ]; then
    echo "✗ Submission verification failed"
    exit 1
fi
echo "✓ Repository created: $REPO_NAME"

# 7. Teacher views submissions
echo "6. Listing submissions..."
SUBMISSION_COUNT=$(fgc submission list --assignment=$ASSIGNMENT_ID \
  --output json | jq '.data | length')
if [ "$SUBMISSION_COUNT" != "1" ]; then
    echo "✗ Expected 1 submission, got $SUBMISSION_COUNT"
    exit 1
fi
echo "✓ Found $SUBMISSION_COUNT submission"

# 8. Test caching - second request should be faster
echo "7. Testing cache..."
START=$(date +%s%N)
fgc assignment view $ASSIGNMENT_ID > /dev/null
END=$(date +%s%N)
DURATION1=$((($END - $START)/1000000))

START=$(date +%s%N)
fgc assignment view $ASSIGNMENT_ID > /dev/null
END=$(date +%s%N)
DURATION2=$((($END - $START)/1000000))

if [ $DURATION2 -lt $DURATION1 ]; then
    echo "✓ Cache working (first: ${DURATION1}ms, second: ${DURATION2}ms)"
else
    echo "⚠ Cache may not be working optimally"
fi

# Cleanup
echo "8. Cleaning up..."
fgc classroom delete $CLASSROOM_ID --force
rm /tmp/students.csv
echo "✓ Cleanup complete"

echo "=== E2E Test: PASSED ==="
```

### **10.6 Performance Testing**

**Load Testing with vegeta:**

```go
// test/load/assignment_acceptance_test.go
package load

import (
    "testing"
    "time"

    vegeta "github.com/tsenart/vegeta/v12/lib"
)

func TestAssignmentAcceptanceLoad(t *testing.T) {
    rate := vegeta.Rate{Freq: 10, Per: time.Second}
    duration := 1 * time.Minute

    targeter := vegeta.NewStaticTargeter(vegeta.Target{
        Method: "POST",
        URL:    "http://localhost:8080/api/classroom/v1/assignments/1/submissions",
        Header: http.Header{
            "Authorization": []string{"token test_token"},
            "Content-Type":  []string{"application/json"},
        },
        Body: []byte(`{"student_identifier": "test_student"}`),
    })

    attacker := vegeta.NewAttacker()

    var metrics vegeta.Metrics
    for res := range attacker.Attack(targeter, rate, duration, "Assignment Acceptance") {
        metrics.Add(res)
    }
    metrics.Close()

    // Assert performance targets
    if metrics.Latencies.P95 > 2*time.Second {
        t.Errorf("P95 latency %v exceeds target of 2s", metrics.Latencies.P95)
    }

    if metrics.Success < 0.99 {
        t.Errorf("Success rate %.2f%% below target of 99%%", metrics.Success*100)
    }

    // Log results
    t.Logf("Requests: %d", metrics.Requests)
    t.Logf("Success: %.2f%%", metrics.Success*100)
    t.Logf("P50 Latency: %v", metrics.Latencies.P50)
    t.Logf("P95 Latency: %v", metrics.Latencies.P95)
    t.Logf("P99 Latency: %v", metrics.Latencies.P99)
    t.Logf("Throughput: %.2f req/s", metrics.Rate)
}
```

---

## **11. Security Considerations**

### **11.1 Authentication and Authorization**

**Authentication Mechanism:**
- Reuse Forgejo's token-based authentication
- Support personal access tokens with appropriate scopes
- Token validation on every API request
- Token expiration and revocation support

**Authorization Model:**

```
Classroom Roles:
├── Owner (Creator)
│   ├── Full control over classroom
│   ├── Manage assignments
│   ├── Manage roster
│   └── Delete classroom
│
├── Admin (Teaching Assistant)
│   ├── Manage assignments
│   ├── View all submissions
│   ├── Grade assignments
│   └── Cannot delete classroom
│
└── Student
    ├── View own assignments
    ├── Accept assignments
    ├── Access own repositories
    └── Submit work
```

**Permission Check Implementation:**

```go
// internal/auth/permission.go
package auth

type Permission string

const (
    PermManageClassroom   Permission = "classroom:manage"
    PermViewClassroom     Permission = "classroom:view"
    PermManageAssignments Permission = "assignments:manage"
    PermGradeAssignments  Permission = "assignments:grade"
    PermViewSubmissions   Permission = "submissions:view"
    PermAcceptAssignment  Permission = "assignments:accept"
)

type Checker struct {
    forgejoClient *forgejo.Client
    repo          *repository.PermissionRepository
    cache         cache.Cache
}

func (c *Checker) Can(ctx context.Context, userID int, classroomID int, perm Permission) (bool, error) {
    // Check cache first
    cacheKey := cache.PermissionKey(userID, classroomID)
    if cached, err := c.cache.Get(ctx, cacheKey); err == nil {
        var perms map[Permission]bool
        if err := json.Unmarshal(cached, &perms); err == nil {
            if allowed, ok := perms[perm]; ok {
                return allowed, nil
            }
        }
    }

    // Check if user is classroom owner
    classroom, err := c.repo.GetClassroom(ctx, classroomID)
    if err != nil {
        return false, err
    }

    if classroom.OwnerID == userID {
        // Owner has all permissions
        c.cachePermissions(ctx, userID, classroomID, allPermissions())
        return true, nil
    }

    // Check if user is organization admin
    isAdmin, err := c.forgejoClient.IsOrganizationAdmin(ctx, classroom.OrganizationID, userID)
    if err != nil {
        return false, err
    }

    if isAdmin {
        // Admins have most permissions except deletion
        allowed := perm != PermManageClassroom
        c.cachePermissions(ctx, userID, classroomID, adminPermissions())
        return allowed, nil
    }

    // Check if user is student in roster
    isStudent, err := c.repo.IsStudentInRoster(ctx, classroomID, userID)
    if err != nil {
        return false, err
    }

    if isStudent {
        // Students can only view and accept
        allowed := perm == PermViewClassroom || perm == PermAcceptAssignment
        c.cachePermissions(ctx, userID, classroomID, studentPermissions())
        return allowed, nil
    }

    return false, nil
}

func (c *Checker) cachePermissions(ctx context.Context, userID, classroomID int, perms map[Permission]bool) {
    cacheKey := cache.PermissionKey(userID, classroomID)
    data, _ := json.Marshal(perms)
    c.cache.Set(ctx, cacheKey, data, 15*time.Minute)
}

// Middleware for API
func RequirePermission(perm Permission) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetInt("user_id")
        classroomID := c.GetInt("classroom_id")

        checker := c.MustGet("permission_checker").(*Checker)

        can, err := checker.Can(c.Request.Context(), userID, classroomID, perm)
        if err != nil {
            RespondError(c, "SYSTEM_INTERNAL_ERROR", "Permission check failed", nil)
            return
        }

        if !can {
            RespondError(c, "AUTHZ_INSUFFICIENT_PERMISSIONS",
                "You do not have permission to perform this action", nil)
            return
        }

        c.Next()
    }
}
```

### **11.2 Input Validation and Sanitization**

**Validation Rules:**

```go
// internal/util/validator.go
package util

import (
    "fmt"
    "html"
    "regexp"
    "strings"
    "unicode/utf8"
)

var (
    slugRegex       = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)
    identifierRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
    emailRegex      = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

type ValidationError struct {
    Field   string
    Message string
    Code    string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidateClassroomName validates classroom name
func ValidateClassroomName(name string) *ValidationError {
    name = strings.TrimSpace(name)

    if name == "" {
        return &ValidationError{
            Field:   "name",
            Message: "cannot be empty",
            Code:    "VALIDATION_MISSING_REQUIRED_FIELD",
        }
    }

    if utf8.RuneCountInString(name) > 255 {
        return &ValidationError{
            Field:   "name",
            Message: "cannot exceed 255 characters",
            Code:    "VALIDATION_TOO_LONG",
        }
    }

    return nil
}

// ValidateSlug validates URL-safe slugs
func ValidateSlug(slug string) *ValidationError {
    if len(slug) < 1 || len(slug) > 100 {
        return &ValidationError{
            Field:   "slug",
            Message: "must be between 1 and 100 characters",
            Code:    "VALIDATION_OUT_OF_RANGE",
        }
    }

    if !slugRegex.MatchString(slug) {
        return &ValidationError{
            Field:   "slug",
            Message: "must contain only lowercase letters, numbers, and hyphens",
            Code:    "VALIDATION_INVALID_FORMAT",
        }
    }

    return nil
}

// ValidateDeadline validates assignment deadline
func ValidateDeadline(deadline time.Time) *ValidationError {
    if deadline.IsZero() {
        return nil // Deadline is optional
    }

    if deadline.Before(time.Now()) {
        return &ValidationError{
            Field:   "deadline",
            Message: "must be in the future",
            Code:    "VALIDATION_INVALID_DATE",
        }
    }

    // Reasonable limit: 2 years in the future
    if deadline.After(time.Now().Add(2 * 365 * 24 * time.Hour)) {
        return &ValidationError{
            Field:   "deadline",
            Message: "cannot be more than 2 years in the future",
            Code:    "VALIDATION_INVALID_DATE",
        }
    }

    return nil
}

// ValidateEmail validates email format
func ValidateEmail(email string) *ValidationError {
    email = strings.TrimSpace(strings.ToLower(email))

    if !emailRegex.MatchString(email) {
        return &ValidationError{
            Field:   "email",
            Message: "invalid email format",
            Code:    "VALIDATION_INVALID_FORMAT",
        }
    }

    return nil
}

// SanitizeHTML removes potentially dangerous HTML
func SanitizeHTML(input string) string {
    // Escape HTML entities
    return html.EscapeString(input)
}

// SanitizeMarkdown allows safe markdown but escapes dangerous HTML
func SanitizeMarkdown(input string) string {
    // Use bluemonday for production
    // For now, escape all HTML
    return html.EscapeString(input)
}
```

**API Request Validation:**

```go
// internal/api/v1/assignment.go
package v1

import (
    "github.com/gin-gonic/gin"
    "code.forgejo.org/forgejo/classroom/internal/api"
    "code.forgejo.org/forgejo/classroom/internal/util"
)

type CreateAssignmentRequest struct {
    Title        string    `json:"title" binding:"required,min=1,max=255"`
    Slug         string    `json:"slug" binding:"required,min=1,max=100"`
    Instructions string    `json:"instructions" binding:"max=10000"`
    TemplateRepo string    `json:"template_repo" binding:"required"`
    Type         string    `json:"type" binding:"required,oneof=individual team"`
    Deadline     time.Time `json:"deadline"`
    MaxTeamSize  *int      `json:"max_team_size" binding:"omitempty,min=2,max=10"`
}

func (h *AssignmentHandler) CreateAssignment(c *gin.Context) {
    classroomID := c.GetInt("classroom_id")

    var req CreateAssignmentRequest

    // Gin automatically validates struct tags
    if err := c.ShouldBindJSON(&req); err != nil {
        fields := parseValidationErrors(err)
        api.RespondValidationError(c, fields)
        return
    }

    // Additional custom validation
    if err := util.ValidateSlug(req.Slug); err != nil {
        api.RespondError(c, err.Code, err.Message, map[string]interface{}{
            "field": err.Field,
        })
        return
    }

    if err := util.ValidateDeadline(req.Deadline); err != nil {
        api.RespondError(c, err.Code, err.Message, map[string]interface{}{
            "field": err.Field,
        })
        return
    }

    // Sanitize user input
    req.Instructions = util.SanitizeMarkdown(req.Instructions)

    // Validate team-specific requirements
    if req.Type == "team" {
        if req.MaxTeamSize == nil {
            api.RespondError(c, "VALIDATION_MISSING_REQUIRED_FIELD",
                "max_team_size is required for team assignments", map[string]interface{}{
                    "field": "max_team_size",
                })
            return
        }
    }

    // Continue with business logic
    assignment, err := h.service.Create(c.Request.Context(), classroomID, req)
    if err != nil {
        handleServiceError(c, err)
        return
    }

    c.JSON(http.StatusCreated, assignment)
}

func parseValidationErrors(err error) []api.ValidationError {
    var validationErrors []api.ValidationError

    if validationErrs, ok := err.(validator.ValidationErrors); ok {
        for _, e := range validationErrs {
            validationErrors = append(validationErrors, api.ValidationError{
                Field:   e.Field(),
                Code:    fmt.Sprintf("VALIDATION_%s", strings.ToUpper(e.Tag())),
                Message: formatValidationMessage(e),
            })
        }
    }

    return validationErrors
}
```

### **11.3 SQL Injection Prevention**

**Always use parameterized queries:**

```go
// internal/repository/assignment.go
package repository

func (r *AssignmentRepository) List(ctx context.Context, classroomID int, filters ListFilters) ([]*model.Assignment, error) {
    query := `
        SELECT id, classroom_id, title, slug, type, deadline, created_at, updated_at
        FROM assignments
        WHERE classroom_id = $1
    `

    args := []interface{}{classroomID}
    paramCount := 1

    // Dynamic filtering with parameterized queries
    if filters.Status != "" {
        paramCount++
        query += fmt.Sprintf(" AND status = $%d", paramCount)
        args = append(args, filters.Status)
    }

    if !filters.DeadlineBefore.IsZero() {
        paramCount++
        query += fmt.Sprintf(" AND deadline < $%d", paramCount)
        args = append(args, filters.DeadlineBefore)
    }

    if filters.Type != "" {
        paramCount++
        query += fmt.Sprintf(" AND type = $%d", paramCount)
        args = append(args, filters.Type)
    }

    // Add sorting
    switch filters.SortBy {
    case "created_at":
        query += " ORDER BY created_at"
    case "deadline":
        query += " ORDER BY deadline"
    case "title":
        query += " ORDER BY title"
    default:
        query += " ORDER BY created_at"
    }

    if filters.SortOrder == "asc" {
        query += " ASC"
    } else {
        query += " DESC"
    }

    // Add pagination
    paramCount++
    query += fmt.Sprintf(" LIMIT $%d", paramCount)
    args = append(args, filters.Limit)

    paramCount++
    query += fmt.Sprintf(" OFFSET $%d", paramCount)
    args = append(args, filters.Offset)

    // ALWAYS use parameterized queries - NEVER string concatenation
    rows, err := r.db.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var assignments []*model.Assignment
    for rows.Next() {
        var a model.Assignment
        if err := rows.Scan(
            &a.ID, &a.ClassroomID, &a.Title, &a.Slug,
            &a.Type, &a.Deadline, &a.CreatedAt, &a.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        assignments = append(assignments, &a)
    }

    return assignments, rows.Err()
}
```

### **11.4 Secure Repository Operations**

**Repository Creation with Proper Security:**

```go
// internal/service/submission.go
package service

import (
    "context"
    "fmt"
    "regexp"
    "time"
)

func (s *SubmissionService) AcceptAssignment(ctx context.Context, assignmentID int, studentIdentifier string) (*model.Submission, error) {
    // Get authenticated user from context
    userID := ctx.Value("user_id").(int)

    // Validate student is in roster
    student, err := s.rosterRepo.GetByIdentifier(ctx, assignmentID, studentIdentifier)
    if err != nil {
        return nil, fmt.Errorf("BUSINESS_ROSTER_NOT_FOUND: student not found in roster")
    }

    // Verify authenticated user matches student
    if student.ForgejoUserID != userID {
        return nil, fmt.Errorf("AUTHZ_FORBIDDEN: cannot accept assignment for another student")
    }

    // Check if already accepted (idempotent operation)
    existing, _ := s.submissionRepo.GetByAssignmentAndStudent(ctx, assignmentID, student.ID)
    if existing != nil {
        return existing, nil
    }

    assignment, err := s.assignmentRepo.Get(ctx, assignmentID)
    if err != nil {
        return nil, err
    }

    // Check deadline if enforced
    if !assignment.AllowLateSubmissions && !assignment.Deadline.IsZero() {
        if time.Now().After(assignment.Deadline) {
            return nil, fmt.Errorf("BUSINESS_DEADLINE_PASSED: assignment deadline has passed")
        }
    }

    // Generate safe repository name
    repoName := s.generateSafeRepoName(assignment.Slug, student)

    // Create repository from template with proper permissions
    var submission *model.Submission
    err = s.repo.WithTransaction(ctx, func(tx *sql.Tx) error {
        repo, err := s.forgejoClient.CreateRepositoryFromTemplate(ctx, forgejo.CreateRepoRequest{
            Owner:        assignment.Classroom.OrganizationName,
            Name:         repoName,
            TemplateRepo: assignment.TemplateRepoID,
            Private:      assignment.Visibility == "private", // Enforce privacy
            Description:  fmt.Sprintf("Assignment: %s - Student: %s", assignment.Title, student.FullName),
        })
        if err != nil {
            return fmt.Errorf("failed to create repository: %w", err)
        }

        // Grant student WRITE access (NOT admin - security boundary)
        err = s.forgejoClient.AddCollaborator(ctx, repo.ID, student.ForgejoUserID, forgejo.PermissionWrite)
        if err != nil {
            // Rollback: delete repository
            s.forgejoClient.DeleteRepository(ctx, repo.ID)
            return fmt.Errorf("failed to grant access: %w", err)
        }

        // Protect default branch if configured
        if assignment.ProtectDefaultBranch {
            err = s.forgejoClient.ProtectBranch(ctx, repo.ID, "main", forgejo.BranchProtection{
                EnablePush:             true,
                EnablePushWhitelist:    true,
                PushWhitelistUsernames: []string{student.ForgejoUsername},
            })
            if err != nil {
                // Log but don't fail - branch protection is optional
                log.Warn("Failed to protect branch", "repo", repo.ID, "error", err)
            }
        }

        // Create submission record
        submission = &model.Submission{
            AssignmentID:   assignmentID,
            StudentID:      student.ID,
            RepositoryID:   repo.ID,
            RepositoryName: repo.FullName,
            RepositoryURL:  repo.CloneURL,
            AcceptedAt:     time.Now(),
            Status:         "in_progress",
        }

        if err := s.submissionRepo.CreateTx(ctx, tx, submission); err != nil {
            // Rollback: delete repository
            s.forgejoClient.DeleteRepository(ctx, repo.ID)
            return err
        }

        return nil
    })

    if err != nil {
        return nil, err
    }

    // Invalidate caches
    s.cache.DeletePattern(ctx, fmt.Sprintf("assignment:%d:submissions*", assignmentID))
    s.cache.DeletePattern(ctx, fmt.Sprintf("stats:assignment:%d", assignmentID))

    return submission, nil
}

func (s *SubmissionService) generateSafeRepoName(assignmentSlug string, student *model.RosterEntry) string {
    // Format: assignment-slug-username
    username := student.ForgejoUsername
    if username == "" {
        username = student.Identifier
    }

    // Sanitize: remove special chars, enforce max length
    safeRegex := regexp.MustCompile(`[^a-zA-Z0-9-_]`)
    username = safeRegex.ReplaceAllString(username, "-")

    if len(username) > 50 {
        username = username[:50]
    }

    // Ensure doesn't start/end with hyphen
    username = strings.Trim(username, "-")

    return fmt.Sprintf("%s-%s", assignmentSlug, username)
}
```

### **11.5 Secrets Management**

**Configuration Security:**

```yaml
# config.yaml - DO NOT commit secrets
server:
  host: 0.0.0.0
  port: 8080
  read_timeout: 30s
  write_timeout: 30s

forgejo:
  url: https://forgejo.example.com
  # Token loaded from environment variable
  token: ${FORGEJO_API_TOKEN}

database:
  host: localhost
  port: 5432
  name: forgejo_classroom
  user: classroom_user
  # Password loaded from environment variable
  password: ${DB_PASSWORD}
  sslmode: require
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 1h

redis:
  host: localhost
  port: 6379
  # Password loaded from environment variable
  password: ${REDIS_PASSWORD}
  db: 0

security:
  # Secret key for signing tokens
  secret_key: ${SECRET_KEY}

rate_limit:
  enabled: true
  redis_enabled: true
```

**Environment Variable Loading:**

```go
// internal/config/config.go
package config

import (
    "errors"
    "os"
    "time"

    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Forgejo  ForgejoConfig  `mapstructure:"forgejo"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    Security SecurityConfig `mapstructure:"security"`
    RateLimit RateLimitConfig `mapstructure:"rate_limit"`
}

type ServerConfig struct {
    Host         string        `mapstructure:"host"`
    Port         int           `mapstructure:"port"`
    ReadTimeout  time.Duration `mapstructure:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type ForgejoConfig struct {
    URL   string `mapstructure:"url"`
    Token string `mapstructure:"token"`
}

type DatabaseConfig struct {
    Host            string        `mapstructure:"host"`
    Port            int           `mapstructure:"port"`
    Name            string        `mapstructure:"name"`
    User            string        `mapstructure:"user"`
    Password        string        `mapstructure:"password"`
    SSLMode         string        `mapstructure:"sslmode"`
    MaxOpenConns    int           `mapstructure:"max_open_conns"`
    MaxIdleConns    int           `mapstructure:"max_idle_conns"`
    ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Password string `mapstructure:"password"`
    DB       int    `mapstructure:"db"`
}

type SecurityConfig struct {
    SecretKey string `mapstructure:"secret_key"`
}

type RateLimitConfig struct {
    Enabled      bool `mapstructure:"enabled"`
    RedisEnabled bool `mapstructure:"redis_enabled"`
}

func Load(path string) (*Config, error) {
    viper.SetConfigFile(path)

    // Enable environment variable substitution
    viper.AutomaticEnv()
    viper.SetEnvPrefix("FGC")

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, err
    }

    // Validate required secrets are set
    if cfg.Forgejo.Token == "" {
        return nil, errors.New("FORGEJO_API_TOKEN not set")
    }

    if cfg.Database.Password == "" {
        return nil, errors.New("DB_PASSWORD not set")
    }

    if cfg.Security.SecretKey == "" {
        return nil, errors.New("SECRET_KEY not set")
    }

    // Set defaults
    if cfg.Server.Port == 0 {
        cfg.Server.Port = 8080
    }

    if cfg.Database.MaxOpenConns == 0 {
        cfg.Database.MaxOpenConns = 25
    }

    if cfg.Database.MaxIdleConns == 0 {
        cfg.Database.MaxIdleConns = 5
    }

    if cfg.Database.ConnMaxLifetime == 0 {
        cfg.Database.ConnMaxLifetime = time.Hour
    }

    return &cfg, nil
}

func (c *DatabaseConfig) ConnectionString() string {
    return fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
    )
}
```

### **11.6 Audit Logging**

**Security Event Logging:**

```go
// internal/audit/logger.go
package audit

import (
    "context"
    "encoding/json"
    "io"
    "time"

    "go.uber.org/zap"
)

type EventType string

const (
    EventClassroomCreated     EventType = "classroom.created"
    EventClassroomDeleted     EventType = "classroom.deleted"
    EventAssignmentCreated    EventType = "assignment.created"
    EventAssignmentAccepted   EventType = "assignment.accepted"
    EventRosterModified       EventType = "roster.modified"
    EventPermissionDenied     EventType = "permission.denied"
    EventAuthenticationFailed EventType = "auth.failed"
    EventRateLimitExceeded    EventType = "rate_limit.exceeded"
)

type Event struct {
    Type      EventType              `json:"type"`
    UserID    int                    `json:"user_id"`
    IP        string                 `json:"ip"`
    Timestamp time.Time              `json:"timestamp"`
    Resource  string                 `json:"resource"`
    Action    string                 `json:"action"`
    Details   map[string]interface{} `json:"details"`
    Success   bool                   `json:"success"`
}

type Logger struct {
    writer io.Writer
    zap    *zap.Logger
}

func NewLogger(writer io.Writer, zapLogger *zap.Logger) *Logger {
    return &Logger{
        writer: writer,
        zap:    zapLogger,
    }
}

func (l *Logger) Log(ctx context.Context, event Event) {
    event.Timestamp = time.Now()

    jsonData, err := json.Marshal(event)
    if err != nil {
        l.zap.Error("Failed to marshal audit event", zap.Error(err))
        return
    }

    l.writer.Write(jsonData)
    l.writer.Write([]byte("\n"))

    // Also log to structured logger for searchability
    l.zap.Info("audit_event",
        zap.String("type", string(event.Type)),
        zap.Int("user_id", event.UserID),
        zap.String("ip", event.IP),
        zap.String("resource", event.Resource),
        zap.String("action", event.Action),
        zap.Bool("success", event.Success),
    )
}

// Middleware to log API requests
func AuditMiddleware(logger *Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        c.Next()

        // Log security-relevant events
        if c.Writer.Status() == 401 {
            logger.Log(c.Request.Context(), Event{
                Type:     EventAuthenticationFailed,
                IP:       c.ClientIP(),
                Resource: c.Request.URL.Path,
                Action:   c.Request.Method,
                Success:  false,
                Details: map[string]interface{}{
                    "status":   401,
                    "duration": time.Since(start).Milliseconds(),
                },
            })
        }

        if c.Writer.Status() == 403 {
            logger.Log(c.Request.Context(), Event{
                Type:     EventPermissionDenied,
                UserID:   c.GetInt("user_id"),
                IP:       c.ClientIP(),
                Resource: c.Request.URL.Path,
                Action:   c.Request.Method,
                Success:  false,
                Details: map[string]interface{}{
                    "status":   403,
                    "duration": time.Since(start).Milliseconds(),
                },
            })
        }

        if c.Writer.Status() == 429 {
            logger.Log(c.Request.Context(), Event{
                Type:     EventRateLimitExceeded,
                UserID:   c.GetInt("user_id"),
                IP:       c.ClientIP(),
                Resource: c.Request.URL.Path,
                Action:   c.Request.Method,
                Success:  false,
                Details: map[string]interface{}{
                    "status":   429,
                    "duration": time.Since(start).Milliseconds(),
                },
            })
        }
    }
}
```

---

## **12. Deployment Guide**

### **12.1 System Requirements**

**Minimum Requirements:**
- OS: Linux (Ubuntu 20.04+, Debian 11+, CentOS 8+)
- CPU: 2 cores
- RAM: 4 GB
- Disk: 20 GB
- Database: PostgreSQL 12+
- Cache: Redis 6+
- Forgejo: v1.20+ (or Gitea v1.20+)

**Recommended for Production:**
- CPU: 4+ cores
- RAM: 8+ GB
- Disk: 100 GB SSD
- Database: PostgreSQL 14+ (separate server)
- Redis: Dedicated instance
- Load balancer (nginx/haproxy)

### **12.2 Installation Methods**

**Method 1: Binary Installation**

```bash
# Download latest release
wget https://github.com/forgejo/classroom/releases/download/v0.1.0/fgc-linux-amd64
chmod +x fgc-linux-amd64
sudo mv fgc-linux-amd64 /usr/local/bin/fgc

# Download server
wget https://github.com/forgejo/classroom/releases/download/v0.1.0/fgc-server-linux-amd64
chmod +x fgc-server-linux-amd64
sudo mv fgc-server-linux-amd64 /usr/local/bin/fgc-server

# Create user
sudo useradd -r -s /bin/false classroom

# Create directories
sudo mkdir -p /etc/forgejo-classroom
sudo mkdir -p /var/log/forgejo-classroom
sudo chown classroom:classroom /var/log/forgejo-classroom

# Create config file
sudo tee /etc/forgejo-classroom/config.yaml > /dev/null <<EOF
server:
  host: 0.0.0.0
  port: 8080

forgejo:
  url: https://forgejo.example.com
  token: \${FORGEJO_API_TOKEN}

database:
  host: localhost
  port: 5432
  name: forgejo_classroom
  user: classroom_user
  password: \${DB_PASSWORD}
  sslmode: require

redis:
  host: localhost
  port: 6379
  password: \${REDIS_PASSWORD}

security:
  secret_key: \${SECRET_KEY}
EOF

# Set environment variables
sudo tee /etc/forgejo-classroom/env > /dev/null <<EOF
FORGEJO_API_TOKEN=your_token_here
DB_PASSWORD=your_db_password
REDIS_PASSWORD=your_redis_password
SECRET_KEY=your_secret_key
EOF

sudo chmod 600 /etc/forgejo-classroom/env

# Run migrations
fgc-server migrate up --config /etc/forgejo-classroom/config.yaml

# Start server (see systemd section)
```

**Method 2: Docker Deployment**

```yaml
# docker-compose.yml
version: '3.8'

services:
  postgres:
    image: postgres:14
    environment:
      POSTGRES_DB: forgejo_classroom
      POSTGRES_USER: classroom_user
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - classroom-net
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U classroom_user"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis_data:/data
    networks:
      - classroom-net
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  classroom-server:
    image: forgejo/classroom:latest
    ports:
      - "8080:8080"
    environment:
      FGC_FORGEJO_URL: https://forgejo.example.com
      FGC_FORGEJO_TOKEN: ${FORGEJO_API_TOKEN}
      FGC_DATABASE_HOST: postgres
      FGC_DATABASE_NAME: forgejo_classroom
      FGC_DATABASE_USER: classroom_user
      FGC_DATABASE_PASSWORD: ${DB_PASSWORD}
      FGC_REDIS_HOST: redis
      FGC_REDIS_PASSWORD: ${REDIS_PASSWORD}
      FGC_SECRET_KEY: ${SECRET_KEY}
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - classroom-net
    volumes:
      - ./config.yaml:/etc/forgejo-classroom/config.yaml:ro
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:

networks:
  classroom-net:
```

```bash
# Start services
docker-compose up -d

# Run migrations
docker-compose exec classroom-server fgc-server migrate up

# View logs
docker-compose logs -f classroom-server

# Stop services
docker-compose down
```

### **12.3 Nginx Reverse Proxy Configuration**

```nginx
# /etc/nginx/sites-available/forgejo-classroom
upstream classroom_backend {
    server 127.0.0.1:8080;
    keepalive 32;
}

server {
    listen 80;
    server_name classroom.example.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name classroom.example.com;

    ssl_certificate /etc/letsencrypt/live/classroom.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/classroom.example.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    client_max_body_size 100M;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;

    location / {
        proxy_pass http://classroom_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket support (future)
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";

        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;

        # Buffering
        proxy_buffering on;
        proxy_buffer_size 4k;
        proxy_buffers 8 4k;
    }

    # Health check endpoint
    location /health {
        proxy_pass http://classroom_backend;
        access_log off;
    }

    # Rate limiting (additional layer)
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;

    location /api/ {
        limit_req zone=api burst=20 nodelay;
        proxy_pass http://classroom_backend;
    }
}
```

### **12.4 Systemd Service**

```ini
# /etc/systemd/system/forgejo-classroom.service
[Unit]
Description=Forgejo Classroom Server
After=network.target postgresql.service redis.service
Wants=postgresql.service redis.service

[Service]
Type=simple
User=classroom
Group=classroom
WorkingDirectory=/opt/forgejo-classroom
ExecStart=/usr/local/bin/fgc-server serve --config /etc/forgejo-classroom/config.yaml
Restart=always
RestartSec=10

# Load environment variables
EnvironmentFile=/etc/forgejo-classroom/env

# Security hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/log/forgejo-classroom
CapabilityBoundingSet=
AmbientCapabilities=
SystemCallFilter=@system-service
SystemCallErrorNumber=EPERM

# Resource limits
LimitNOFILE=65536
LimitNPROC=512

# Logging
StandardOutput=append:/var/log/forgejo-classroom/server.log
StandardError=append:/var/log/forgejo-classroom/error.log

[Install]
WantedBy=multi-user.target
```

```bash
# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable forgejo-classroom
sudo systemctl start forgejo-classroom
sudo systemctl status forgejo-classroom

# View logs
sudo journalctl -u forgejo-classroom -f

# Restart service
sudo systemctl restart forgejo-classroom

# Stop service
sudo systemctl stop forgejo-classroom
```

---

## **13. Monitoring and Observability**

### **13.1 Health Check Endpoints**

```go
// internal/api/health.go
package api

import (
    "context"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

type HealthHandler struct {
    db            *sql.DB
    cache         cache.Cache
    forgejoClient *forgejo.Client
}

func (h *HealthHandler) Health(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
    defer cancel()

    health := map[string]interface{}{
        "status":    "healthy",
        "timestamp": time.Now(),
        "checks":    map[string]interface{}{},
    }

    // Database check
    if err := h.db.PingContext(ctx); err != nil {
        health["status"] = "unhealthy"
        health["checks"].(map[string]interface{})["database"] = map[string]interface{}{
            "status": "down",
            "error":  err.Error(),
        }
    } else {
        var result int
        err := h.db.QueryRowContext(ctx, "SELECT 1").Scan(&result)
        if err != nil {
            health["status"] = "degraded"
            health["checks"].(map[string]interface{})["database"] = map[string]interface{}{
                "status": "degraded",
                "error":  err.Error(),
            }
        } else {
            health["checks"].(map[string]interface{})["database"] = map[string]interface{}{
                "status": "up",
            }
        }
    }

    // Redis/Cache check
    if err := h.cache.Ping(ctx); err != nil {
        health["status"] = "degraded"
        health["checks"].(map[string]interface{})["cache"] = map[string]interface{}{
            "status": "down",
            "error":  err.Error(),
        }
    } else {
        health["checks"].(map[string]interface{})["cache"] = map[string]interface{}{
            "status": "up",
        }
    }

    // Forgejo API check
    if _, err := h.forgejoClient.GetVersion(ctx); err != nil {
        health["status"] = "degraded"
        health["checks"].(map[string]interface{})["forgejo"] = map[string]interface{}{
            "status": "down",
            "error":  err.Error(),
        }
    } else {
        health["checks"].(map[string]interface{})["forgejo"] = map[string]interface{}{
            "status": "up",
        }
    }

    statusCode := http.StatusOK
    if health["status"] == "unhealthy" {
        statusCode = http.StatusServiceUnavailable
    } else if health["status"] == "degraded" {
        statusCode = http.StatusOK // Still serve traffic
    }

    c.JSON(statusCode, health)
}

func (h *HealthHandler) Ready(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
    defer cancel()

    // Check if service is ready to accept traffic
    if err := h.db.PingContext(ctx); err != nil {
        c.JSON(http.StatusServiceUnavailable, gin.H{
            "ready": false,
            "error": "database not ready",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "ready": true,
    })
}

func (h *HealthHandler) Live(c *gin.Context) {
    // Simple liveness check
    c.JSON(http.StatusOK, gin.H{
        "alive": true,
    })
}
```

### **13.2 Prometheus Metrics**

```go
// internal/metrics/metrics.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    HTTPRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "fgc_http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )

    HTTPRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "fgc_http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
        },
        []string{"method", "path"},
    )

    AssignmentsCreated = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "fgc_assignments_created_total",
            Help: "Total number of assignments created",
        },
    )

    AssignmentsAccepted = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "fgc_assignments_accepted_total",
            Help: "Total number of assignments accepted by students",
        },
    )

    SubmissionsTotal = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "fgc_submissions_total",
            Help: "Total number of submissions by status",
        },
        []string{"status"},
    )

    RepositoryCreationDuration = promauto.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "fgc_repository_creation_duration_seconds",
            Help:    "Duration of repository creation operations",
            Buckets: []float64{0.5, 1, 2, 5, 10, 30, 60},
        },
    )

    DatabaseQueries = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "fgc_database_queries_total",
            Help: "Total number of database queries",
        },
        []string{"operation", "table"},
    )

    DatabaseQueryDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "fgc_database_query_duration_seconds",
            Help:    "Database query duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"operation", "table"},
    )

    CacheOperations = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "fgc_cache_operations_total",
            Help: "Total number of cache operations",
        },
        []string{"operation", "result"}, // operation: get/set/delete, result: hit/miss/error
    )

    ForgejoAPIRequests = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "fgc_forgejo_api_requests_total",
            Help: "Total number of Forgejo API requests",
        },
        []string{"endpoint", "status"},
    )

    QueueJobsProcessed = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "fgc_queue_jobs_processed_total",
            Help: "Total number of queue jobs processed",
        },
        []string{"queue", "status"}, // status: success/failure
    )

    ActiveConnections = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "fgc_active_connections",
            Help: "Number of active HTTP connections",
        },
    )
)

// Middleware to record HTTP metrics
func MetricsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        ActiveConnections.Inc()
        defer ActiveConnections.Dec()

        c.Next()

        duration := time.Since(start).Seconds()
        status := strconv.Itoa(c.Writer.Status())

        HTTPRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), status).Inc()
        HTTPRequestDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration)
    }
}

// Database metrics wrapper
func RecordDatabaseQuery(operation, table string, duration time.Duration) {
    DatabaseQueries.WithLabelValues(operation, table).Inc()
    DatabaseQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}

// Cache metrics wrapper
func RecordCacheOperation(operation, result string) {
    CacheOperations.WithLabelValues(operation, result).Inc()
}

// Forgejo API metrics wrapper
func RecordForgejoRequest(endpoint, status string) {
    ForgejoAPIRequests.WithLabelValues(endpoint, status).Inc()
}

// Queue job metrics wrapper
func RecordQueueJob(queue, status string) {
    QueueJobsProcessed.WithLabelValues(queue, status).Inc()
}
```

**Expose Metrics Endpoint:**

```go
// cmd/fgc-server/main.go
import (
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func setupRouter() *gin.Engine {
    r := gin.Default()

    // Add metrics middleware
    r.Use(metrics.MetricsMiddleware())

    // Expose metrics endpoint
    r.GET("/metrics", gin.WrapH(promhttp.Handler()))

    // Health checks
    r.GET("/health", healthHandler.Health)
    r.GET("/ready", healthHandler.Ready)
    r.GET("/live", healthHandler.Live)

    // API routes
    v1 := r.Group("/api/classroom/v1")
    {
        // ... routes
    }

    return r
}
```

### **13.3 Grafana Dashboard**

**Prometheus Configuration:**

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'forgejo-classroom'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
```

**Grafana Dashboard JSON:**

```json
{
  "dashboard": {
    "title": "Forgejo Classroom Monitoring",
    "timezone": "browser",
    "panels": [
      {
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(fgc_http_requests_total[5m])",
            "legendFormat": "{{method}} {{path}} - {{status}}"
          }
        ],
        "yaxes": [
          {
            "format": "reqps",
            "label": "Requests/sec"
          }
        ]
      },
      {
        "title": "Request Latency (P95)",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(fgc_http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "{{method}} {{path}}"
          }
        ],
        "yaxes": [
          {
            "format": "s",
            "label": "Duration"
          }
        ]
      },
      {
        "title": "Error Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(fgc_http_requests_total{status=~\"5..\"}[5m])",
            "legendFormat": "5xx errors"
          },
          {
            "expr": "rate(fgc_http_requests_total{status=~\"4..\"}[5m])",
            "legendFormat": "4xx errors"
          }
        ]
      },
      {
        "title": "Assignment Activity",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(fgc_assignments_created_total[1h])",
            "legendFormat": "Assignments Created"
          },
          {
            "expr": "rate(fgc_assignments_accepted_total[1h])",
            "legendFormat": "Assignments Accepted"
          }
        ]
      },
      {
        "title": "Repository Creation Time",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(fgc_repository_creation_duration_seconds_bucket[5m]))",
            "legendFormat": "P95"
          },
          {
            "expr": "histogram_quantile(0.99, rate(fgc_repository_creation_duration_seconds_bucket[5m]))",
            "legendFormat": "P99"
          }
        ]
      },
      {
        "title": "Cache Hit Rate",
        "type": "singlestat",
        "targets": [
          {
            "expr": "rate(fgc_cache_operations_total{result=\"hit\"}[5m]) / rate(fgc_cache_operations_total{operation=\"get\"}[5m])"
          }
        ],
        "format": "percentunit",
        "thresholds": "0.7,0.8"
      },
      {
        "title": "Database Query Duration",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(fgc_database_query_duration_seconds_bucket[5m]))",
            "legendFormat": "{{operation}} {{table}}"
          }
        ]
      },
      {
        "title": "Active Connections",
        "type": "graph",
        "targets": [
          {
            "expr": "fgc_active_connections",
            "legendFormat": "Active"
          }
        ]
      },
      {
        "title": "Queue Job Success Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(fgc_queue_jobs_processed_total{status=\"success\"}[5m]) / rate(fgc_queue_jobs_processed_total[5m])",
            "legendFormat": "{{queue}}"
          }
        ],
        "yaxes": [
          {
            "format": "percentunit"
          }
        ]
      },
      {
        "title": "Forgejo API Health",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(fgc_forgejo_api_requests_total{status=~\"2..\"}[5m])",
            "legendFormat": "Success"
          },
          {
            "expr": "rate(fgc_forgejo_api_requests_total{status=~\"[45]..\"}[5m])",
            "legendFormat": "Errors"
          }
        ]
      }
    ]
  }
}
```

### **13.4 Logging Strategy**

```go
// internal/logging/logger.go
package logging

import (
    "os"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func NewLogger(level string, development bool) (*zap.Logger, error) {
    var config zap.Config

    if development {
        config = zap.NewDevelopmentConfig()
        config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
    } else {
        config = zap.NewProductionConfig()
        config.EncoderConfig.TimeKey = "timestamp"
        config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    }

    // Parse level
    var zapLevel zapcore.Level
    if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
        zapLevel = zapcore.InfoLevel
    }
    config.Level = zap.NewAtomicLevelAt(zapLevel)

    // Add caller information
    config.EncoderConfig.CallerKey = "caller"
    config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

    return config.Build()
}

// Context keys for request tracking
type contextKey string

const (
    RequestIDKey contextKey = "request_id"
    UserIDKey    contextKey = "user_id"
)

// Middleware to add request ID and structured logging
func RequestLoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        requestID := uuid.New().String()

        c.Set("request_id", requestID)
        c.Set("logger", logger.With(
            zap.String("request_id", requestID),
            zap.String("method", c.Request.Method),
            zap.String("path", c.Request.URL.Path),
            zap.String("ip", c.ClientIP()),
        ))

        c.Next()

        duration := time.Since(start)

        loggerWithFields := logger.With(
            zap.String("request_id", requestID),
            zap.String("method", c.Request.Method),
            zap.String("path", c.Request.URL.Path),
            zap.Int("status", c.Writer.Status()),
            zap.Duration("duration", duration),
            zap.String("ip", c.ClientIP()),
            zap.Int("user_id", c.GetInt("user_id")),
        )

        if len(c.Errors) > 0 {
            loggerWithFields.Error("Request completed with errors",
                zap.String("errors", c.Errors.String()),
            )
        } else if c.Writer.Status() >= 500 {
            loggerWithFields.Error("Request failed")
        } else if c.Writer.Status() >= 400 {
            loggerWithFields.Warn("Request error")
        } else {
            loggerWithFields.Info("Request completed")
        }
    }
}
```

### **13.5 Alerting Rules**

**Prometheus Alert Rules:**

```yaml
# alerts.yml
groups:
  - name: forgejo_classroom
    interval: 30s
    rules:
      # High error rate
      - alert: HighErrorRate
        expr: rate(fgc_http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value }} req/s for {{ $labels.path }}"

      # High latency
      - alert: HighLatency
        expr: histogram_quantile(0.95, rate(fgc_http_request_duration_seconds_bucket[5m])) > 1
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "High latency detected"
          description: "P95 latency is {{ $value }}s for {{ $labels.path }}"

      # Database connection issues
      - alert: DatabaseDown
        expr: up{job="forgejo-classroom"} == 0 or fgc_database_queries_total == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Database connectivity issue"
          description: "Cannot connect to database"

      # Cache unavailable
      - alert: CacheDown
        expr: rate(fgc_cache_operations_total{result="error"}[5m]) > 0.1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Cache service degraded"
          description: "Cache error rate is high"

      # Forgejo API issues
      - alert: ForgejoAPIDown
        expr: rate(fgc_forgejo_api_requests_total{status=~"[45].."}[5m]) > 0.5
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Forgejo API connectivity issues"
          description: "Cannot communicate with Forgejo API"

      # Low cache hit rate
      - alert: LowCacheHitRate
        expr: rate(fgc_cache_operations_total{result="hit"}[10m]) / rate(fgc_cache_operations_total{operation="get"}[10m]) < 0.5
        for: 15m
        labels:
          severity: warning
        annotations:
          summary: "Low cache hit rate"
          description: "Cache hit rate is {{ $value | humanizePercentage }}"

      # High queue job failure rate
      - alert: HighQueueFailureRate
        expr: rate(fgc_queue_jobs_processed_total{status="failure"}[10m]) / rate(fgc_queue_jobs_processed_total[10m]) > 0.1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High queue job failure rate"
          description: "Queue {{ $labels.queue }} has high failure rate"
```

---

## **14. API Documentation Standards**

### **14.1 Complete OpenAPI Specification**

```yaml
# docs/api/openapi.yaml
openapi: 3.0.3
info:
  title: Forgejo Classroom API
  version: 1.0.0
  description: |
    Forgejo Classroom API for managing educational assignments and coursework.

    ## Authentication
    All endpoints require authentication using a Forgejo personal access token.
    Include the token in the Authorization header:
    ```
    Authorization: token YOUR_TOKEN_HERE
    ```

    ## Rate Limiting
    API requests are rate limited based on user role:
    - Owners: 1000 reads/min, 300 writes/min, 50 bulk/min
    - Admins: 500 reads/min, 150 writes/min, 30 bulk/min
    - Students: 200 reads/min, 50 writes/min, 5 bulk/min

    Rate limit information is included in response headers:
    - `X-RateLimit-Limit`: Request limit per window
    - `X-RateLimit-Remaining`: Requests remaining
    - `X-RateLimit-Reset`: Unix timestamp when limit resets

    ## Pagination
    List endpoints support pagination using query parameters:
    - `page`: Page number (default: 1)
    - `per_page`: Items per page (default: 30, max: 100)

    Pagination metadata is included in response headers:
    - `X-Total-Count`: Total number of items
    - `X-Page`: Current page number
    - `X-Per-Page`: Items per page
    - `Link`: Pagination links (next, prev, first, last)

    ## Error Handling
    Errors follow a standardized format with error codes, messages, and details.
    See the Error schema for complete structure.

    ## Versioning
    The API uses URL versioning (e.g., `/api/classroom/v1`).
    Breaking changes will result in a new major version.

```yaml
contact:
  name: Forgejo Classroom Team
  url: https://code.forgejo.org/forgejo/classroom
license:
  name: MIT
  url: https://opensource.org/licenses/MIT

servers:
  - url: https://forgejo.example.com/api/classroom/v1
    description: Production server
  - url: http://localhost:8080/api/classroom/v1
    description: Development server

security:
  - TokenAuth: []

components:
  securitySchemes:
    TokenAuth:
      type: apiKey
      in: header
      name: Authorization
      description: "Format: token YOUR_TOKEN_HERE"

  parameters:
    PageParam:
      name: page
      in: query
      description: Page number for pagination
      schema:
        type: integer
        minimum: 1
        default: 1

    PerPageParam:
      name: per_page
      in: query
      description: Number of items per page
      schema:
        type: integer
        minimum: 1
        maximum: 100
        default: 30

    SortParam:
      name: sort
      in: query
      description: Field to sort by
      schema:
        type: string
        enum: [created_at, updated_at, name, deadline, title]
        default: created_at

    OrderParam:
      name: order
      in: query
      description: Sort direction
      schema:
        type: string
        enum: [asc, desc]
        default: desc

  schemas:
    Classroom:
      type: object
      required:
        - id
        - name
        - slug
        - organization_id
        - owner_id
        - status
        - created_at
        - updated_at
      properties:
        id:
          type: integer
          example: 1
        name:
          type: string
          example: "CS101 Fall 2025"
        slug:
          type: string
          example: "cs101-fall2025"
        description:
          type: string
          nullable: true
        organization_id:
          type: integer
          example: 42
        organization_name:
          type: string
          example: "cs101-fall2025"
        owner_id:
          type: integer
          example: 5
        status:
          type: string
          enum: [active, archived]
          example: "active"
        student_count:
          type: integer
          example: 45
        assignment_count:
          type: integer
          example: 8
        created_at:
          type: string
          format: date-time
        updated_at:
          type: string
          format: date-time

    Assignment:
      type: object
      required:
        - id
        - classroom_id
        - title
        - slug
        - type
        - invitation_code
        - created_at
        - updated_at
      properties:
        id:
          type: integer
          example: 10
        classroom_id:
          type: integer
          example: 1
        title:
          type: string
          example: "Homework 1: Variables and Operators"
        slug:
          type: string
          example: "hw01"
        instructions:
          type: string
          nullable: true
        type:
          type: string
          enum: [individual, team]
        template_repo_id:
          type: integer
          example: 100
        template_repo_name:
          type: string
          example: "cs101-templates/hw01-starter"
        deadline:
          type: string
          format: date-time
          nullable: true
        max_team_size:
          type: integer
          nullable: true
        allow_late_submissions:
          type: boolean
          default: true
        visibility:
          type: string
          enum: [private, public]
          default: private
        invitation_code:
          type: string
          example: "abc123xyz"
        invitation_url:
          type: string
          example: "https://forgejo.example.com/classroom/accept/hw01/abc123"
        acceptance_count:
          type: integer
          example: 42
        submission_count:
          type: integer
          example: 38
        created_at:
          type: string
          format: date-time
        updated_at:
          type: string
          format: date-time

    Submission:
      type: object
      required:
        - id
        - assignment_id
        - student_id
        - repository_id
        - status
        - accepted_at
        - created_at
        - updated_at
      properties:
        id:
          type: integer
          example: 42
        assignment_id:
          type: integer
          example: 10
        student_id:
          type: integer
          example: 123
        repository_id:
          type: integer
          example: 500
        repository_name:
          type: string
          example: "cs101-fall2025/hw01-jdoe"
        repository_url:
          type: string
          example: "https://forgejo.example.com/cs101-fall2025/hw01-jdoe"
        status:
          type: string
          enum: [pending, in_progress, submitted, graded]
        accepted_at:
          type: string
          format: date-time
        last_commit_at:
          type: string
          format: date-time
          nullable: true
        last_commit_sha:
          type: string
          nullable: true
        submitted_at:
          type: string
          format: date-time
          nullable: true
        deadline_tag:
          type: string
          nullable: true
        is_late:
          type: boolean
        commit_count:
          type: integer
        grade:
          type: number
          format: float
          nullable: true
        graded_at:
          type: string
          format: date-time
          nullable: true
        graded_by:
          type: integer
          nullable: true
        feedback:
          type: string
          nullable: true
        created_at:
          type: string
          format: date-time
        updated_at:
          type: string
          format: date-time

    RosterEntry:
      type: object
      required:
        - id
        - classroom_id
        - identifier
        - email
        - full_name
        - status
        - created_at
        - updated_at
      properties:
        id:
          type: integer
          example: 10
        classroom_id:
          type: integer
          example: 1
        identifier:
          type: string
          example: "student001"
        email:
          type: string
          format: email
          example: "john.doe@university.edu"
        full_name:
          type: string
          example: "John Doe"
        forgejo_user_id:
          type: integer
          nullable: true
        forgejo_username:
          type: string
          nullable: true
        status:
          type: string
          enum: [pending, linked]
        created_at:
          type: string
          format: date-time
        updated_at:
          type: string
          format: date-time

    Team:
      type: object
      required:
        - id
        - assignment_id
        - name
        - slug
        - member_ids
        - created_at
        - updated_at
      properties:
        id:
          type: integer
          example: 5
        assignment_id:
          type: integer
          example: 10
        name:
          type: string
          example: "Team Alpha"
        slug:
          type: string
          example: "team-alpha"
        member_ids:
          type: array
          items:
            type: integer
          example: [123, 124, 125]
        repository_id:
          type: integer
          nullable: true
        repository_name:
          type: string
          nullable: true
        created_at:
          type: string
          format: date-time
        updated_at:
          type: string
          format: date-time

    Error:
      type: object
      required:
        - code
        - message
        - request_id
        - timestamp
      properties:
        code:
          type: string
          example: "VALIDATION_INVALID_INPUT"
        message:
          type: string
          example: "Invalid input provided"
        details:
          type: object
          additionalProperties: true
        request_id:
          type: string
          example: "req_abc123xyz"
        timestamp:
          type: string
          format: date-time
        documentation_url:
          type: string
          example: "https://docs.forgejo-classroom.org/errors/VALIDATION_INVALID_INPUT"

    PaginatedResponse:
      type: object
      required:
        - data
        - pagination
      properties:
        data:
          type: array
          items: {}
        pagination:
          type: object
          required:
            - page
            - per_page
            - total_count
            - total_pages
          properties:
            page:
              type: integer
            per_page:
              type: integer
            total_count:
              type: integer
            total_pages:
              type: integer

  responses:
    UnauthorizedError:
      description: Authentication token is missing or invalid
      content:
        application/json:
          schema:
            type: object
            properties:
              error:
                $ref: '#/components/schemas/Error'
          example:
            error:
              code: "AUTH_INVALID_TOKEN"
              message: "Invalid or missing authentication token"
              request_id: "req_abc123"
              timestamp: "2025-10-29T12:34:56Z"

    ForbiddenError:
      description: Insufficient permissions to perform this action
      content:
        application/json:
          schema:
            type: object
            properties:
              error:
                $ref: '#/components/schemas/Error'
          example:
            error:
              code: "AUTHZ_INSUFFICIENT_PERMISSIONS"
              message: "You do not have permission to access this resource"
              request_id: "req_abc123"
              timestamp: "2025-10-29T12:34:56Z"

    NotFoundError:
      description: The requested resource was not found
      content:
        application/json:
          schema:
            type: object
            properties:
              error:
                $ref: '#/components/schemas/Error'
          example:
            error:
              code: "RESOURCE_NOT_FOUND"
              message: "The requested classroom does not exist"
              request_id: "req_abc123"
              timestamp: "2025-10-29T12:34:56Z"

    ValidationError:
      description: Request validation failed
      content:
        application/json:
          schema:
            type: object
            properties:
              error:
                $ref: '#/components/schemas/Error'
          example:
            error:
              code: "VALIDATION_INVALID_INPUT"
              message: "Validation failed"
              details:
                fields:
                  - field: "title"
                    code: "VALIDATION_TOO_SHORT"
                    message: "must not be empty"
                  - field: "deadline"
                    code: "VALIDATION_INVALID_DATE"
                    message: "must be in the future"
              request_id: "req_abc123"
              timestamp: "2025-10-29T12:34:56Z"

    RateLimitError:
      description: Rate limit exceeded
      content:
        application/json:
          schema:
            type: object
            properties:
              error:
                $ref: '#/components/schemas/Error'
          example:
            error:
              code: "RATE_LIMIT_EXCEEDED"
              message: "Rate limit exceeded"
              details:
                retry_after: 1730208896
              request_id: "req_abc123"
              timestamp: "2025-10-29T12:34:56Z"

paths:
  /classrooms:
    get:
      summary: List classrooms
      description: Retrieve a paginated list of classrooms accessible to the authenticated user
      tags:
        - Classrooms
      parameters:
        - $ref: '#/components/parameters/PageParam'
        - $ref: '#/components/parameters/PerPageParam'
        - $ref: '#/components/parameters/SortParam'
        - $ref: '#/components/parameters/OrderParam'
        - name: status
          in: query
          schema:
            type: string
            enum: [active, archived]
      responses:
        '200':
          description: Successful response
          headers:
            X-Total-Count:
              schema:
                type: integer
              description: Total number of classrooms
            X-Page:
              schema:
                type: integer
            X-Per-Page:
              schema:
                type: integer
            Link:
              schema:
                type: string
              description: Pagination links
          content:
            application/json:
              schema:
                allOf:
                  - $ref: '#/components/schemas/PaginatedResponse'
                  - type: object
                    properties:
                      data:
                        type: array
                        items:
                          $ref: '#/components/schemas/Classroom'
              examples:
                success:
                  value:
                    data:
                      - id: 1
                        name: "CS101 Fall 2025"
                        slug: "cs101-fall2025"
                        organization_id: 42
                        organization_name: "cs101-fall2025"
                        owner_id: 5
                        status: "active"
                        student_count: 45
                        assignment_count: 8
                        created_at: "2025-09-01T10:00:00Z"
                        updated_at: "2025-10-29T12:00:00Z"
                    pagination:
                      page: 1
                      per_page: 30
                      total_count: 1
                      total_pages: 1
        '401':
          $ref: '#/components/responses/UnauthorizedError'
        '429':
          $ref: '#/components/responses/RateLimitError'

    post:
      summary: Create classroom
      description: Create a new classroom with an associated Forgejo organization
      tags:
        - Classrooms
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - name
                - organization_name
              properties:
                name:
                  type: string
                  minLength: 1
                  maxLength: 255
                  example: "CS101 Fall 2025"
                description:
                  type: string
                  maxLength: 1000
                  example: "Introduction to Computer Science"
                organization_name:
                  type: string
                  pattern: "^[a-z0-9][a-z0-9-]*[a-z0-9]$"
                  example: "cs101-fall2025"
                visibility:
                  type: string
                  enum: [private, public]
                  default: private
            examples:
              minimal:
                value:
                  name: "CS101 Fall 2025"
                  organization_name: "cs101-fall2025"
              full:
                value:
                  name: "CS101 Fall 2025"
                  description: "Introduction to Computer Science"
                  organization_name: "cs101-fall2025"
                  visibility: "private"
      responses:
        '201':
          description: Classroom created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Classroom'
              example:
                id: 1
                name: "CS101 Fall 2025"
                slug: "cs101-fall2025"
                description: "Introduction to Computer Science"
                organization_id: 42
                organization_name: "cs101-fall2025"
                owner_id: 5
                status: "active"
                student_count: 0
                assignment_count: 0
                created_at: "2025-10-29T12:34:56Z"
                updated_at: "2025-10-29T12:34:56Z"
        '400':
          $ref: '#/components/responses/ValidationError'
        '401':
          $ref: '#/components/responses/UnauthorizedError'
        '409':
          description: Organization name already exists
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    $ref: '#/components/schemas/Error'
              example:
                error:
                  code: "RESOURCE_ALREADY_EXISTS"
                  message: "Organization name already exists"
                  details:
                    field: "organization_name"
                  request_id: "req_abc123"
                  timestamp: "2025-10-29T12:34:56Z"
        '429':
          $ref: '#/components/responses/RateLimitError'

  /classrooms/{classroom_id}:
    parameters:
      - name: classroom_id
        in: path
        required: true
        schema:
          type: integer
        description: Classroom ID

    get:
      summary: Get classroom
      description: Retrieve detailed information about a specific classroom
      tags:
        - Classrooms
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Classroom'
        '401':
          $ref: '#/components/responses/UnauthorizedError'
        '403':
          $ref: '#/components/responses/ForbiddenError'
        '404':
          $ref: '#/components/responses/NotFoundError'

    patch:
      summary: Update classroom
      description: Update classroom properties
      tags:
        - Classrooms
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                name:
                  type: string
                  minLength: 1
                  maxLength: 255
                description:
                  type: string
                  maxLength: 1000
            example:
              name: "CS101 Fall 2025 - Updated"
              description: "Updated description"
      responses:
        '200':
          description: Classroom updated successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Classroom'
        '400':
          $ref: '#/components/responses/ValidationError'
        '401':
          $ref: '#/components/responses/UnauthorizedError'
        '403':
          $ref: '#/components/responses/ForbiddenError'
        '404':
          $ref: '#/components/responses/NotFoundError'

    delete:
      summary: Delete classroom
      description: |
        Permanently delete a classroom and its associated organization.
        **Warning**: This action cannot be undone.
      tags:
        - Classrooms
      responses:
        '204':
          description: Classroom deleted successfully
        '401':
          $ref: '#/components/responses/UnauthorizedError'
        '403':
          $ref: '#/components/responses/ForbiddenError'
        '404':
          $ref: '#/components/responses/NotFoundError'

  /classrooms/{classroom_id}/assignments:
    parameters:
      - name: classroom_id
        in: path
        required: true
        schema:
          type: integer

    get:
      summary: List assignments
      description: List all assignments for a classroom
      tags:
        - Assignments
      parameters:
        - $ref: '#/components/parameters/PageParam'
        - $ref: '#/components/parameters/PerPageParam'
        - name: type
          in: query
          schema:
            type: string
            enum: [individual, team]
        - name: status
          in: query
          schema:
            type: string
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                allOf:
                  - $ref: '#/components/schemas/PaginatedResponse'
                  - type: object
                    properties:
                      data:
                        type: array
                        items:
                          $ref: '#/components/schemas/Assignment'
        '401':
          $ref: '#/components/responses/UnauthorizedError'
        '403':
          $ref: '#/components/responses/ForbiddenError'
        '404':
          $ref: '#/components/responses/NotFoundError'

    post:
      summary: Create assignment
      description: Create a new assignment in the classroom
      tags:
        - Assignments
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - title
                - slug
                - template_repo
                - type
              properties:
                title:
                  type: string
                  minLength: 1
                  maxLength: 255
                slug:
                  type: string
                  pattern: "^[a-z0-9]([a-z0-9-]*[a-z0-9])?$"
                  minLength: 1
                  maxLength: 100
                instructions:
                  type: string
                  maxLength: 10000
                template_repo:
                  type: string
                  description: "Format: owner/repo"
                type:
                  type: string
                  enum: [individual, team]
                deadline:
                  type: string
                  format: date-time
                max_team_size:
                  type: integer
                  minimum: 2
                  maximum: 10
                allow_late_submissions:
                  type: boolean
                  default: true
                visibility:
                  type: string
                  enum: [private, public]
                  default: private
            examples:
              individual:
                summary: Individual assignment
                value:
                  title: "Homework 1: Variables"
                  slug: "hw01"
                  template_repo: "cs101-templates/hw01-starter"
                  type: "individual"
                  deadline: "2025-11-15T23:59:59Z"
                  instructions: "Complete all exercises"
              team:
                summary: Team assignment
                value:
                  title: "Final Project"
                  slug: "final-project"
                  template_repo: "cs101-templates/project-starter"
                  type: "team"
                  max_team_size: 4
                  deadline: "2025-12-15T23:59:59Z"
      responses:
        '201':
          description: Assignment created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Assignment'
        '400':
          $ref: '#/components/responses/ValidationError'
        '401':
          $ref: '#/components/responses/UnauthorizedError'
        '403':
          $ref: '#/components/responses/ForbiddenError'
        '404':
          $ref: '#/components/responses/NotFoundError'
        '422':
          description: Business logic error
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    $ref: '#/components/schemas/Error'
              example:
                error:
                  code: "BUSINESS_TEMPLATE_NOT_FOUND"
                  message: "Template repository not found"
                  details:
                    template_repo: "cs101-templates/nonexistent"
                  request_id: "req_abc123"
                  timestamp: "2025-10-29T12:34:56Z"

  /assignments/{assignment_id}:
    parameters:
      - name: assignment_id
        in: path
        required: true
        schema:
          type: integer

    get:
      summary: Get assignment
      description: Get detailed information about an assignment
      tags:
        - Assignments
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Assignment'
        '401':
          $ref: '#/components/responses/UnauthorizedError'
        '403':
          $ref: '#/components/responses/ForbiddenError'
        '404':
          $ref: '#/components/responses/NotFoundError'

    patch:
      summary: Update assignment
      description: Update assignment properties
      tags:
        - Assignments
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                title:
                  type: string
                instructions:
                  type: string
                deadline:
                  type: string
                  format: date-time
                allow_late_submissions:
                  type: boolean
      responses:
        '200':
          description: Assignment updated successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Assignment'
        '400':
          $ref: '#/components/responses/ValidationError'
        '401':
          $ref: '#/components/responses/UnauthorizedError'
        '403':
          $ref: '#/components/responses/ForbiddenError'
        '404':
          $ref: '#/components/responses/NotFoundError'

    delete:
      summary: Delete assignment
      description: Delete an assignment (only if no submissions exist)
      tags:
        - Assignments
      responses:
        '204':
          description: Assignment deleted successfully
        '401':
          $ref: '#/components/responses/UnauthorizedError'
        '403':
          $ref: '#/components/responses/ForbiddenError'
        '404':
          $ref: '#/components/responses/NotFoundError'
        '409':
          description: Cannot delete assignment with submissions
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    $ref: '#/components/schemas/Error'

  /assignments/{assignment_id}/submissions:
    parameters:
      - name: assignment_id
        in: path
        required: true
        schema:
          type: integer

    post:
      summary: Accept assignment (create submission)
      description: |
        Student accepts an assignment, creating a personal repository from the template.
        This operation is idempotent - accepting twice returns the same submission.
      tags:
        - Submissions
      responses:
        '201':
          description: Submission created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Submission'
              example:
                id: 42
                assignment_id: 10
                student_id: 123
                repository_id: 500
                repository_name: "cs101-fall2025/hw01-jdoe"
                repository_url: "https://forgejo.example.com/cs101-fall2025/hw01-jdoe"
                status: "in_progress"
                accepted_at: "2025-10-29T12:34:56Z"
                is_late: false
                commit_count: 0
                created_at: "2025-10-29T12:34:56Z"
                updated_at: "2025-10-29T12:34:56Z"
        '200':
          description: Assignment already accepted (idempotent)
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Submission'
        '401':
          $ref: '#/components/responses/UnauthorizedError'
        '403':
          $ref: '#/components/responses/ForbiddenError'
        '404':
          $ref: '#/components/responses/NotFoundError'
        '422':
          description: Business logic error
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    $ref: '#/components/schemas/Error'
              examples:
                deadline_passed:
                  value:
                    error:
                      code: "BUSINESS_DEADLINE_PASSED"
                      message: "Assignment deadline has passed"
                      request_id: "req_abc123"
                      timestamp: "2025-10-29T12:34:56Z"
                not_in_roster:
                  value:
                    error:
                      code: "BUSINESS_ROSTER_NOT_FOUND"
                      message: "Student not found in classroom roster"
                      request_id: "req_abc123"
                      timestamp: "2025-10-29T12:34:56Z"

  /submissions:
    get:
      summary: List submissions
      description: |
        List submissions with flexible filtering.
        At least one filter parameter must be provided.
      tags:
        - Submissions
      parameters:
        - $ref: '#/components/parameters/PageParam'
        - $ref: '#/components/parameters/PerPageParam'
        - name: assignment_id
          in: query
          schema:
            type: integer
          description: Filter by assignment
        - name: student_id
          in: query
          schema:
            type: integer
          description: Filter by student
        - name: classroom_id
          in: query
          schema:
            type: integer
          description: Filter by classroom
        - name: status
          in: query
          schema:
            type: string
            enum: [pending, in_progress, submitted, graded]
          description: Filter by status
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                allOf:
                  - $ref: '#/components/schemas/PaginatedResponse'
                  - type: object
                    properties:
                      data:
                        type: array
                        items:
                          $ref: '#/components/schemas/Submission'
        '400':
          description: Missing required filter parameter
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    $ref: '#/components/schemas/Error'
        '401':
          $ref: '#/components/responses/UnauthorizedError'

  /submissions/{submission_id}:
    parameters:
      - name: submission_id
        in: path
        required: true
        schema:
          type: integer

    get:
      summary: Get submission
      description: Get detailed information about a submission
      tags:
        - Submissions
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Submission'
        '401':
          $ref: '#/components/responses/UnauthorizedError'
        '403':
          $ref: '#/components/responses/ForbiddenError'
        '404':
          $ref: '#/components/responses/NotFoundError'

    patch:
      summary: Update submission
      description: Update submission (grading, feedback, etc.)
      tags:
        - Submissions
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                status:
                  type: string
                  enum: [submitted, graded]
                grade:
                  type: number
                  format: float
                feedback:
                  type: string
            example:
              status: "graded"
              grade: 95.5
              feedback: "Excellent work!"
      responses:
        '200':
          description: Submission updated successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Submission'
        '400':
          $ref: '#/components/responses/ValidationError'
        '401':
          $ref: '#/components/responses/UnauthorizedError'
        '403':
          $ref: '#/components/responses/ForbiddenError'
        '404':
          $ref: '#/components/responses/NotFoundError'

  /classrooms/{classroom_id}/roster:
    parameters:
      - name: classroom_id
        in: path
        required: true
        schema:
          type: integer

    get:
      summary: List roster entries
      description: List all students in the classroom roster
      tags:
        - Roster
      parameters:
        - $ref: '#/components/parameters/PageParam'
        - $ref: '#/components/parameters/PerPageParam'
        - name: status
          in: query
          schema:
            type: string
            enum: [pending, linked]
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                allOf:
                  - $ref: '#/components/schemas/PaginatedResponse'
                  - type: object
                    properties:
                      data:
                        type: array
                        items:
                          $ref: '#/components/schemas/RosterEntry'
        '401':
          $ref: '#/components/responses/UnauthorizedError'
        '403':
          $ref: '#/components/responses/ForbiddenError'
        '404':
          $ref: '#/components/responses/NotFoundError'

    post:
      summary: Add roster entry
      description: Add a single student to the roster
      tags:
        - Roster
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - identifier
                - email
                - full_name
              properties:
                identifier:
                  type: string
                  example: "student001"
                email:
                  type: string
                  format: email
                full_name:
                  type: string
      responses:
        '201':
          description: Roster entry created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/RosterEntry'
        '400':
          $ref: '#/components/responses/ValidationError'
        '401':
          $ref: '#/components/responses/UnauthorizedError'
        '403':
          $ref: '#/components/responses/ForbiddenError'
        '409':
          description: Duplicate identifier
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    $ref: '#/components/schemas/Error'

  /classrooms/{classroom_id}/roster/bulk:
    parameters:
      - name: classroom_id
        in: path
        required: true
        schema:
          type: integer

    post:
      summary: Bulk roster operations
      description: |
        Perform multiple roster operations in a single request.
        For small operations (<100 entries), synchronous.
      tags:
        - Roster
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - operations
              properties:
                operations:
                  type: array
                  items:
                    type: object
                    required:
                      - action
                    properties:
                      action:
                        type: string
                        enum: [add, remove]
                      identifier:
                        type: string
                      email:
                        type: string
                      full_name:
                        type: string
            example:
              operations:
                - action: "add"
                  identifier: "student001"
                  email: "s1@edu"
                  full_name: "Student One"
                - action: "add"
                  identifier: "student002"
                  email: "s2@edu"
                  full_name: "Student Two"
                - action: "remove"
                  identifier: "student003"
      responses:
        '200':
          description: Operations completed
          content:
            application/json:
              schema:
                type: object
                properties:
                  results:
                    type: array
                    items:
                      type: object
                      properties:
                        index:
                          type: integer
                        status:
                          type: string
                          enum: [success, error]
                        id:
                          type: integer
                          nullable: true
                        error:
                          type: object
                          nullable: true
                  summary:
                    type: object
                    properties:
                      total:
                        type: integer
                      succeeded:
                        type: integer
                      failed:
                        type: integer
              example:
                results:
                  - index: 0
                    status: "success"
                    id: 123
                  - index: 1
                    status: "error"
                    error:
                      code: "RESOURCE_ALREADY_EXISTS"
                      message: "Duplicate identifier"
                  - index: 2
                    status: "success"
                    id: null
                summary:
                  total: 3
                  succeeded: 2
                  failed: 1
        '400':
          $ref: '#/components/responses/ValidationError'
        '401':
          $ref: '#/components/responses/UnauthorizedError'
        '403':
          $ref: '#/components/responses/ForbiddenError'

  /classrooms/{classroom_id}/roster/import:
    parameters:
      - name: classroom_id
        in: path
        required: true
        schema:
          type: integer

    post:
      summary: Import roster from CSV
      description: |
        Import roster from CSV file. For large files (100+ entries),
        this operation is asynchronous and returns a job ID.
      tags:
        - Roster
      requestBody:
        required: true
        content:
          multipart/form-data:
            schema:
              type: object
              required:
                - file
              properties:
                file:
                  type: string
                  format: binary
      responses:
        '202':
          description: Import job created (async)
          content:
            application/json:
              schema:
                type: object
                properties:
                  job_id:
                    type: string
                  status:
                    type: string
                    enum: [processing]
                  status_url:
                    type: string
              example:
                job_id: "import_abc123"
                status: "processing"
                status_url: "/jobs/import_abc123"
        '400':
          $ref: '#/components/responses/ValidationError'
        '401':
          $ref: '#/components/responses/UnauthorizedError'
        '403':
          $ref: '#/components/responses/ForbiddenError'

  /jobs/{job_id}:
    parameters:
      - name: job_id
        in: path
        required: true
        schema:
          type: string

    get:
      summary: Get job status
      description: Check the status of an asynchronous job
      tags:
        - Jobs
      responses:
        '200':
          description: Job status retrieved
          content:
            application/json:
              schema:
                type: object
                properties:
                  id:
                    type: string
                  status:
                    type: string
                    enum: [processing, completed, failed]
                  progress:
                    type: object
                    properties:
                      total:
                        type: integer
                      processed:
                        type: integer
                      succeeded:
                        type: integer
                      failed:
                        type: integer
                  errors_url:
                    type: string
                    nullable: true
                  created_at:
                    type: string
                    format: date-time
                  completed_at:
                    type: string
                    format: date-time
                    nullable: true
              example:
                id: "import_abc123"
                status: "completed"
                progress:
                  total: 500
                  processed: 500
                  succeeded: 485
                  failed: 15
                errors_url: "/jobs/import_abc123/errors"
                created_at: "2025-10-29T10:00:00Z"
                completed_at: "2025-10-29T10:02:34Z"
        '404':
          $ref: '#/components/responses/NotFoundError'
```

---

## **15. Migration Strategy**

*(Content remains the same as original document, Section 14)*

---

## **16. Future Roadmap**

*(Content remains the same as original document, Section 15)*

---

## **17. Success Metrics**

*(Content remains the same as original document, Section 16)*

---

## **18. Risk Assessment and Mitigation**

*(Content remains the same as original document, Section 17)*

---

## **Appendices**

### **Appendix A: Glossary**

- **Classroom**: A logical grouping of students and assignments, mapped to a Forgejo organization
- **Assignment**: A task distributed to students, based on a template repository
- **Submission**: A student's acceptance of an assignment, resulting in a personal repository
- **Roster**: The list of students enrolled in a classroom
- **Template Repository**: A repository containing starter code for an assignment
- **Deadline Snapshot**: A Git tag created at assignment deadline time
- **Invitation Code**: A unique token used in assignment invitation URLs
- **Forgejo Organization**: A multi-user account in Forgejo that owns repositories

### **Appendix B: References**

1. Forgejo Documentation: https://forgejo.org/docs/latest/
2. Forgejo API Reference: https://forgejo.your.host/api/swagger
3. Go Project Layout: https://github.com/golang-standards/project-layout
4. REST API Design Best Practices: https://stackoverflow.blog/2020/03/02/best-practices-for-rest-api-design/
5. OpenAPI Specification: https://swagger.io/specification/
6. Prometheus Best Practices: https://prometheus.io/docs/practices/naming/
7. Semantic Versioning: https://semver.org/

### **Appendix C: Database Schema**

*(Detailed schema would be included here - summary version below)*

```sql
-- Core tables
CREATE TABLE classrooms (...);
CREATE TABLE assignments (...);
CREATE TABLE submissions (...);
CREATE TABLE roster_entries (...);
CREATE TABLE teams (...);
CREATE TABLE team_members (...);

-- Indexes for performance
CREATE INDEX idx_submissions_assignment ON submissions(assignment_id);
CREATE INDEX idx_submissions_student ON submissions(student_id);
CREATE INDEX idx_submissions_status ON submissions(status);
CREATE INDEX idx_roster_classroom ON roster_entries(classroom_id);
CREATE INDEX idx_roster_identifier ON roster_entries(identifier);
```

---

**End of Technical Implementation Proposal**

**Version**: 2.0
**Date**: October 29, 2025
**Status**: Design Review - Revised
**Next Review**: After Phase 1 completion

**Summary of Changes from v1.0:**
- Added comprehensive language choice rationale (Section 2)
- Improved RESTful API design with proper resource modeling (Section 5)
- Standardized error handling with comprehensive taxonomy (Section 6)
- Added API versioning strategy (Section 7)
- Enhanced infrastructure with caching and queue layers (Section 8)
- Added contract testing methodology (Section 10)
- Improved security implementations (Section 11)
- Expanded monitoring and observability (Section 13)
- Complete OpenAPI 3.0 specification (Section 14)
