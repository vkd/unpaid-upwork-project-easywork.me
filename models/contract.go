package models

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type Contract struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`

	ContractBase `bson:"inline"`
}

type ContractBase struct {
	Status ContractStatus `json:"status" bson:"status"`

	ProjectID primitive.ObjectID `json:"project_id" bson:"project_id"`
	// ProjectTitle       string             `json:"project_title" bson:"project_title"`
	// ProjectDescription string             `json:"project_description" bson:"project_description"`
	ContractorID    UserID             `json:"contractor_id" bson:"contractor_id"`
	OwnerID         UserID             `json:"owner_id" bson:"owner_id"`
	TermsID         primitive.ObjectID `json:"terms_id" bson:"terms_id"`
	CreatedDateTime time.Time          `json:"created_date_time" bson:"created_date_time"`
}

type ContractStatus string

const (
	NotStarted ContractStatus = "not_started"
	Started    ContractStatus = "started"
	Paused     ContractStatus = "paused"
	Ended      ContractStatus = "ended"
)

func NewContractBase() *ContractBase {
	var c ContractBase
	c.Status = NotStarted
	return &c
}

func (c *ContractBase) FromInvitation(i *Invitation) *ContractBase {
	c.ProjectID = i.ProjectID
	c.ContractorID = i.InviteeID
	c.TermsID = i.TermsID
	return c
}
