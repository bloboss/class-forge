# Integration Tests

This directory contains integration tests for Forgejo Classroom that require external dependencies.

## Forgejo Client Integration Tests

The Forgejo client integration tests verify the actual API interactions with a Forgejo instance.

### Prerequisites

- A running Forgejo instance (can be local or remote)
- An API token with appropriate permissions
- Network connectivity to the Forgejo instance

### Running Integration Tests

#### Using Environment Variables

```bash
export FORGEJO_BASE_URL="https://your-forgejo-instance.com"
export FORGEJO_TOKEN="your-api-token"
go test ./test/integration/... -v
```

#### Using Docker Compose

We provide a Docker Compose setup for running a local Forgejo instance for testing:

```bash
# Start Forgejo in Docker
docker-compose -f test/integration/docker-compose.yml up -d

# Wait for Forgejo to be ready
sleep 10

# Run integration tests
export FORGEJO_BASE_URL="http://localhost:3000"
export FORGEJO_TOKEN="your-generated-token"
go test ./test/integration/... -v

# Cleanup
docker-compose -f test/integration/docker-compose.yml down -v
```

#### Skipping Integration Tests

Integration tests are automatically skipped if:
- The `FORGEJO_BASE_URL` or `FORGEJO_TOKEN` environment variables are not set
- The `-short` flag is used with `go test`

```bash
# Skip integration tests
go test ./test/integration/... -short
```

### Test Coverage

The integration tests cover:

1. **Health Check & Authentication**
   - Server connectivity
   - API token validation
   - Version retrieval

2. **Organization Operations**
   - Creating organizations
   - Retrieving organization details
   - Checking organization existence
   - Deleting organizations

3. **Repository Operations**
   - Creating user repositories
   - Creating organization repositories
   - Repository retrieval and existence checks
   - Repository deletion

4. **User Operations**
   - Getting user information
   - Checking user existence
   - User search

5. **Team Operations**
   - Creating teams
   - Adding/removing team members
   - Team retrieval

### Setting Up a Test Forgejo Instance

#### Option 1: Using Docker

```bash
docker run -d \
  --name forgejo-test \
  -p 3000:3000 \
  -e FORGEJO__security__INSTALL_LOCK=true \
  -e FORGEJO__database__DB_TYPE=sqlite3 \
  codeberg.org/forgejo/forgejo:latest
```

After starting:
1. Navigate to http://localhost:3000
2. Complete the initial setup
3. Create an API token: Settings → Applications → Generate New Token
4. Export the token as `FORGEJO_TOKEN`

#### Option 2: Using Existing Instance

If you have an existing Forgejo instance:
1. Generate an API token from your user settings
2. Export the base URL and token:
   ```bash
   export FORGEJO_BASE_URL="https://your-instance.com"
   export FORGEJO_TOKEN="your-token"
   ```

### Continuous Integration

In CI environments, integration tests should:
1. Spin up a Forgejo instance in a container
2. Wait for it to be ready
3. Generate or use a pre-configured API token
4. Run the integration tests
5. Clean up the container

Example GitHub Actions workflow:

```yaml
integration-tests:
  runs-on: ubuntu-latest
  services:
    forgejo:
      image: codeberg.org/forgejo/forgejo:latest
      ports:
        - 3000:3000
      env:
        FORGEJO__security__INSTALL_LOCK: true
        FORGEJO__database__DB_TYPE: sqlite3
  steps:
    - uses: actions/checkout@v3
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.23'
    - name: Run integration tests
      env:
        FORGEJO_BASE_URL: http://localhost:3000
        FORGEJO_TOKEN: ${{ secrets.FORGEJO_TEST_TOKEN }}
      run: go test ./test/integration/... -v
```

### Best Practices

1. **Cleanup**: Always clean up resources (organizations, repositories) created during tests
2. **Isolation**: Use unique names (e.g., timestamp-based) to avoid conflicts
3. **Timeouts**: Set appropriate timeouts for API calls
4. **Error Handling**: Check both success and error cases
5. **Idempotency**: Tests should be idempotent and not depend on each other

### Troubleshooting

**Connection Refused**
- Ensure Forgejo is running and accessible at the specified URL
- Check firewall/network settings

**Authentication Failed**
- Verify the API token is valid
- Ensure the token has appropriate permissions

**Tests Hanging**
- Check if Forgejo is responsive
- Verify network connectivity
- Review timeout settings

**Cleanup Issues**
- Manually clean up test resources via Forgejo UI if needed
- Check logs for deletion errors
