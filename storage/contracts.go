package storage

import (
	"context"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"

	"gitlab.com/easywork.me/backend/models"
)

const (
	colContracts = "contracts"
)

// ContractsCreate - create contract
func (s *Storage) ContractsCreate(ctx context.Context, c *models.ContractBase, ownerID models.UserID) (*models.Contract, error) {
	if c == nil {
		c = models.NewContractBase()
	}

	c.OwnerID = ownerID
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
func (s *Storage) ContractsUpdateStatus(ctx context.Context, cID bson.ObjectId, role models.Role, status models.ContractStatus, userID models.UserID) error {
	c, err := s.ContractGet(ctx, cID, role)
	if err != nil {
		return errors.Wrapf(err, "cannot update status of contract (id: %v)", cID)
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

	_, err = s.ProjectGetByOwner(ctx, c.ProjectID, userID)
	if err != nil {
		return errors.Wrapf(err, "cannot update status of contract (id: %v)", c.ID)
	}

	res, err := s.contracts().UpdateOne(ctx, bson.M{"_id": cID}, bson.M{"status": status})
	if err != nil {
		return errors.Wrapf(err, "error on update status of contract (id: %v)", cID)
	}
	if res.ModifiedCount != 1 {
		return &models.ContractNotChangedStatus
	}
	return nil
}

func (s *Storage) ContractGet(ctx context.Context, cID bson.ObjectId, role models.Role) (*models.Contract, error) {
	panic("Not implemented - auth depends, with join")
	var c models.Contract
	err := s.contracts().FindOne(ctx, bson.M{"_id": cID}).Decode(&c)
	if err != nil {
		return nil, errors.Wrapf(err, "error on get contract (id: %v)", cID)
	}
	return &c, nil
}

func (s *Storage) contracts() *mongo.Collection {
	return s.c(colContracts)
}
