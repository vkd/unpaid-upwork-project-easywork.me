package models

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type Contract struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	Status ContractStatus     `json:"status" bson:"status"`

	ProjectID          primitive.ObjectID `json:"project_id" bson:"project_id"`
	ProjectTitle       string             `json:"project_title" bson:"project_title"`
	ProjectDescription string             `json:"project_description" bson:"project_description"`
	ContractorID       UserID             `json:"contractor_id" bson:"contractor_id"`
	OwnerID            string             `json:"owner_id" bson:"owner_id"`
	TermsID            primitive.ObjectID `json:"terms_id" bson:"terms_id"`
	CreatedDateTime    time.Time          `json:"created_date_time" bson:"created_date_time"`
}

type ContractStatus string

const (
	NotStarted ContractStatus = "not_started"
	Started    ContractStatus = "started"
	Paused     ContractStatus = "paused"
	Ended      ContractStatus = "ended"
)

func NewContract() *Contract {
	var c Contract
	c.Status = NotStarted
	c.CreatedDateTime = time.Now()
	return &c
}

func (c *Contract) FromInvitation(i *Invitation) *Contract {
	c.ProjectID = i.ProjectID
	c.ContractorID = i.InviteeID
	c.TermsID = i.TermsID
	return c
}

func (c *Contract) SetStatus(status ContractStatus) *Contract {
	c.Status = status
	return c
}
