package models

type UserID string

type User struct {
	ID    UserID   `json:"id,omitempty" bson:"id"`
	Email string   `json:"email,omitempty" bson:"email"`
	Role  UserRole `json:"role,omitempty" bson:"role"`
}
