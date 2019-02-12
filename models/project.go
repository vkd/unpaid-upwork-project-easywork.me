package models

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type Project struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`

	ProjectBase `bson:"inline"`
}

type ProjectBase struct {
	OwnerID         UserID    `json:"owner_id" bson:"owner_id"`
	Title           string    `json:"title" bson:"title"`
	Description     string    `json:"description" bson:"description"`
	CreatedDateTime time.Time `json:"created_date_time" bson:"created_date_time"`
}

func NewProject(ownerID UserID) *Project {
	return &Project{
		ProjectBase: ProjectBase{
			OwnerID:         ownerID,
			CreatedDateTime: time.Now(),
		},
	}
}
