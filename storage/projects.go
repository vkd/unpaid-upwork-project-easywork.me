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

// ProjectGetByOwner - get project
func (s *Storage) ProjectGetByOwner(ctx context.Context, pID primitive.ObjectID, userID models.UserID) (*models.Project, error) {
	var p models.Project
	err := s.projects().FindOne(ctx, bson.M{"_id": pID, "owner_id": userID}).Decode(&p)
	if err != nil {
		return nil, errors.Wrapf(err, "error on get project (id: %v)", pID)
	}
	return &p, nil
}

// ProjectCreate - create new project
func (s *Storage) ProjectCreate(ctx context.Context, p *models.ProjectBase, uID models.UserID) (*models.Project, error) {
	p.OwnerID = uID
	p.CreatedDateTime = time.Now()
	res, err := s.projects().InsertOne(ctx, p)
	if err != nil {
		return nil, errors.Wrapf(err, "error on create project")
	}
	var out models.Project
	out.ProjectBase = *p
	out.ID = res.InsertedID.(primitive.ObjectID)
	return &out, nil
}

func (s *Storage) projects() *mongo.Collection {
	return s.db().Collection("projects")
}
