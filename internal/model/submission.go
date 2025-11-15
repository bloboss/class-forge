package model

import (
	"time"
)

// Submission represents a student's assignment submission
type Submission struct {
	ID                int64      `json:"id" db:"id"`
	AssignmentID      int64      `json:"assignment_id" db:"assignment_id"`
	StudentID         *int64     `json:"student_id,omitempty" db:"student_id"`
	TeamID            *int64     `json:"team_id,omitempty" db:"team_id"`
	RepositoryName    string     `json:"repository_name" db:"repository_name"`
	RepositoryID      int64      `json:"repository_id" db:"repository_id"`
	RepositoryURL     string     `json:"repository_url" db:"repository_url"`
	Status            string     `json:"status" db:"status"` // pending, accepted, late
	AcceptedAt        *time.Time `json:"accepted_at,omitempty" db:"accepted_at"`
	LastCommitSHA     *string    `json:"last_commit_sha,omitempty" db:"last_commit_sha"`
	LastCommitMessage *string    `json:"last_commit_message,omitempty" db:"last_commit_message"`
	CommitCount       int        `json:"commit_count" db:"commit_count"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// SubmissionListRequest represents the request to list submissions
type SubmissionListRequest struct {
	AssignmentID   *int64 `form:"assignment_id" json:"assignment_id,omitempty"`
	Status         string `form:"status" json:"status,omitempty"`
	TeamOnly       bool   `form:"team_only" json:"team_only,omitempty"`
	IndividualOnly bool   `form:"individual_only" json:"individual_only,omitempty"`
	Page           int    `form:"page" json:"page,omitempty"`
	PerPage        int    `form:"per_page" json:"per_page,omitempty"`
}

// SubmissionListResponse represents the response for listing submissions
type SubmissionListResponse struct {
	Submissions []Submission `json:"submissions"`
	Total       int          `json:"total"`
	Page        int          `json:"page"`
	PerPage     int          `json:"per_page"`
	TotalPages  int          `json:"total_pages"`
}

// IsTeamSubmission returns true if this is a team submission
func (s *Submission) IsTeamSubmission() bool {
	return s.TeamID != nil
}

// IsIndividualSubmission returns true if this is an individual submission
func (s *Submission) IsIndividualSubmission() bool {
	return s.StudentID != nil
}

// IsAccepted returns true if the submission has been accepted
func (s *Submission) IsAccepted() bool {
	return s.Status == "accepted"
}

// IsLate returns true if the submission was submitted after deadline
func (s *Submission) IsLate() bool {
	return s.Status == "late"
}
