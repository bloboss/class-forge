package forgejo

import (
	"context"
	"fmt"

	"code.gitea.io/sdk/gitea"
	"go.uber.org/zap"
)

// GetUser retrieves a user by username
func (c *Client) GetUser(ctx context.Context, username string) (*gitea.User, error) {
	c.logger.Debug("getting user", zap.String("username", username))

	user, _, err := c.client.GetUserInfo(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// UserExists checks if a user exists
func (c *Client) UserExists(ctx context.Context, username string) (bool, error) {
	_, _, err := c.client.GetUserInfo(username)
	if err != nil {
		if gitea.IsErrNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return true, nil
}

// SearchUsers searches for users by query
func (c *Client) SearchUsers(ctx context.Context, query string, limit int) ([]*gitea.User, error) {
	c.logger.Debug("searching users",
		zap.String("query", query),
		zap.Int("limit", limit),
	)

	users, _, err := c.client.SearchUsers(gitea.SearchUsersOption{
		KeyWord: query,
		ListOptions: gitea.ListOptions{
			Page:     1,
			PageSize: limit,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	return users, nil
}

// ListUsers lists all users (admin only)
func (c *Client) ListUsers(ctx context.Context, page, limit int) ([]*gitea.User, error) {
	c.logger.Debug("listing users",
		zap.Int("page", page),
		zap.Int("limit", limit),
	)

	users, _, err := c.client.AdminListUsers(gitea.AdminListUsersOptions{
		ListOptions: gitea.ListOptions{
			Page:     page,
			PageSize: limit,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// CreateUser creates a new user (admin only)
func (c *Client) CreateUser(ctx context.Context, opts gitea.CreateUserOption) (*gitea.User, error) {
	c.logger.Debug("creating user",
		zap.String("username", opts.Username),
		zap.String("email", opts.Email),
	)

	user, _, err := c.client.AdminCreateUser(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	c.logger.Info("user created",
		zap.String("username", user.UserName),
		zap.Int64("id", user.ID),
	)

	return user, nil
}

// DeleteUser deletes a user (admin only)
func (c *Client) DeleteUser(ctx context.Context, username string) error {
	c.logger.Debug("deleting user", zap.String("username", username))

	_, err := c.client.AdminDeleteUser(username)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	c.logger.Info("user deleted", zap.String("username", username))
	return nil
}

// GetUserEmails retrieves emails for a user
func (c *Client) GetUserEmails(ctx context.Context) ([]*gitea.Email, error) {
	c.logger.Debug("getting user emails")

	emails, _, err := c.client.ListEmails(gitea.ListEmailsOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get user emails: %w", err)
	}

	return emails, nil
}

// IsUserOrgMember checks if a user is a member of an organization
func (c *Client) IsUserOrgMember(ctx context.Context, orgName, username string) (bool, error) {
	c.logger.Debug("checking organization membership",
		zap.String("org", orgName),
		zap.String("username", username),
	)

	isMember, _, err := c.client.CheckOrgMembership(orgName, username)
	if err != nil {
		return false, fmt.Errorf("failed to check organization membership: %w", err)
	}

	return isMember, nil
}

// AddOrgMember adds a user to an organization
func (c *Client) AddOrgMember(ctx context.Context, orgName, username string) error {
	c.logger.Debug("adding organization member",
		zap.String("org", orgName),
		zap.String("username", username),
	)

	_, err := c.client.AddOrgMembership(orgName, username)
	if err != nil {
		return fmt.Errorf("failed to add organization member: %w", err)
	}

	c.logger.Info("organization member added",
		zap.String("org", orgName),
		zap.String("username", username),
	)

	return nil
}

// RemoveOrgMember removes a user from an organization
func (c *Client) RemoveOrgMember(ctx context.Context, orgName, username string) error {
	c.logger.Debug("removing organization member",
		zap.String("org", orgName),
		zap.String("username", username),
	)

	_, err := c.client.RemoveOrgMembership(orgName, username)
	if err != nil {
		return fmt.Errorf("failed to remove organization member: %w", err)
	}

	c.logger.Info("organization member removed",
		zap.String("org", orgName),
		zap.String("username", username),
	)

	return nil
}

// ListOrgMembers lists members of an organization
func (c *Client) ListOrgMembers(ctx context.Context, orgName string, page, limit int) ([]*gitea.User, error) {
	c.logger.Debug("listing organization members",
		zap.String("org", orgName),
		zap.Int("page", page),
		zap.Int("limit", limit),
	)

	members, _, err := c.client.ListOrgMembers(orgName, gitea.ListOrgMembersOptions{
		ListOptions: gitea.ListOptions{
			Page:     page,
			PageSize: limit,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list organization members: %w", err)
	}

	return members, nil
}
