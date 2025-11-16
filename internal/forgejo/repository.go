package forgejo

import (
	"context"
	"encoding/base64"
	"fmt"

	"code.gitea.io/sdk/gitea"
	"go.uber.org/zap"
)

// CreateRepositoryOptions holds options for creating a repository
type CreateRepositoryOptions struct {
	Name          string
	Description   string
	Private       bool
	AutoInit      bool
	DefaultBranch string
	Gitignores    string
	License       string
	Readme        string
	Template      bool
	TrustModel    string
}

// CreateRepository creates a new repository for the authenticated user
func (c *Client) CreateRepository(ctx context.Context, opts CreateRepositoryOptions) (*gitea.Repository, error) {
	c.logger.Debug("creating repository",
		zap.String("name", opts.Name),
		zap.Bool("private", opts.Private),
	)

	repo, _, err := c.client.CreateRepo(gitea.CreateRepoOption{
		Name:          opts.Name,
		Description:   opts.Description,
		Private:       opts.Private,
		AutoInit:      opts.AutoInit,
		DefaultBranch: opts.DefaultBranch,
		Gitignores:    opts.Gitignores,
		License:       opts.License,
		Readme:        opts.Readme,
		Template:      opts.Template,
		TrustModel:    gitea.TrustModel(opts.TrustModel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	c.logger.Info("repository created",
		zap.String("name", repo.Name),
		zap.Int64("id", repo.ID),
		zap.String("full_name", repo.FullName),
	)

	return repo, nil
}

// CreateOrgRepository creates a new repository in an organization
func (c *Client) CreateOrgRepository(ctx context.Context, orgName string, opts CreateRepositoryOptions) (*gitea.Repository, error) {
	c.logger.Debug("creating organization repository",
		zap.String("org", orgName),
		zap.String("name", opts.Name),
		zap.Bool("private", opts.Private),
	)

	repo, _, err := c.client.CreateOrgRepo(orgName, gitea.CreateRepoOption{
		Name:          opts.Name,
		Description:   opts.Description,
		Private:       opts.Private,
		AutoInit:      opts.AutoInit,
		DefaultBranch: opts.DefaultBranch,
		Gitignores:    opts.Gitignores,
		License:       opts.License,
		Readme:        opts.Readme,
		Template:      opts.Template,
		TrustModel:    gitea.TrustModel(opts.TrustModel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create organization repository: %w", err)
	}

	c.logger.Info("organization repository created",
		zap.String("org", orgName),
		zap.String("name", repo.Name),
		zap.Int64("id", repo.ID),
		zap.String("full_name", repo.FullName),
	)

	return repo, nil
}

// GetRepository retrieves a repository by owner and repo name
func (c *Client) GetRepository(ctx context.Context, owner, repo string) (*gitea.Repository, error) {
	c.logger.Debug("getting repository",
		zap.String("owner", owner),
		zap.String("repo", repo),
	)

	repository, _, err := c.client.GetRepo(owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}

	return repository, nil
}

// RepositoryExists checks if a repository exists
func (c *Client) RepositoryExists(ctx context.Context, owner, repo string) (bool, error) {
	_, resp, err := c.client.GetRepo(owner, repo)
	if err != nil {
		// Check if it's a 404 error
		if resp != nil && resp.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("failed to check repository existence: %w", err)
	}
	return true, nil
}

// DeleteRepository deletes a repository
func (c *Client) DeleteRepository(ctx context.Context, owner, repo string) error {
	c.logger.Debug("deleting repository",
		zap.String("owner", owner),
		zap.String("repo", repo),
	)

	_, err := c.client.DeleteRepo(owner, repo)
	if err != nil {
		return fmt.Errorf("failed to delete repository: %w", err)
	}

	c.logger.Info("repository deleted",
		zap.String("owner", owner),
		zap.String("repo", repo),
	)

	return nil
}

// ForkRepository creates a fork of a repository
func (c *Client) ForkRepository(ctx context.Context, owner, repo, orgName string) (*gitea.Repository, error) {
	c.logger.Debug("forking repository",
		zap.String("owner", owner),
		zap.String("repo", repo),
		zap.String("org", orgName),
	)

	forkedRepo, _, err := c.client.CreateFork(owner, repo, gitea.CreateForkOption{
		Organization: &orgName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fork repository: %w", err)
	}

	c.logger.Info("repository forked",
		zap.String("source", fmt.Sprintf("%s/%s", owner, repo)),
		zap.String("destination", forkedRepo.FullName),
		zap.Int64("id", forkedRepo.ID),
	)

	return forkedRepo, nil
}

// GenerateRepository creates a new repository from a template
func (c *Client) GenerateRepository(ctx context.Context, templateOwner, templateRepo string, opts CreateRepositoryOptions) (*gitea.Repository, error) {
	c.logger.Debug("generating repository from template",
		zap.String("template_owner", templateOwner),
		zap.String("template_repo", templateRepo),
		zap.String("new_name", opts.Name),
	)

	repo, _, err := c.client.CreateRepoFromTemplate(templateOwner, templateRepo, gitea.CreateRepoFromTemplateOption{
		Name:        opts.Name,
		Description: opts.Description,
		Private:     opts.Private,
		GitContent:  true,
		Topics:      true,
		GitHooks:    false,
		Webhooks:    false,
		Avatar:      true,
		Labels:      true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate repository from template: %w", err)
	}

	c.logger.Info("repository generated from template",
		zap.String("template", fmt.Sprintf("%s/%s", templateOwner, templateRepo)),
		zap.String("new_repo", repo.FullName),
		zap.Int64("id", repo.ID),
	)

	return repo, nil
}

// GenerateOrgRepository creates a new repository from a template in an organization
func (c *Client) GenerateOrgRepository(ctx context.Context, templateOwner, templateRepo, orgName string, opts CreateRepositoryOptions) (*gitea.Repository, error) {
	c.logger.Debug("generating organization repository from template",
		zap.String("template_owner", templateOwner),
		zap.String("template_repo", templateRepo),
		zap.String("org", orgName),
		zap.String("new_name", opts.Name),
	)

	repo, _, err := c.client.CreateRepoFromTemplate(templateOwner, templateRepo, gitea.CreateRepoFromTemplateOption{
		Owner:       orgName,
		Name:        opts.Name,
		Description: opts.Description,
		Private:     opts.Private,
		GitContent:  true,
		Topics:      true,
		GitHooks:    false,
		Webhooks:    false,
		Avatar:      true,
		Labels:      true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate organization repository from template: %w", err)
	}

	c.logger.Info("organization repository generated from template",
		zap.String("template", fmt.Sprintf("%s/%s", templateOwner, templateRepo)),
		zap.String("new_repo", repo.FullName),
		zap.Int64("id", repo.ID),
	)

	return repo, nil
}

// AddCollaborator adds a collaborator to a repository
func (c *Client) AddCollaborator(ctx context.Context, owner, repo, collaborator string, permission gitea.AccessMode) error {
	c.logger.Debug("adding collaborator",
		zap.String("owner", owner),
		zap.String("repo", repo),
		zap.String("collaborator", collaborator),
		zap.String("permission", string(permission)),
	)

	_, err := c.client.AddCollaborator(owner, repo, collaborator, gitea.AddCollaboratorOption{
		Permission: &permission,
	})
	if err != nil {
		return fmt.Errorf("failed to add collaborator: %w", err)
	}

	c.logger.Info("collaborator added",
		zap.String("repo", fmt.Sprintf("%s/%s", owner, repo)),
		zap.String("collaborator", collaborator),
		zap.String("permission", string(permission)),
	)

	return nil
}

// CreateBranch creates a new branch in a repository
func (c *Client) CreateBranch(ctx context.Context, owner, repo, branchName, ref string) error {
	c.logger.Debug("creating branch",
		zap.String("owner", owner),
		zap.String("repo", repo),
		zap.String("branch", branchName),
		zap.String("ref", ref),
	)

	_, _, err := c.client.CreateBranch(owner, repo, gitea.CreateBranchOption{
		BranchName:    branchName,
		OldBranchName: ref,
	})
	if err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	c.logger.Info("branch created",
		zap.String("repo", fmt.Sprintf("%s/%s", owner, repo)),
		zap.String("branch", branchName),
	)

	return nil
}

// ProtectBranch adds branch protection to a repository branch
func (c *Client) ProtectBranch(ctx context.Context, owner, repo, branch string, opts gitea.BranchProtection) error {
	c.logger.Debug("protecting branch",
		zap.String("owner", owner),
		zap.String("repo", repo),
		zap.String("branch", branch),
	)

	_, _, err := c.client.CreateBranchProtection(owner, repo, gitea.CreateBranchProtectionOption{
		BranchName:              branch,
		EnablePush:              opts.EnablePush,
		EnablePushWhitelist:     opts.EnablePushWhitelist,
		PushWhitelistUsernames:  opts.PushWhitelistUsernames,
		PushWhitelistTeams:      opts.PushWhitelistTeams,
		EnableMergeWhitelist:    opts.EnableMergeWhitelist,
		MergeWhitelistUsernames: opts.MergeWhitelistUsernames,
		MergeWhitelistTeams:     opts.MergeWhitelistTeams,
	})
	if err != nil {
		return fmt.Errorf("failed to protect branch: %w", err)
	}

	c.logger.Info("branch protected",
		zap.String("repo", fmt.Sprintf("%s/%s", owner, repo)),
		zap.String("branch", branch),
	)

	return nil
}

// CreateTag creates a new tag in a repository
func (c *Client) CreateTag(ctx context.Context, owner, repo, tagName, target, message string) error {
	c.logger.Debug("creating tag",
		zap.String("owner", owner),
		zap.String("repo", repo),
		zap.String("tag", tagName),
		zap.String("target", target),
	)

	_, _, err := c.client.CreateTag(owner, repo, gitea.CreateTagOption{
		TagName: tagName,
		Target:  target,
		Message: message,
	})
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	c.logger.Info("tag created",
		zap.String("repo", fmt.Sprintf("%s/%s", owner, repo)),
		zap.String("tag", tagName),
	)

	return nil
}

// ListRepositoryBranches lists all branches in a repository
func (c *Client) ListRepositoryBranches(ctx context.Context, owner, repo string) ([]*gitea.Branch, error) {
	c.logger.Debug("listing repository branches",
		zap.String("owner", owner),
		zap.String("repo", repo),
	)

	branches, _, err := c.client.ListRepoBranches(owner, repo, gitea.ListRepoBranchesOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list repository branches: %w", err)
	}

	return branches, nil
}

// GetRepositoryFile retrieves file content from a repository
func (c *Client) GetRepositoryFile(ctx context.Context, owner, repo, filepath, ref string) ([]byte, error) {
	c.logger.Debug("getting repository file",
		zap.String("owner", owner),
		zap.String("repo", repo),
		zap.String("filepath", filepath),
		zap.String("ref", ref),
	)

	contents, _, err := c.client.GetContents(owner, repo, ref, filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository file: %w", err)
	}

	if contents == nil || contents.Content == nil {
		return nil, fmt.Errorf("file not found: %s", filepath)
	}

	// Decode base64 content
	content, err := base64.StdEncoding.DecodeString(*contents.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file content: %w", err)
	}

	return content, nil
}
