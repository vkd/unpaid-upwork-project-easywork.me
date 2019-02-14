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

// ProjectsGet - get projects
func (s *Storage) ProjectsGet(ctx context.Context, userID models.UserID) ([]models.Project, error) {
	var out = []models.Project{}

	cur, err := s.projects().Find(ctx, bson.M{"owner_id": userID})
	if err != nil {
		return nil, errors.Wrapf(err, "error on get projects")
	}

	for cur.Next(ctx) {
		var p models.Project
		err = cur.Decode(&p)
		if err != nil {
			return nil, errors.Wrapf(err, "error on decode project")
		}
		out = append(out, p)
	}

	if err = cur.Err(); err != nil {
		return nil, errors.Wrapf(err, "error on iterate projects")
	}

	return out, nil
}

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

// ProjectDelete - delete project
func (s *Storage) ProjectDelete(ctx context.Context, pID primitive.ObjectID, userID models.UserID) error {
	res, err := s.projects().DeleteOne(ctx, bson.M{"_id": pID, "owner_id": userID})
	if err != nil {
		return errors.Wrapf(err, "error on delete project (id: %v)", pID)
	}
	if res.DeletedCount < 1 {
		return errors.Errorf("project is not deleted (id: %v)", pID)
	}
	return nil
}

func (s *Storage) projects() *mongo.Collection {
	return s.db().Collection("projects")
}
