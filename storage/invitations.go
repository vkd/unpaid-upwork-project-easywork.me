package storage

import (
	"context"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/pkg/errors"
	"gitlab.com/easywork.me/backend/models"
	"gopkg.in/mgo.v2/bson"
)

// InvitationCreate - create invitation
func (s *Storage) InvitationCreate(ctx context.Context, inv *models.InvitationBase, tmsb *models.TermsSetBase, userID models.UserID) (*models.Invitation, error) {
	proj, err := s.ProjectGetByOwner(ctx, inv.ProjectID, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "invitation not created")
	}

	u, err := s.UserGetByIDOrEmail(ctx, inv.InviteeID, inv.InviteeEmail)
	if err != nil {
		return nil, errors.Wrapf(err, "invitation not created")
	}
	if u.Role != models.Work {
		return nil, errors.Errorf("cannot allow to invite non-work type of user (id: %v, email: %v)", u.ID, u.Email)
	}

	inv.InviteeID = u.ID
	inv.InviteeEmail = u.Email

	if proj.OwnerID == inv.InviteeID {
		return nil, &models.UserCannotBeInvitedToHisOwnProject
	}

	tms, err := s.TermsCreate(ctx, tmsb)
	if err != nil {
		return nil, errors.Wrapf(err, "error on create terms set")
	}

	inv = inv.PreCreate()
	inv.OwnerID = userID
	inv.TermsID = tms.ID

	res, err := s.invitations().InsertOne(ctx, inv)
	if err != nil {
		return nil, errors.Wrapf(err, "error on create invitation")
	}

	var out models.Invitation
	out.InvitationBase = *inv
	out.ID = res.InsertedID.(primitive.ObjectID)
	return &out, nil
}

// InvitationUpdateStatus - update status of invitation
func (s *Storage) InvitationUpdateStatus(ctx context.Context, iID primitive.ObjectID, uID models.UserID, status models.InvitationStatus) (*models.Invitation, error) {
	i, err := s.invitationGet(ctx, iID)
	if err != nil {
		return nil, errors.Wrapf(err, "invitation not updated (id: %v)", iID)
	}

	if i.InviteeID != uID {
		return nil, &models.AccessForbidden
	}

	res, err := s.invitations().UpdateOne(ctx, bson.M{"_id": iID}, bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		return nil, errors.Wrapf(err, "error on update status on invitation (id: %v)", iID)
	}
	if res.ModifiedCount == 0 {
		return nil, ErrNoUpdated
	}

	i.Status = status
	return i, nil
}

// InvitationGet - get one invitation by id
func (s *Storage) InvitationGet(ctx context.Context, iID primitive.ObjectID, ownerID models.UserID) (*models.Invitation, error) {
	return s.invitationGet(ctx, iID, ownerID)
}

func (s *Storage) invitationGet(ctx context.Context, iID primitive.ObjectID, ownerID ...models.UserID) (*models.Invitation, error) {
	filter := bson.M{"_id": iID}
	if len(ownerID) > 0 {
		filter["owner_id"] = ownerID[0]
	}

	var i models.Invitation
	err := s.invitations().FindOne(ctx, filter).Decode(&i)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &models.InvitationNotFound
		}
		return nil, errors.Wrapf(err, "error on get invitation (id: %v)", iID)
	}
	return &i, nil
}

func (s *Storage) invitations() *mongo.Collection {
	return s.db().Collection("invitations")
}
