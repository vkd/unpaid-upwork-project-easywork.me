package storage

import (
	"context"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/pkg/errors"
	"gitlab.com/easywork.me/backend/models"
	"gopkg.in/mgo.v2/bson"
)

// InvitationCreate - create invitation
func (s *Storage) InvitationCreate(ctx context.Context, pID bson.ObjectId, inviteeID models.UserID, tms *models.TermsSet, user *models.User) (*models.Invitation, error) {
	proj, err := s.ProjectGet(ctx, pID)
	if err != nil {
		return nil, errors.Wrapf(err, "invitation not created")
	}

	u, err := s.UserGet(ctx, inviteeID)
	if err != nil {
		return nil, errors.Wrapf(err, "invitation not created")
	}

	if proj.OwnerId == u.ID {
		return nil, &models.UserCannotBeInvitedToHisOwnProject
	}

	tms, err = s.TermsCreate(ctx, tms)
	if err != nil {
		return nil, errors.Wrapf(err, "error on create terms set")
	}

	inv := models.NewInvitation()
	inv.ProjectId = proj.Id
	inv.OwnerID = user.ID
	inv.InviteeID = u.ID
	inv.InviteeEmail = u.Email
	inv.TermsId = tms.ID

	res, err := s.invitations().InsertOne(ctx, inv)
	if err != nil {
		return nil, errors.Wrapf(err, "error on create invitation")
	}

	inv.ID = res.InsertedID.(bson.ObjectId)
	return inv, nil
}

// InvitationUpdateStatus - update status of invitation
func (s *Storage) InvitationUpdateStatus(ctx context.Context, iID bson.ObjectId, uID models.UserID, status models.InvitationStatus) (*models.Invitation, error) {
	i, err := s.InvitationGet(ctx, iID)
	if err != nil {
		return nil, errors.Wrapf(err, "invitation not updated (id: %v)", iID)
	}

	if i.InviteeID != uID {
		return nil, &models.AccessForbidden
	}

	res, err := s.invitations().UpdateOne(ctx, iID, bson.M{"status": status})
	if err != nil {
		return nil, errors.Wrapf(err, "error on update status on invitation (id: %v)", iID)
	}
	if res.ModifiedCount != 1 {
		return nil, errors.Errorf("invitation (id: %v) not updated", iID)
	}

	i.Status = status
	return i, nil
}

// InvitationGet - get one invitation by id
func (s *Storage) InvitationGet(ctx context.Context, iID bson.ObjectId) (*models.Invitation, error) {
	var i models.Invitation
	err := s.invitations().FindOne(ctx, iID).Decode(&i)
	if err != nil {
		return nil, errors.Wrapf(err, "error on get invitation (id: %v)", iID)
	}
	return &i, nil
}

func (s *Storage) invitations() *mongo.Collection {
	return s.db().Collection("invitations")
}
