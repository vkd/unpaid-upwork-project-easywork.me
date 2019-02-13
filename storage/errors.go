package storage

import (
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/pkg/errors"
)

var ErrNotFound = errors.New("not found")
var ErrNoUpdated = errors.New("not updated")

func IsNotFound(err error) bool {
	err = errors.Cause(err)
	switch err {
	case ErrNotFound, mongo.ErrNoDocuments:
		return true
	}
	return false
}
