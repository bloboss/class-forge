package model

import (
	"time"
)

// Team represents a team for team-based assignments
type Team struct {
	ID           int64     `json:"id" db:"id"`
	AssignmentID int64     `json:"assignment_id" db:"assignment_id"`
	Name         string    `json:"name" db:"name"`
	Slug         string    `json:"slug" db:"slug"`
	Description  string    `json:"description" db:"description"`
	LeaderID     int64     `json:"leader_id" db:"leader_id"`
	MemberCount  int       `json:"member_count" db:"member_count"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// TeamMember represents a member of a team
type TeamMember struct {
	ID         int64     `json:"id" db:"id"`
	TeamID     int64     `json:"team_id" db:"team_id"`
	StudentID  int64     `json:"student_id" db:"student_id"`
	Role       string    `json:"role" db:"role"` // leader, member
	JoinedAt   time.Time `json:"joined_at" db:"joined_at"`
}

// CreateTeamRequest represents the request to create a team
type CreateTeamRequest struct {
	AssignmentID int64    `json:"assignment_id" binding:"required"`
	Name         string   `json:"name" binding:"required"`
	Description  string   `json:"description"`
	Members      []string `json:"members,omitempty"` // Initial member usernames
}

// JoinTeamRequest represents the request to join a team
type JoinTeamRequest struct {
	// No additional fields needed - user comes from auth context
}

// TeamListRequest represents the request to list teams
type TeamListRequest struct {
	AssignmentID int64 `form:"assignment_id" json:"assignment_id,omitempty"`
	ShowMembers  bool  `form:"show_members" json:"show_members,omitempty"`
	Page         int   `form:"page" json:"page,omitempty"`
	PerPage      int   `form:"per_page" json:"per_page,omitempty"`
}

// TeamListResponse represents the response for listing teams
type TeamListResponse struct {
	Teams      []TeamWithMembers `json:"teams"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	PerPage    int               `json:"per_page"`
	TotalPages int               `json:"total_pages"`
}

// TeamWithMembers represents a team with its members
type TeamWithMembers struct {
	Team
	Members []TeamMemberInfo `json:"members,omitempty"`
}

// TeamMemberInfo represents detailed information about a team member
type TeamMemberInfo struct {
	TeamMember
	StudentName     string `json:"student_name"`
	ForgejoUsername string `json:"forgejo_username"`
}

// CanAddMember returns true if the team can accept more members
func (t *Team) CanAddMember(maxTeamSize int) bool {
	return t.MemberCount < maxTeamSize
}

// IsFull returns true if the team is at capacity
func (t *Team) IsFull(maxTeamSize int) bool {
	return t.MemberCount >= maxTeamSize
}