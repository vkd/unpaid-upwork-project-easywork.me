package storage

import (
	"context"

	"gitlab.com/easywork.me/backend/models"
	"gopkg.in/mgo.v2/bson"
)

func (s *Storage) ProjectGet(ctx context.Context, pID bson.ObjectId) (*models.Project, error) {
	panic("Not implemented")
}
