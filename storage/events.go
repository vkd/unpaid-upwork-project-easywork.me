package storage

import (
	"context"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/pkg/errors"
	"gitlab.com/easywork.me/backend/models"
	"gopkg.in/mgo.v2/bson"
)

func (s *Storage) EventCreate(ctx context.Context, e *models.EventBase, user *models.User) (*models.Event, error) {
	contract, err := s.ContractGet(ctx, e.ContractID, user)
	if err != nil {
		return nil, errors.Wrapf(err, "error on get contract", err)
	}

	if contract.Status != models.Started {
		return nil, &models.ContractIsNotStarted
	}

	col := s.events(e.ContractID)

	var latestEvent *models.Event = &models.Event{}
	err = col.FindOne(ctx, bson.M{}, &options.FindOneOptions{Sort: bson.M{"_id": -1}}).Decode(&latestEvent)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			latestEvent = nil
		} else {
			return nil, errors.Wrapf(err, "error on find last event (contract: %v)", e.ContractID)
		}
	}

	if err = checkAllowedToCreateNewEvent(latestEvent, e.EventType); err != nil {
		return nil, err
	}

	res, err := col.InsertOne(ctx, e)
	if err != nil {
		return nil, errors.Wrapf(err, "error on insert new event (contract: %v)", e.ContractID)
	}

	var out models.Event
	out.EventBase = *e
	out.ID = res.InsertedID.(primitive.ObjectID)
	return &out, nil
}

func (s *Storage) events(contractID primitive.ObjectID) *mongo.Collection {
	return s.db().Collection("events_" + contractID.Hex())
}

func checkAllowedToCreateNewEvent(oldEvent *models.Event, newType models.EventType) error {
	if oldEvent == nil {
		if newType == models.EventStart {
			return nil
		}
		return errors.Errorf(`%q event cannot be saved when last event was nil`, newType)
	}

	switch oldEvent.EventType {
	case models.EventStart:
		switch newType {
		case models.EventStop, models.EventLog:
			return nil
		}
	case models.EventStop:
		switch newType {
		case models.EventStart:
			return nil
		}
	case models.EventLog:
		switch newType {
		case models.EventLog, models.EventStop:
			return nil
		}
	}
	return errors.Errorf(`%q event cannot be saved when last event was %q`, newType, oldEvent.EventType)
}

// func (s *Storage) GetLastProject() {
// 	ctx := context.Background()
// 	var p *models.Project
// 	err := s.projects().FindOne(ctx, bson.M{}, &options.FindOneOptions{Sort: bson.M{"_id": -1}}).Decode(p)
// 	if err != nil {
// 		log.Fatalf("Error on get last: %v", err)
// 	}
// 	log.Printf("project: %#v", p)
// }
