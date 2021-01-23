package database

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/metaclips/LetsTalk/api/pkg/appcontext"
	"github.com/metaclips/LetsTalk/api/values"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoStore struct {
	collection *mongo.Collection
}

var _ Database = &mongoStore{}

func CreateMongoInstance(ctx context.Context) (*mongo.Database, error) {
	mongoDB, err := mongo.Connect(ctx, options.Client().ApplyURI(values.Config.DbHost))
	if err != nil {
		log.Errorln("error connecting mongodb instance:", err)
		return nil, err
	}

	db := mongoDB.Database(values.Config.DbName)

	// Ping mongo database continuously if down.
	go func() {
		for {
			if err := mongoDB.Ping(ctx, readpref.Primary()); err != nil {
				// Call context cancel to all child context.
				appcontext.Cancel()
				log.Fatalln("Database down", err)
			}

			time.Sleep(time.Second * 5)
		}
	}()

	return db, nil
}

func NewMongoCollection(db *mongo.Database, collectionName string) Database {
	collection := db.Collection(collectionName)
	return &mongoStore{
		collection: collection,
	}
}

func (m *mongoStore) Save(ctx context.Context, payload, options interface{}) error {
	_, err := m.collection.InsertOne(ctx, payload)
	return err
}

func (m *mongoStore) SaveMany(ctx context.Context, payload []interface{}, option interface{}) error {
	opt, ok := option.(*options.InsertManyOptions)
	opts := make([]*options.InsertManyOptions, 0)
	if ok {
		opts = append(opts, opt)
	}

	_, err := m.collection.InsertMany(ctx, payload, opts...)
	return err
}

func (m *mongoStore) FindAndUpdateOne(ctx context.Context, filter, payload, option interface{}) error {
	opt := options.FindOneAndUpdate()
	opt.Projection = option

	result := m.collection.FindOneAndUpdate(ctx, filter, payload, opt)
	err := result.Err()

	return err
}

func (m *mongoStore) UpdateMany(ctx context.Context, filter interface{}, payload, option interface{}) error {
	opt, ok := option.(*options.UpdateOptions)
	opts := make([]*options.UpdateOptions, 0)
	if ok {
		opts = append(opts, opt)
	}

	_, err := m.collection.UpdateMany(ctx, filter, payload, opts...)
	return err
}

func (m *mongoStore) DeleteOne(ctx context.Context, filter, option interface{}) error {
	opt, ok := option.(*options.DeleteOptions)
	opts := make([]*options.DeleteOptions, 0)
	if ok {
		opts = append(opts, opt)
	}

	_, err := m.collection.DeleteOne(ctx, filter, opts...)
	return err
}

func (m *mongoStore) DeleteMany(ctx context.Context, filter, option interface{}) error {
	opt, ok := option.(*options.DeleteOptions)
	opts := make([]*options.DeleteOptions, 0)
	if ok {
		opts = append(opts, opt)
	}

	_, err := m.collection.DeleteMany(ctx, filter, opts...)
	return err
}

func (m *mongoStore) FindOne(ctx context.Context, filter, option, result interface{}) error {
	opt := options.FindOne()
	opt.Projection = option

	a := m.collection.FindOne(ctx, filter, opt)
	if err := a.Err(); err != nil {
		return err
	}

	return a.Decode(&result)
}

func (m *mongoStore) FindMany(ctx context.Context, filter, option, results interface{}) error {
	opt := options.Find()
	opt.Projection = option

	cursor, err := m.collection.Find(ctx, filter, opt)
	if err != nil {
		return err
	}

	return cursor.All(ctx, results)
}
