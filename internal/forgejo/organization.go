package forgejo

import (
	"context"
	"fmt"

	"code.gitea.io/sdk/gitea"
	"go.uber.org/zap"
)

// CreateOrganizationOptions holds options for creating an organization
type CreateOrganizationOptions struct {
	Name        string
	FullName    string
	Description string
	Website     string
	Location    string
	Visibility  string // public, limited, private
}

// CreateOrganization creates a new organization in Forgejo
func (c *Client) CreateOrganization(ctx context.Context, opts CreateOrganizationOptions) (*gitea.Organization, error) {
	c.logger.Debug("creating organization",
		zap.String("name", opts.Name),
		zap.String("visibility", opts.Visibility),
	)

	// Set default visibility
	visibility := gitea.VisibleTypePublic
	switch opts.Visibility {
	case "private":
		visibility = gitea.VisibleTypePrivate
	case "limited":
		visibility = gitea.VisibleTypeLimited
	}

	org, _, err := c.client.CreateOrg(gitea.CreateOrgOption{
		Name:        opts.Name,
		FullName:    opts.FullName,
		Description: opts.Description,
		Website:     opts.Website,
		Location:    opts.Location,
		Visibility:  visibility,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	c.logger.Info("organization created",
		zap.String("name", org.UserName),
		zap.Int64("id", org.ID),
	)

	return org, nil
}

// GetOrganization retrieves an organization by name
func (c *Client) GetOrganization(ctx context.Context, name string) (*gitea.Organization, error) {
	c.logger.Debug("getting organization", zap.String("name", name))

	org, _, err := c.client.GetOrg(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	return org, nil
}

// OrganizationExists checks if an organization exists
func (c *Client) OrganizationExists(ctx context.Context, name string) (bool, error) {
	_, _, err := c.client.GetOrg(name)
	if err != nil {
		// Check if it's a 404 error
		if gitea.IsErrNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check organization existence: %w", err)
	}
	return true, nil
}

// DeleteOrganization deletes an organization
func (c *Client) DeleteOrganization(ctx context.Context, name string) error {
	c.logger.Debug("deleting organization", zap.String("name", name))

	_, err := c.client.DeleteOrg(name)
	if err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}

	c.logger.Info("organization deleted", zap.String("name", name))
	return nil
}

// ListOrganizations lists all organizations visible to the authenticated user
func (c *Client) ListOrganizations(ctx context.Context, page, limit int) ([]*gitea.Organization, error) {
	c.logger.Debug("listing organizations",
		zap.Int("page", page),
		zap.Int("limit", limit),
	)

	orgs, _, err := c.client.ListMyOrgs(gitea.ListOrgsOptions{
		ListOptions: gitea.ListOptions{
			Page:     page,
			PageSize: limit,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}

	return orgs, nil
}

// CreateTeam creates a team in an organization
func (c *Client) CreateTeam(ctx context.Context, orgName string, opts gitea.CreateTeamOption) (*gitea.Team, error) {
	c.logger.Debug("creating team",
		zap.String("org", orgName),
		zap.String("team", opts.Name),
	)

	team, _, err := c.client.CreateTeam(orgName, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	c.logger.Info("team created",
		zap.String("org", orgName),
		zap.String("team", team.Name),
		zap.Int64("team_id", team.ID),
	)

	return team, nil
}

// AddTeamMember adds a user to a team
func (c *Client) AddTeamMember(ctx context.Context, teamID int64, username string) error {
	c.logger.Debug("adding team member",
		zap.Int64("team_id", teamID),
		zap.String("username", username),
	)

	_, err := c.client.AddTeamMember(teamID, username)
	if err != nil {
		return fmt.Errorf("failed to add team member: %w", err)
	}

	c.logger.Info("team member added",
		zap.Int64("team_id", teamID),
		zap.String("username", username),
	)

	return nil
}

// RemoveTeamMember removes a user from a team
func (c *Client) RemoveTeamMember(ctx context.Context, teamID int64, username string) error {
	c.logger.Debug("removing team member",
		zap.Int64("team_id", teamID),
		zap.String("username", username),
	)

	_, err := c.client.RemoveTeamMember(teamID, username)
	if err != nil {
		return fmt.Errorf("failed to remove team member: %w", err)
	}

	c.logger.Info("team member removed",
		zap.Int64("team_id", teamID),
		zap.String("username", username),
	)

	return nil
}

// GetTeam retrieves a team by ID
func (c *Client) GetTeam(ctx context.Context, teamID int64) (*gitea.Team, error) {
	c.logger.Debug("getting team", zap.Int64("team_id", teamID))

	team, _, err := c.client.GetTeam(teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	return team, nil
}

// ListTeamMembers lists members of a team
func (c *Client) ListTeamMembers(ctx context.Context, teamID int64) ([]*gitea.User, error) {
	c.logger.Debug("listing team members", zap.Int64("team_id", teamID))

	members, _, err := c.client.ListTeamMembers(teamID, gitea.ListTeamMembersOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list team members: %w", err)
	}

	return members, nil
}
