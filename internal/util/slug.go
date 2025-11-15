package util

import (
	"regexp"
	"strings"
	"unicode"
)

// GenerateSlug generates a URL-friendly slug from a string
func GenerateSlug(input string) string {
	// Convert to lowercase
	slug := strings.ToLower(input)

	// Replace non-alphanumeric characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")

	// Limit length to 50 characters
	if len(slug) > 50 {
		slug = slug[:50]
		slug = strings.TrimRight(slug, "-")
	}

	return slug
}

// IsValidSlug checks if a string is a valid slug
func IsValidSlug(slug string) bool {
	if slug == "" {
		return false
	}

	// Check format: lowercase alphanumeric with hyphens
	slugRegex := regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	return slugRegex.MatchString(slug)
}

// SanitizeIdentifier removes non-alphanumeric characters and ensures valid identifier format
func SanitizeIdentifier(input string) string {
	// Remove non-alphanumeric characters except hyphens and underscores
	reg := regexp.MustCompile(`[^a-zA-Z0-9_-]+`)
	sanitized := reg.ReplaceAllString(input, "")

	// Ensure it doesn't start with a number
	if len(sanitized) > 0 && unicode.IsDigit(rune(sanitized[0])) {
		sanitized = "repo-" + sanitized
	}

	return sanitized
}

// GenerateRepositoryName generates a repository name from classroom and assignment info
func GenerateRepositoryName(classroomSlug, assignmentSlug, studentIdentifier string) string {
	parts := []string{classroomSlug, assignmentSlug, SanitizeIdentifier(studentIdentifier)}
	return strings.Join(parts, "-")
}

// GenerateTeamRepositoryName generates a repository name for team assignments
func GenerateTeamRepositoryName(classroomSlug, assignmentSlug, teamSlug string) string {
	parts := []string{classroomSlug, assignmentSlug, "team", teamSlug}
	return strings.Join(parts, "-")
}
