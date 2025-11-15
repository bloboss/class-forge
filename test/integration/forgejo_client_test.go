package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"code.gitea.io/sdk/gitea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"code.forgejo.org/forgejo/classroom/internal/forgejo"
)

// ForgejoClientTestSuite is an integration test suite for the Forgejo client
// These tests require a running Forgejo instance
type ForgejoClientTestSuite struct {
	suite.Suite
	client  *forgejo.Client
	baseURL string
	token   string
	ctx     context.Context
}

// SetupSuite runs once before all tests in the suite
func (s *ForgejoClientTestSuite) SetupSuite() {
	// Get Forgejo connection details from environment
	s.baseURL = os.Getenv("FORGEJO_BASE_URL")
	s.token = os.Getenv("FORGEJO_TOKEN")

	if s.baseURL == "" || s.token == "" {
		s.T().Skip("Skipping integration tests: FORGEJO_BASE_URL and FORGEJO_TOKEN must be set")
	}

	s.ctx = context.Background()

	// Create client
	logger, err := zap.NewDevelopment()
	require.NoError(s.T(), err)

	client, err := forgejo.NewClient(forgejo.ClientConfig{
		BaseURL:   s.baseURL,
		Token:     s.token,
		Timeout:   30 * time.Second,
		Logger:    logger,
		UserAgent: "forgejo-classroom-integration-test/1.0",
	})
	require.NoError(s.T(), err)
	s.client = client
}

// TearDownSuite runs once after all tests in the suite
func (s *ForgejoClientTestSuite) TearDownSuite() {
	if s.client != nil {
		s.client.Close()
	}
}

// TestHealthCheck verifies Forgejo connectivity
func (s *ForgejoClientTestSuite) TestHealthCheck() {
	err := s.client.HealthCheck(s.ctx)
	s.NoError(err, "Health check should succeed")
}

// TestGetVersion retrieves the Forgejo server version
func (s *ForgejoClientTestSuite) TestGetVersion() {
	version, err := s.client.GetVersion(s.ctx)
	s.NoError(err)
	s.NotEmpty(version, "Version should not be empty")
	s.T().Logf("Forgejo version: %s", version)
}

// TestGetCurrentUser retrieves the authenticated user
func (s *ForgejoClientTestSuite) TestGetCurrentUser() {
	user, err := s.client.GetCurrentUser(s.ctx)
	s.NoError(err)
	s.NotNil(user)
	s.NotEmpty(user.UserName, "Username should not be empty")
	s.T().Logf("Authenticated as user: %s (ID: %d)", user.UserName, user.ID)
}

// TestOrganizationLifecycle tests creating, retrieving, and deleting an organization
func (s *ForgejoClientTestSuite) TestOrganizationLifecycle() {
	orgName := "test-org-" + s.generateTimestamp()

	// Create organization
	org, err := s.client.CreateOrganization(s.ctx, forgejo.CreateOrganizationOptions{
		Name:        orgName,
		FullName:    "Test Organization",
		Description: "Integration test organization",
		Visibility:  "public",
	})
	s.NoError(err)
	s.NotNil(org)
	s.Equal(orgName, org.UserName)

	// Get organization
	retrievedOrg, err := s.client.GetOrganization(s.ctx, orgName)
	s.NoError(err)
	s.NotNil(retrievedOrg)
	s.Equal(orgName, retrievedOrg.UserName)

	// Check organization exists
	exists, err := s.client.OrganizationExists(s.ctx, orgName)
	s.NoError(err)
	s.True(exists)

	// Cleanup: Delete organization
	err = s.client.DeleteOrganization(s.ctx, orgName)
	s.NoError(err)

	// Verify deletion
	exists, err = s.client.OrganizationExists(s.ctx, orgName)
	s.NoError(err)
	s.False(exists)
}

// TestRepositoryLifecycle tests creating, retrieving, and deleting a repository
func (s *ForgejoClientTestSuite) TestRepositoryLifecycle() {
	repoName := "test-repo-" + s.generateTimestamp()

	// Get current user to use as owner
	user, err := s.client.GetCurrentUser(s.ctx)
	s.NoError(err)

	// Create repository
	repo, err := s.client.CreateRepository(s.ctx, forgejo.CreateRepositoryOptions{
		Name:        repoName,
		Description: "Integration test repository",
		Private:     false,
		AutoInit:    true,
		Readme:      "Default",
	})
	s.NoError(err)
	s.NotNil(repo)
	s.Equal(repoName, repo.Name)

	// Get repository
	retrievedRepo, err := s.client.GetRepository(s.ctx, user.UserName, repoName)
	s.NoError(err)
	s.NotNil(retrievedRepo)
	s.Equal(repoName, retrievedRepo.Name)

	// Check repository exists
	exists, err := s.client.RepositoryExists(s.ctx, user.UserName, repoName)
	s.NoError(err)
	s.True(exists)

	// Cleanup: Delete repository
	err = s.client.DeleteRepository(s.ctx, user.UserName, repoName)
	s.NoError(err)

	// Verify deletion
	exists, err = s.client.RepositoryExists(s.ctx, user.UserName, repoName)
	s.NoError(err)
	s.False(exists)
}

// TestOrgRepositoryCreation tests creating a repository in an organization
func (s *ForgejoClientTestSuite) TestOrgRepositoryCreation() {
	orgName := "test-org-" + s.generateTimestamp()
	repoName := "test-repo-" + s.generateTimestamp()

	// Create organization
	org, err := s.client.CreateOrganization(s.ctx, forgejo.CreateOrganizationOptions{
		Name:     orgName,
		FullName: "Test Organization for Repo",
	})
	s.NoError(err)
	s.NotNil(org)

	defer func() {
		// Cleanup: Delete organization (will also delete repositories)
		s.client.DeleteOrganization(s.ctx, orgName)
	}()

	// Create repository in organization
	repo, err := s.client.CreateOrgRepository(s.ctx, orgName, forgejo.CreateRepositoryOptions{
		Name:        repoName,
		Description: "Test repository in organization",
		Private:     true,
		AutoInit:    true,
	})
	s.NoError(err)
	s.NotNil(repo)
	s.Equal(repoName, repo.Name)
	s.Equal(orgName, repo.Owner.UserName)
}

// TestUserOperations tests user-related operations
func (s *ForgejoClientTestSuite) TestUserOperations() {
	// Get current user
	user, err := s.client.GetCurrentUser(s.ctx)
	s.NoError(err)
	s.NotNil(user)

	// Get user by username
	retrievedUser, err := s.client.GetUser(s.ctx, user.UserName)
	s.NoError(err)
	s.NotNil(retrievedUser)
	s.Equal(user.UserName, retrievedUser.UserName)

	// Check user exists
	exists, err := s.client.UserExists(s.ctx, user.UserName)
	s.NoError(err)
	s.True(exists)

	// Check non-existent user
	exists, err = s.client.UserExists(s.ctx, "nonexistent-user-"+s.generateTimestamp())
	s.NoError(err)
	s.False(exists)
}

// TestTeamOperations tests team-related operations
func (s *ForgejoClientTestSuite) TestTeamOperations() {
	orgName := "test-org-" + s.generateTimestamp()

	// Create organization
	org, err := s.client.CreateOrganization(s.ctx, forgejo.CreateOrganizationOptions{
		Name:     orgName,
		FullName: "Test Organization for Teams",
	})
	s.NoError(err)

	defer func() {
		s.client.DeleteOrganization(s.ctx, orgName)
	}()

	// Create team
	team, err := s.client.CreateTeam(s.ctx, orgName, gitea.CreateTeamOption{
		Name:        "test-team",
		Description: "Test team",
		Permission:  gitea.AccessModeWrite,
	})
	s.NoError(err)
	s.NotNil(team)
	s.Equal("test-team", team.Name)

	// Get team
	retrievedTeam, err := s.client.GetTeam(s.ctx, team.ID)
	s.NoError(err)
	s.NotNil(retrievedTeam)
	s.Equal(team.ID, retrievedTeam.ID)
}

// Helper function to generate timestamp-based unique identifier
func (s *ForgejoClientTestSuite) generateTimestamp() string {
	return time.Now().Format("20060102-150405")
}

// TestForgejoClientIntegration runs the integration test suite
func TestForgejoClientIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	suite.Run(t, new(ForgejoClientTestSuite))
}

// Standalone test for basic client creation (doesn't require Forgejo)
func TestClientCreation(t *testing.T) {
	logger := zap.NewNop()

	t.Run("create client with valid config", func(t *testing.T) {
		client, err := forgejo.NewClient(forgejo.ClientConfig{
			BaseURL: "https://forgejo.example.com",
			Token:   "test-token",
			Logger:  logger,
		})
		require.NoError(t, err)
		assert.NotNil(t, client)
		client.Close()
	})

	t.Run("create client with missing URL", func(t *testing.T) {
		client, err := forgejo.NewClient(forgejo.ClientConfig{
			Token:  "test-token",
			Logger: logger,
		})
		require.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "base URL is required")
	})

	t.Run("create client with missing token", func(t *testing.T) {
		client, err := forgejo.NewClient(forgejo.ClientConfig{
			BaseURL: "https://forgejo.example.com",
			Logger:  logger,
		})
		require.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "API token is required")
	})
}
