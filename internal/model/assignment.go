package model

import (
	"time"
)

// Assignment represents an assignment entity
type Assignment struct {
	ID                   int64      `json:"id" db:"id"`
	ClassroomID          int64      `json:"classroom_id" db:"classroom_id"`
	Name                 string     `json:"name" db:"name"`
	Slug                 string     `json:"slug" db:"slug"`
	Description          string     `json:"description" db:"description"`
	TemplateRepository   string     `json:"template_repository" db:"template_repository"`
	TemplateRepositoryID int64      `json:"template_repository_id" db:"template_repository_id"`
	Deadline             *time.Time `json:"deadline,omitempty" db:"deadline"`
	MaxTeamSize          int        `json:"max_team_size" db:"max_team_size"`
	AutoAccept           bool       `json:"auto_accept" db:"auto_accept"`
	Public               bool       `json:"public" db:"public"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
}

// CreateAssignmentRequest represents the request to create an assignment
type CreateAssignmentRequest struct {
	ClassroomID        int64  `json:"classroom_id" binding:"required"`
	Name               string `json:"name" binding:"required"`
	Description        string `json:"description"`
	TemplateRepository string `json:"template_repository" binding:"required"`
	Deadline           string `json:"deadline,omitempty"` // RFC3339 format
	MaxTeamSize        int    `json:"max_team_size"`
	AutoAccept         bool   `json:"auto_accept"`
	Public             bool   `json:"public"`
}

// UpdateAssignmentRequest represents the request to update an assignment
type UpdateAssignmentRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Deadline    *string `json:"deadline,omitempty"` // RFC3339 format
	MaxTeamSize *int    `json:"max_team_size,omitempty"`
	AutoAccept  *bool   `json:"auto_accept,omitempty"`
	Public      *bool   `json:"public,omitempty"`
}

// AssignmentListRequest represents the request to list assignments
type AssignmentListRequest struct {
	ClassroomID int64  `form:"classroom_id" json:"classroom_id,omitempty"`
	Status      string `form:"status" json:"status,omitempty"` // active, past, all
	Page        int    `form:"page" json:"page,omitempty"`
	PerPage     int    `form:"per_page" json:"per_page,omitempty"`
}

// AssignmentListResponse represents the response for listing assignments
type AssignmentListResponse struct {
	Assignments []Assignment `json:"assignments"`
	Total       int          `json:"total"`
	Page        int          `json:"page"`
	PerPage     int          `json:"per_page"`
	TotalPages  int          `json:"total_pages"`
}

// AssignmentStats represents statistics for an assignment
type AssignmentStats struct {
	AssignmentID      int64   `json:"assignment_id"`
	TotalStudents     int     `json:"total_students"`
	AcceptedCount     int     `json:"accepted_count"`
	SubmissionCount   int     `json:"submission_count"`
	TeamCount         int     `json:"team_count"`
	AcceptanceRate    float64 `json:"acceptance_rate"`
	SubmissionRate    float64 `json:"submission_rate"`
	AverageCommits    float64 `json:"average_commits"`
	OnTimeSubmissions int     `json:"on_time_submissions"`
	LateSubmissions   int     `json:"late_submissions"`
}

// AcceptAssignmentRequest represents the request to accept an assignment
type AcceptAssignmentRequest struct {
	TeamName string `json:"team_name,omitempty"` // For team assignments
}

// IsTeamAssignment returns true if this assignment allows teams
func (a *Assignment) IsTeamAssignment() bool {
	return a.MaxTeamSize > 1
}

// IsIndividualAssignment returns true if this assignment is individual only
func (a *Assignment) IsIndividualAssignment() bool {
	return a.MaxTeamSize == 1
}

// IsActive returns true if the assignment deadline hasn't passed
func (a *Assignment) IsActive() bool {
	if a.Deadline == nil {
		return true // No deadline means always active
	}
	return time.Now().Before(*a.Deadline)
}

// IsPast returns true if the assignment deadline has passed
func (a *Assignment) IsPast() bool {
	if a.Deadline == nil {
		return false // No deadline means never past
	}
	return time.Now().After(*a.Deadline)
}

// Validate validates the create assignment request
func (req *CreateAssignmentRequest) Validate() error {
	// TODO: Implement validation logic
	return nil
}

// Validate validates the update assignment request
func (req *UpdateAssignmentRequest) Validate() error {
	// TODO: Implement validation logic
	return nil
}
