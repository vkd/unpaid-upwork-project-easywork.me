package storage

import (
	"context"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/pkg/errors"
	"gitlab.com/easywork.me/backend/models"
	"gopkg.in/mgo.v2/bson"
)

func (s *Storage) UserGet(ctx context.Context, uID models.UserID) (*models.User, error) {
	var u models.User
	err := s.users().FindOne(ctx, bson.M{"id": uID}).Decode(&u)
	if err != nil {
		return nil, errors.Wrapf(err, "error on get user (id: %v)", uID)
	}
	return &u, nil
}

func (s *Storage) users() *mongo.Collection {
	return s.db().Collection("users")
}
