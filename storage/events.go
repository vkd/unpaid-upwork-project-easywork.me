package storage

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/pkg/errors"
	"gitlab.com/easywork.me/backend/models"
	"gopkg.in/mgo.v2/bson"
)

func (s *Storage) EventsGet(ctx context.Context, cID primitive.ObjectID, et *models.EventType, from, to *time.Time) ([]models.Event, error) {
	var out = []models.Event{}

	filter := bson.M{}
	if et != nil {
		filter["event_type"] = *et
	}
	if from != nil || to != nil {
		byDate := bson.M{}
		if from != nil {
			byDate["$gte"] = from
		}
		if to != nil {
			byDate["$lt"] = to
		}
		filter["created_date_time"] = byDate
	}

	cur, err := s.events(cID).Find(ctx, filter)
	if err != nil {
		return nil, errors.Wrapf(err, "error on get events")
	}

	for cur.Next(ctx) {
		var e models.Event
		err = cur.Decode(&e)
		if err != nil {
			return nil, errors.Wrapf(err, "error on decode event")
		}
		out = append(out, e)
	}

	if err = cur.Err(); err != nil {
		return nil, errors.Wrapf(err, "error on cursor of events")
	}

	return out, nil
}

func (s *Storage) EventsGetCountLogs(ctx context.Context, cID primitive.ObjectID, from, to *time.Time) (int64, error) {
	filter := bson.M{"event_type": models.EventLog}
	if from != nil || to != nil {
		byDate := bson.M{}
		if from != nil {
			byDate["$gte"] = from
		}
		if to != nil {
			byDate["$lt"] = to
		}
		filter["created_date_time"] = byDate
	}

	out, err := s.events(cID).Count(ctx, filter)
	if err != nil {
		return 0, errors.Wrapf(err, "error on get count of events")
	}
	return out, nil
}

func (s *Storage) EventsGetCountLogsLast24H(ctx context.Context, cID primitive.ObjectID) (int64, error) {
	to := time.Now()
	from := to.AddDate(0, 0, -1)
	return s.EventsGetCountLogs(ctx, cID, &from, &to)
}

func (s *Storage) EventsGetCountLogsCurrentWeek(ctx context.Context, cID primitive.ObjectID) (int64, error) {
	now := time.Now()

	currentWeekDay := int(now.Weekday())

	from := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1-currentWeekDay)
	to := from.AddDate(0, 0, 7) // last day not included
	return s.EventsGetCountLogs(ctx, cID, &from, &to)
}

func (s *Storage) EventsGetCountLogsPrevWeek(ctx context.Context, cID primitive.ObjectID) (int64, error) {
	now := time.Now()

	currentWeekDay := int(now.Weekday())

	from := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1-currentWeekDay-7)
	to := from.AddDate(0, 0, 7) // last day not included
	return s.EventsGetCountLogs(ctx, cID, &from, &to)
}

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
