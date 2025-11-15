package forgejo

import (
	"context"

	"code.gitea.io/sdk/gitea"
	"github.com/stretchr/testify/mock"
)

// MockClient is a mock implementation of ForgejoClient for testing
type MockClient struct {
	mock.Mock
}

// Ensure MockClient implements ForgejoClient interface
var _ ForgejoClient = (*MockClient)(nil)

// Health and authentication

func (m *MockClient) HealthCheck(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockClient) GetVersion(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockClient) GetCurrentUser(ctx context.Context) (*gitea.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gitea.User), args.Error(1)
}

func (m *MockClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Organization operations

func (m *MockClient) CreateOrganization(ctx context.Context, opts CreateOrganizationOptions) (*gitea.Organization, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gitea.Organization), args.Error(1)
}

func (m *MockClient) GetOrganization(ctx context.Context, name string) (*gitea.Organization, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gitea.Organization), args.Error(1)
}

func (m *MockClient) OrganizationExists(ctx context.Context, name string) (bool, error) {
	args := m.Called(ctx, name)
	return args.Bool(0), args.Error(1)
}

func (m *MockClient) DeleteOrganization(ctx context.Context, name string) error {
	args := m.Called(ctx, name)
	return args.Error(0)
}

func (m *MockClient) ListOrganizations(ctx context.Context, page, limit int) ([]*gitea.Organization, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*gitea.Organization), args.Error(1)
}

// Team operations

func (m *MockClient) CreateTeam(ctx context.Context, orgName string, opts gitea.CreateTeamOption) (*gitea.Team, error) {
	args := m.Called(ctx, orgName, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gitea.Team), args.Error(1)
}

func (m *MockClient) AddTeamMember(ctx context.Context, teamID int64, username string) error {
	args := m.Called(ctx, teamID, username)
	return args.Error(0)
}

func (m *MockClient) RemoveTeamMember(ctx context.Context, teamID int64, username string) error {
	args := m.Called(ctx, teamID, username)
	return args.Error(0)
}

func (m *MockClient) GetTeam(ctx context.Context, teamID int64) (*gitea.Team, error) {
	args := m.Called(ctx, teamID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gitea.Team), args.Error(1)
}

func (m *MockClient) ListTeamMembers(ctx context.Context, teamID int64) ([]*gitea.User, error) {
	args := m.Called(ctx, teamID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*gitea.User), args.Error(1)
}

// Repository operations

func (m *MockClient) CreateRepository(ctx context.Context, opts CreateRepositoryOptions) (*gitea.Repository, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gitea.Repository), args.Error(1)
}

func (m *MockClient) CreateOrgRepository(ctx context.Context, orgName string, opts CreateRepositoryOptions) (*gitea.Repository, error) {
	args := m.Called(ctx, orgName, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gitea.Repository), args.Error(1)
}

func (m *MockClient) GetRepository(ctx context.Context, owner, repo string) (*gitea.Repository, error) {
	args := m.Called(ctx, owner, repo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gitea.Repository), args.Error(1)
}

func (m *MockClient) RepositoryExists(ctx context.Context, owner, repo string) (bool, error) {
	args := m.Called(ctx, owner, repo)
	return args.Bool(0), args.Error(1)
}

func (m *MockClient) DeleteRepository(ctx context.Context, owner, repo string) error {
	args := m.Called(ctx, owner, repo)
	return args.Error(0)
}

func (m *MockClient) ForkRepository(ctx context.Context, owner, repo, orgName string) (*gitea.Repository, error) {
	args := m.Called(ctx, owner, repo, orgName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gitea.Repository), args.Error(1)
}

func (m *MockClient) GenerateRepository(ctx context.Context, templateOwner, templateRepo string, opts CreateRepositoryOptions) (*gitea.Repository, error) {
	args := m.Called(ctx, templateOwner, templateRepo, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gitea.Repository), args.Error(1)
}

func (m *MockClient) GenerateOrgRepository(ctx context.Context, templateOwner, templateRepo, orgName string, opts CreateRepositoryOptions) (*gitea.Repository, error) {
	args := m.Called(ctx, templateOwner, templateRepo, orgName, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gitea.Repository), args.Error(1)
}

func (m *MockClient) AddCollaborator(ctx context.Context, owner, repo, collaborator string, permission gitea.AccessMode) error {
	args := m.Called(ctx, owner, repo, collaborator, permission)
	return args.Error(0)
}

func (m *MockClient) CreateBranch(ctx context.Context, owner, repo, branchName, ref string) error {
	args := m.Called(ctx, owner, repo, branchName, ref)
	return args.Error(0)
}

func (m *MockClient) ProtectBranch(ctx context.Context, owner, repo, branch string, opts gitea.BranchProtection) error {
	args := m.Called(ctx, owner, repo, branch, opts)
	return args.Error(0)
}

func (m *MockClient) CreateTag(ctx context.Context, owner, repo, tagName, target, message string) error {
	args := m.Called(ctx, owner, repo, tagName, target, message)
	return args.Error(0)
}

func (m *MockClient) ListRepositoryBranches(ctx context.Context, owner, repo string) ([]*gitea.Branch, error) {
	args := m.Called(ctx, owner, repo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*gitea.Branch), args.Error(1)
}

func (m *MockClient) GetRepositoryFile(ctx context.Context, owner, repo, filepath, ref string) ([]byte, error) {
	args := m.Called(ctx, owner, repo, filepath, ref)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

// User operations

func (m *MockClient) GetUser(ctx context.Context, username string) (*gitea.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gitea.User), args.Error(1)
}

func (m *MockClient) UserExists(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockClient) SearchUsers(ctx context.Context, query string, limit int) ([]*gitea.User, error) {
	args := m.Called(ctx, query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*gitea.User), args.Error(1)
}

func (m *MockClient) ListUsers(ctx context.Context, page, limit int) ([]*gitea.User, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*gitea.User), args.Error(1)
}

func (m *MockClient) CreateUser(ctx context.Context, opts gitea.CreateUserOption) (*gitea.User, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gitea.User), args.Error(1)
}

func (m *MockClient) DeleteUser(ctx context.Context, username string) error {
	args := m.Called(ctx, username)
	return args.Error(0)
}

func (m *MockClient) GetUserEmails(ctx context.Context) ([]*gitea.Email, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*gitea.Email), args.Error(1)
}

func (m *MockClient) IsUserOrgMember(ctx context.Context, orgName, username string) (bool, error) {
	args := m.Called(ctx, orgName, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockClient) AddOrgMember(ctx context.Context, orgName, username string) error {
	args := m.Called(ctx, orgName, username)
	return args.Error(0)
}

func (m *MockClient) RemoveOrgMember(ctx context.Context, orgName, username string) error {
	args := m.Called(ctx, orgName, username)
	return args.Error(0)
}

func (m *MockClient) ListOrgMembers(ctx context.Context, orgName string, page, limit int) ([]*gitea.User, error) {
	args := m.Called(ctx, orgName, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*gitea.User), args.Error(1)
}
