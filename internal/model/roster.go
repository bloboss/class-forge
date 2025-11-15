package model

import (
	"time"
)

// RosterEntry represents a student in a classroom roster
type RosterEntry struct {
	ID              int64      `json:"id" db:"id"`
	ClassroomID     int64      `json:"classroom_id" db:"classroom_id"`
	StudentName     string     `json:"student_name" db:"student_name"`
	StudentEmail    string     `json:"student_email" db:"student_email"`
	StudentID       string     `json:"student_id" db:"student_id"`
	ForgejoUsername *string    `json:"forgejo_username,omitempty" db:"forgejo_username"`
	ForgejoUserID   *int64     `json:"forgejo_user_id,omitempty" db:"forgejo_user_id"`
	Role            string     `json:"role" db:"role"` // student, assistant, instructor
	LinkedAt        *time.Time `json:"linked_at,omitempty" db:"linked_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// AddStudentRequest represents the request to add a student to roster
type AddStudentRequest struct {
	StudentName  string `json:"student_name" binding:"required"`
	StudentEmail string `json:"student_email" binding:"required,email"`
	StudentID    string `json:"student_id" binding:"required"`
	Role         string `json:"role"` // defaults to "student"
}

// LinkStudentRequest represents the request to link a student account
type LinkStudentRequest struct {
	ForgejoUsername string `json:"forgejo_username" binding:"required"`
}

// RosterListRequest represents the request to list roster entries
type RosterListRequest struct {
	LinkedOnly   bool `form:"linked_only" json:"linked_only,omitempty"`
	UnlinkedOnly bool `form:"unlinked_only" json:"unlinked_only,omitempty"`
	Page         int  `form:"page" json:"page,omitempty"`
	PerPage      int  `form:"per_page" json:"per_page,omitempty"`
}

// RosterListResponse represents the response for listing roster entries
type RosterListResponse struct {
	Students   []RosterEntry `json:"students"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PerPage    int           `json:"per_page"`
	TotalPages int           `json:"total_pages"`
}

// IsLinked returns true if the student has a linked Forgejo account
func (r *RosterEntry) IsLinked() bool {
	return r.ForgejoUsername != nil && *r.ForgejoUsername != ""
}
