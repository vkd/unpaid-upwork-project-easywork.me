package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Project struct {
	Id              bson.ObjectId `json:"id" bson:"id"`
	CreatedDateTime time.Time     `json:"created_date_time" bson:"created_date_time"`
	OwnerId         UserID        `json:"owner_id" bson:"owner_id"`
	Title           string        `json:"title" bson:"title"`
	Description     string        `json:"description" bson:"description"`
}
