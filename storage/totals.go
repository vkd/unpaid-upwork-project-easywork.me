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

// TotalDaily - get total by daily
func (s *Storage) TotalDaily(ctx context.Context, cID primitive.ObjectID, from, to string) ([]models.Total, error) {
	filter := bson.M{"contract_id": cID}
	if from != "" || to != "" {
		byDate := bson.M{}
		if from != "" {
			byDate["$gte"] = from
		}
		if to != "" {
			byDate["$lte"] = to
		}
		filter["date"] = byDate
	}

	var out = []models.Total{}

	cur, err := s.totals().Find(ctx, filter)
	if err != nil {
		return nil, errors.Wrapf(err, "error on get daily total (contract_id: %v)", cID.Hex())
	}
	for cur.Next(ctx) {
		var t models.Total
		err = cur.Decode(&t)
		if err != nil {
			return nil, errors.Wrapf(err, "error on decode totals")
		}
		out = append(out, t)
	}
	if err = cur.Err(); err != nil {
		return nil, errors.Wrapf(err, "error on cursor of totals result")
	}
	return out, nil
}

// func (s *Storage) TotalGet(ctx context.Context, cID primitive.ObjectID) (int, error) {
// 	filter := bson.M{"contract_id": cID}
// }

// TotalsUpdate - update total value for contract_id by date
func (s *Storage) TotalsUpdate(ctx context.Context, cID primitive.ObjectID, createdAt time.Time, value int) error {
	date := createdAt.Format("2006-01-02")

	res, err := s.totals().UpdateOne(ctx,
		bson.M{"contract_id": cID, "date": date},
		bson.M{"$inc": bson.M{"value": value}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return errors.Wrapf(err, "error on update total value (contract_id: %v, date: %v)", cID, date)
	}
	if res.ModifiedCount == 0 && res.UpsertedCount == 0 {
		return errors.Errorf("total value not updated (contract_id: %v, date: %v)", cID, date)
	}
	return nil
}

func (s *Storage) totals() *mongo.Collection {
	return s.db().Collection("totals")
}
