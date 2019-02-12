package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Invitation struct {
	ID     bson.ObjectId    `json:"id" bson:"_id"`
	Status InvitationStatus `json:"status" bson:"status"`

	OwnerID   UserID `json:"owner_id" bson:"owner_id"`
	InviteeID UserID `json:"invitee_id" bson:"invitee_id"`

	ProjectId    bson.ObjectId `json:"project_id" bson:"project_id"`
	InviteeEmail string        `json:"invitee_email"`
	TermsId      bson.ObjectId `json:"terms_id" bson:"terms_id"`

	CreatedDateTime time.Time `json:"created_date_time" bson:"created_date_time"`
}

type InvitationStatus string

const InvitationStatusPending InvitationStatus = "pending"
const InvitationStatusAccepted InvitationStatus = "accepted"
const InvitationStatusDeclined InvitationStatus = "declined"

func NewInvitation() *Invitation {
	return &Invitation{
		Status:          InvitationStatusPending,
		CreatedDateTime: time.Now(),
	}
}
