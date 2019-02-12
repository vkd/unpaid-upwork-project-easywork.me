package storage

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

// TotalsUpdate - update total value for contract_id by date
func (s *Storage) TotalsUpdate(ctx context.Context, cID bson.ObjectId, createdAt time.Time, value int) error {
	res, err := s.totals().UpdateOne(ctx,
		bson.M{"contract_id": cID, "date": createdAt},
		bson.M{"$inc": bson.M{"value": 10}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return errors.Wrapf(err, "error on update total value (contract_id: %v, date: %v)", cID, createdAt)
	}
	if res.ModifiedCount == 0 || res.UpsertedCount == 0 {
		return errors.Errorf("total value not updated (contract_id: %v, date: %v)", cID, createdAt)
	}
	return nil
}

func (s *Storage) totals() *mongo.Collection {
	return s.db().Collection("totals")
}
