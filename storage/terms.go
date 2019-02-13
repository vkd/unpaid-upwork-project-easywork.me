package storage

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/pkg/errors"
	"gitlab.com/easywork.me/backend/models"
	"gopkg.in/mgo.v2/bson"
)

// TermsCreate - create term
func (s *Storage) TermsCreate(ctx context.Context, ts *models.TermsSetBase) (*models.TermsSet, error) {
	ts = ts.PreCreate()
	ts.CreatedDateTime = time.Now()
	res, err := s.terms().InsertOne(ctx, ts)
	if err != nil {
		return nil, errors.Wrapf(err, "error on create terms")
	}

	var out models.TermsSet
	out.ID = res.InsertedID.(primitive.ObjectID)
	out.TermsSetBase = *ts
	return &out, nil
}

// TermsDelete - delete term of invitation
func (s *Storage) TermsDelete(ctx context.Context, tID primitive.ObjectID) error {
	res, err := s.terms().DeleteOne(ctx, bson.M{"_id": tID})
	if err != nil {
		return errors.Wrapf(err, "error on delete term (id: %v)", tID)
	}
	if res.DeletedCount < 1 {
		return errors.Errorf("term is not deleted (id: %v)", tID)
	}
	return nil
}

func (s *Storage) terms() *mongo.Collection {
	return s.db().Collection("terms")
}
