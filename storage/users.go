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

// UserGetByIDOrEmail - get user by id or email
func (s *Storage) UserGetByIDOrEmail(ctx context.Context, uID models.UserID, email string) (*models.User, error) {
	if len(email) > 0 {
		return s.UserGetByEmail(ctx, email)
	}
	return s.UserGet(ctx, uID)
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

// UserDelete - delete user
func (s *Storage) UserDelete(ctx context.Context, userID models.UserID) error {
	res, err := s.users().DeleteOne(ctx, bson.M{"_id": userID})
	if err != nil {
		return errors.Wrapf(err, "error on delete user (id: %v)", userID)
	}
	if res.DeletedCount < 1 {
		return errors.Errorf("user not deleted (id: %v)", userID)
	}
	return nil
}

func (s *Storage) users() *mongo.Collection {
	return s.db().Collection("users")
}
