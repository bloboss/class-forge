package forgejo

import (
	"context"

	"code.gitea.io/sdk/gitea"
)

// ForgejoClient defines the interface for Forgejo operations
// This interface is used for dependency injection and mocking in tests
type ForgejoClient interface {
	// Health and authentication
	HealthCheck(ctx context.Context) error
	GetVersion(ctx context.Context) (string, error)
	GetCurrentUser(ctx context.Context) (*gitea.User, error)
	Close() error

	// Organization operations
	CreateOrganization(ctx context.Context, opts CreateOrganizationOptions) (*gitea.Organization, error)
	GetOrganization(ctx context.Context, name string) (*gitea.Organization, error)
	OrganizationExists(ctx context.Context, name string) (bool, error)
	DeleteOrganization(ctx context.Context, name string) error
	ListOrganizations(ctx context.Context, page, limit int) ([]*gitea.Organization, error)

	// Team operations
	CreateTeam(ctx context.Context, orgName string, opts gitea.CreateTeamOption) (*gitea.Team, error)
	AddTeamMember(ctx context.Context, teamID int64, username string) error
	RemoveTeamMember(ctx context.Context, teamID int64, username string) error
	GetTeam(ctx context.Context, teamID int64) (*gitea.Team, error)
	ListTeamMembers(ctx context.Context, teamID int64) ([]*gitea.User, error)

	// Repository operations
	CreateRepository(ctx context.Context, opts CreateRepositoryOptions) (*gitea.Repository, error)
	CreateOrgRepository(ctx context.Context, orgName string, opts CreateRepositoryOptions) (*gitea.Repository, error)
	GetRepository(ctx context.Context, owner, repo string) (*gitea.Repository, error)
	RepositoryExists(ctx context.Context, owner, repo string) (bool, error)
	DeleteRepository(ctx context.Context, owner, repo string) error
	ForkRepository(ctx context.Context, owner, repo, orgName string) (*gitea.Repository, error)
	GenerateRepository(ctx context.Context, templateOwner, templateRepo string, opts CreateRepositoryOptions) (*gitea.Repository, error)
	GenerateOrgRepository(ctx context.Context, templateOwner, templateRepo, orgName string, opts CreateRepositoryOptions) (*gitea.Repository, error)
	AddCollaborator(ctx context.Context, owner, repo, collaborator string, permission gitea.AccessMode) error
	CreateBranch(ctx context.Context, owner, repo, branchName, ref string) error
	ProtectBranch(ctx context.Context, owner, repo, branch string, opts gitea.BranchProtection) error
	CreateTag(ctx context.Context, owner, repo, tagName, target, message string) error
	ListRepositoryBranches(ctx context.Context, owner, repo string) ([]*gitea.Branch, error)
	GetRepositoryFile(ctx context.Context, owner, repo, filepath, ref string) ([]byte, error)

	// User operations
	GetUser(ctx context.Context, username string) (*gitea.User, error)
	UserExists(ctx context.Context, username string) (bool, error)
	SearchUsers(ctx context.Context, query string, limit int) ([]*gitea.User, error)
	ListUsers(ctx context.Context, page, limit int) ([]*gitea.User, error)
	CreateUser(ctx context.Context, opts gitea.CreateUserOption) (*gitea.User, error)
	DeleteUser(ctx context.Context, username string) error
	GetUserEmails(ctx context.Context) ([]*gitea.Email, error)
	IsUserOrgMember(ctx context.Context, orgName, username string) (bool, error)
	AddOrgMember(ctx context.Context, orgName, username string) error
	RemoveOrgMember(ctx context.Context, orgName, username string) error
	ListOrgMembers(ctx context.Context, orgName string, page, limit int) ([]*gitea.User, error)
}

// Ensure Client implements ForgejoClient interface
var _ ForgejoClient = (*Client)(nil)
