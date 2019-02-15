package storage

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"

	"gitlab.com/easywork.me/backend/models"
)

func (s *Storage) ContractsGet(ctx context.Context, user *models.User) ([]models.Contract, error) {
	var out = []models.Contract{}

	filter := bson.M{}
	switch user.Role {
	case models.Work:
		filter["contractor_id"] = user.ID
	case models.Hire:
		filter["owner_id"] = user.ID
	default:
		return nil, errors.Errorf("unsupported user role: %q", user.Role)
	}

	cur, err := s.contracts().Find(ctx, filter)
	if err != nil {
		return nil, errors.Wrapf(err, "error on get contracts")
	}

	for cur.Next(ctx) {
		var c models.Contract
		err = cur.Decode(&c)
		if err != nil {
			return nil, errors.Wrapf(err, "error on decode contract")
		}
		out = append(out, c)
	}

	if err = cur.Err(); err != nil {
		return nil, errors.Wrapf(err, "error on iterate contracts")
	}

	return out, nil
}

// ContractsCreate - create contract
func (s *Storage) ContractsCreate(ctx context.Context, c *models.ContractBase, ownerID models.UserID) (*models.Contract, error) {
	if c == nil {
		c = models.NewContractBase()
	}

	c.OwnerID = ownerID
	c.CreatedDateTime = time.Now()
	res, err := s.contracts().InsertOne(ctx, c)
	if err != nil {
		return nil, errors.Wrapf(err, "error on insert new contract")
	}

	var out models.Contract
	out.ID = res.InsertedID.(primitive.ObjectID)
	out.ContractBase = *c
	return &out, nil
}

// ContractsUpdateStatus - update status of contract
func (s *Storage) ContractsUpdateStatus(ctx context.Context, cID primitive.ObjectID, status models.ContractStatus, user *models.User) error {
	c, err := s.ContractGet(ctx, cID, user)
	if err != nil {
		return errors.Wrapf(err, "cannot update status of contract (id: %v)", cID.Hex())
	}

	if c.Status == models.Ended {
		return errors.Errorf("contract is ended (id: %v)", cID.Hex())
	}

	if c.Status == status {
		switch status {
		case models.Started:
			return &models.ContractAlreadyStarted
		case models.Paused:
			return &models.ContractAlreadyPaused
		}
		return &models.ContractAlreadySameStatus
	}

	_, err = s.ProjectGetByOwner(ctx, c.ProjectID, user.ID)
	if err != nil {
		return errors.Wrapf(err, "cannot update status of contract (id: %v)", c.ID)
	}

	res, err := s.contracts().UpdateOne(ctx, bson.M{"_id": cID}, bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		return errors.Wrapf(err, "error on update status of contract (id: %v)", cID)
	}
	if res.ModifiedCount != 1 {
		return &models.ContractNotChangedStatus
	}
	return nil
}

// ContractGet - get contract
func (s *Storage) ContractGet(ctx context.Context, cID primitive.ObjectID, user *models.User) (*models.Contract, error) {
	filter := bson.M{"_id": cID}

	switch user.Role {
	case models.Work:
		filter["contractor_id"] = user.ID
	case models.Hire:
		filter["owner_id"] = user.ID
	default:
		return nil, errors.Errorf("unsupported user role: %q", user.Role)
	}

	var c models.Contract
	err := s.contracts().FindOne(ctx, filter).Decode(&c)
	if err != nil {
		return nil, errors.Wrapf(err, "error on get contract (id: %v)", cID)
	}
	return &c, nil
}

func (s *Storage) contracts() *mongo.Collection {
	return s.c("contracts")
}
