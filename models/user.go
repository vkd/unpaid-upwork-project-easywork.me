package models

import "github.com/pkg/errors"

type UserID string

type User struct {
	ID    UserID `json:"id,omitempty" bson:"id"`
	Email string `json:"email,omitempty" bson:"email"`
	Role  Role   `json:"role" bson:"role"`

	UserProfile `bson:"inline"`
}

type UserPassword struct {
	User     `bson:"inline"`
	Password string `json:"password" bson:"password"`
}

type UserProfile struct {
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
}

func NewUser() *User {
	var u User
	u.Role = Work
	return &u
}

type Role string

const (
	Hire Role = "hire"
	Work Role = "work"
)

func CheckRole(r Role) error {
	switch r {
	case Hire, Work:
		return nil
	}
	return errors.Errorf("wrong user role (%q)", r)
}

type PublicUser struct {
	ID    UserID `json:"id,omitempty" bson:"id"`
	Email string `json:"email,omitempty" bson:"email"`
	Role  Role   `json:"role" bson:"role"`
}
