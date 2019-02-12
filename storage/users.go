package storage

import (
	"context"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/pkg/errors"
	"gitlab.com/easywork.me/backend/models"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

// UserGet - get user
func (s *Storage) UserGet(ctx context.Context, uID models.UserID) (*models.User, error) {
	var u models.User
	err := s.users().FindOne(ctx, bson.M{"id": uID}).Decode(&u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, errors.Wrapf(err, "error on get user (id: %v)", uID)
	}
	return &u, nil
}

// UserGetByEmail - get user by email
func (s *Storage) UserGetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := s.users().FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, errors.Wrapf(err, "error on get user (email: %v)", email)
	}
	return &u, nil
}

// UserCreate - create user
func (s *Storage) UserCreate(ctx context.Context, u *models.UserPassword) (*models.User, error) {
	if u.Role == "" {
		u.Role = models.Work
	}
	if err := models.CheckRole(u.Role); err != nil {
		return nil, err
	}
	passBs, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(passBs)
	_, err := s.users().InsertOne(ctx, u)
	if err != nil {
		return nil, errors.Wrapf(err, "error on create user (email: %v)", u.Email)
	}
	return &u.User, nil
}

func (s *Storage) users() *mongo.Collection {
	return s.db().Collection("users")
}
