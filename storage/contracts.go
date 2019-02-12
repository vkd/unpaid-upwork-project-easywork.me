package storage

import (
	"context"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"

	"gitlab.com/easywork.me/backend/models"
)

const (
	colContracts = "contracts"
)

// ContractsCreate - create contract
func (s *Storage) ContractsCreate(ctx context.Context, c *models.Contract) (*models.Contract, error) {
	if c == nil {
		c = models.NewContract()
	}

	res, err := s.contracts().InsertOne(ctx, c)
	if err != nil {
		return nil, errors.Wrapf(err, "error on insert new contract")
	}

	c.ID = res.InsertedID.(bson.ObjectId)

	return c, nil
}

// ContractsUpdateStatus - update status of contract
func (s *Storage) ContractsUpdateStatus(ctx context.Context, cID bson.ObjectId, role models.Role, status models.ContractStatus) error {
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

	_, err = s.ProjectGet(ctx, c.ProjectID)
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
