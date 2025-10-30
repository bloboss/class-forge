package model

import (
	"time"
)

// Classroom represents a classroom entity
type Classroom struct {
	ID               int64     `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	Slug             string    `json:"slug" db:"slug"`
	Description      string    `json:"description" db:"description"`
	OrganizationName string    `json:"organization_name" db:"organization_name"`
	OrganizationID   int64     `json:"organization_id" db:"organization_id"`
	InstructorID     int64     `json:"instructor_id" db:"instructor_id"`
	InstructorLogin  string    `json:"instructor_login" db:"instructor_login"`
	Public           bool      `json:"public" db:"public"`
	Archived         bool      `json:"archived" db:"archived"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
	ArchivedAt       *time.Time `json:"archived_at,omitempty" db:"archived_at"`
}

// CreateClassroomRequest represents the request to create a classroom
type CreateClassroomRequest struct {
	Name             string `json:"name" binding:"required"`
	Description      string `json:"description"`
	OrganizationName string `json:"organization_name" binding:"required"`
	Public           bool   `json:"public"`
}

// UpdateClassroomRequest represents the request to update a classroom
type UpdateClassroomRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Public      *bool   `json:"public,omitempty"`
}

// ClassroomListRequest represents the request to list classrooms
type ClassroomListRequest struct {
	OrganizationName string `form:"organization" json:"organization,omitempty"`
	IncludeArchived  bool   `form:"archived" json:"archived,omitempty"`
	Page             int    `form:"page" json:"page,omitempty"`
	PerPage          int    `form:"per_page" json:"per_page,omitempty"`
}

// ClassroomListResponse represents the response for listing classrooms
type ClassroomListResponse struct {
	Classrooms []Classroom `json:"classrooms"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	TotalPages int         `json:"total_pages"`
}

// ClassroomStats represents statistics for a classroom
type ClassroomStats struct {
	ClassroomID     int64 `json:"classroom_id"`
	TotalStudents   int   `json:"total_students"`
	LinkedStudents  int   `json:"linked_students"`
	TotalAssignments int  `json:"total_assignments"`
	ActiveAssignments int `json:"active_assignments"`
	TotalSubmissions int  `json:"total_submissions"`
}

// Validate validates the create classroom request
func (req *CreateClassroomRequest) Validate() error {
	// TODO: Implement validation logic
	// This would be called by the validation utility
	return nil
}

// Validate validates the update classroom request
func (req *UpdateClassroomRequest) Validate() error {
	// TODO: Implement validation logic
	return nil
}