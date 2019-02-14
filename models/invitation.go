package models

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type Invitation struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`

	InvitationBase `bson:"inline"`
}

type InvitationBase struct {
	Status InvitationStatus `json:"status" bson:"status"`

	OwnerID   UserID `json:"owner_id" bson:"owner_id"`
	InviteeID UserID `json:"invitee_id" bson:"invitee_id"`

	ProjectID    primitive.ObjectID `json:"project_id" bson:"project_id"`
	InviteeEmail string             `json:"invitee_email"`
	TermsID      primitive.ObjectID `json:"terms_id" bson:"terms_id"`

	CreatedDateTime time.Time `json:"created_date_time" bson:"created_date_time"`

	TermsBase `bson:"inline"`
}

type InvitationStatus string

const InvitationStatusPending InvitationStatus = "pending"
const InvitationStatusAccepted InvitationStatus = "accepted"
const InvitationStatusDeclined InvitationStatus = "declined"

func (i *InvitationBase) PreCreate() *InvitationBase {
	i.Status = InvitationStatusPending
	i.TermsBase = *i.TermsBase.PreCreate()
	return i
}
