package storage

import (
	"context"

	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
	"github.com/pkg/errors"

	"github.com/mongodb/mongo-go-driver/mongo"
)

type Storage struct {
	client *mongo.Client
}

func NewMongoDB(ctx context.Context, uri string) (*Storage, error) {
	cli, err := mongo.NewClient(uri)
	if err != nil {
		return nil, errors.Wrapf(err, "error on create mongo client (uri: %v)", uri)
	}

	err = cli.Connect(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error on connect to mongo (uri: %v)", err)
	}

	err = cli.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, errors.Wrapf(err, "error on ping to mongo (uri: %v)", uri)
	}

	return &Storage{client: cli}, nil
}

func (s *Storage) Init(ctx context.Context) error {
	err := s.initCreateIndexes(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) db() *mongo.Database {
	return s.client.Database("easywork")
}

func (s *Storage) c(name string) *mongo.Collection {
	return s.db().Collection(name)
}

func (s *Storage) initCreateIndexes(ctx context.Context) error {
	type fs map[string]int
	for _, i := range []struct {
		db       string
		col      string
		fields   fs
		isUnique bool
	}{} {
		err := s.createIndex(ctx, i.db, i.col, i.fields, i.isUnique)
		if err != nil {
			return errors.Wrapf(err, "error on create index (%s.%s)", i.db, i.col)
		}
	}
	return nil
}

func (s *Storage) createIndex(ctx context.Context, db string, col string, fields map[string]int, isUnique bool) error {
	_, err := s.client.Database(db).Collection(col).Indexes().CreateOne(ctx, makeIndexModel(fields, isUnique))
	return err
}

func makeIndexModel(fields map[string]int, isUnique bool) mongo.IndexModel {
	var keys bsonx.Doc
	for name, order := range fields {
		keys = keys.Append(name, bsonx.Int32(int32(order)))
	}

	var im mongo.IndexModel
	im.Keys = keys
	if isUnique {
		im.Options = options.Index().SetUnique(isUnique)
	}
	return im
}
