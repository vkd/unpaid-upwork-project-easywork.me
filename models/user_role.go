package models

type UserRole string

const (
	Hire UserRole = "hire"
	Work UserRole = "work"
)

// IsWork - check role is work
func IsWork(role UserRole) bool {
	return role == Work
}
