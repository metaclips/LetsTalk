package database

import (
	"context"
	"encoding/json"
	"time"

	"github.com/metaclips/LetsTalk/backend/helpers/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type mongoDB struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewMonStore(db *mongo.Database, collection string) Database {

	ctx := context.TODO()

	_, _ = db.Collection(constants.RoomRequestCollection).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bsonx.Doc{bsonx.Elem{Key: "userID", Value: bsonx.Int32(1)}, bsonx.Elem{Key: "roomID", Value: bsonx.Int32(1)}},
		Options: options.Index().SetUnique(true),
	},
	)

	return &mongoDB{
		db,
		db.Collection(collection),
	}
}

func (d *mongoDB) Save(ctx context.Context, payload interface{}, out interface{}) error {
	result, err := d.collection.InsertOne(ctx, payload)

	if err != nil {
		return err
	}

	if out == nil {
		return nil
	}

	return d.collection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(out)
}

func (d *mongoDB) SaveMany(ctx context.Context, payload []interface{}) error {
	_, err := d.collection.InsertMany(ctx, payload)

	return err
}

func (d *mongoDB) FindByID(ctx context.Context, id interface{}, projection map[string]interface{}, result interface{}) error {
	ops := options.FindOne()
	if projection != nil {
		ops.Projection = projection
	}

	if err := d.collection.FindOne(ctx, bson.M{"_id": id, "deletedAt": nil}, ops).Decode(result); err != nil {
		return err
	}

	return nil
}

func (d *mongoDB) FindOne(ctx context.Context, fields, projection map[string]interface{}, result interface{}) error {
	ops := options.FindOne()
	ops.Projection = projection

	fields["deletedAt"] = nil

	err := d.collection.FindOne(ctx, fields, ops).Decode(result)
	return err
}

func (d *mongoDB) FindMany(ctx context.Context, fields, projection, sort map[string]interface{}, limit, skip int64, results interface{}) error {
	ops := options.Find()
	if limit > 0 {
		ops.Limit = &limit
	}
	if skip > 0 {
		ops.Skip = &skip
	}
	if projection != nil {
		ops.Projection = projection
	}
	if sort != nil {
		ops.Sort = sort
	}

	fields["deletedAt"] = nil

	cursor, err := d.collection.Find(ctx, fields, ops)
	if err != nil {
		return err
	}

	var output []map[string]interface{}
	for cursor.Next(ctx) {
		var item map[string]interface{}
		if err := cursor.Decode(&item); err != nil {
			return err
		}

		output = append(output, item)
	}

	if b, err := json.Marshal(output); err == nil {
		_ = json.Unmarshal(b, &results)
	} else {
		return err
	}

	return nil
}

func (d *mongoDB) UpsertOne(ctx context.Context, fields map[string]interface{}, payload interface{}) error {
	result := d.collection.FindOneAndUpdate(ctx, fields, bson.M{
		"$set": payload,
	}, options.FindOneAndUpdate().SetUpsert(true))

	// If result gives a document does not exist error, record is inserted.
	if err := result.Err(); err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	return nil
}

func (d *mongoDB) UpdateByID(ctx context.Context, id interface{}, payload interface{}) error {
	result := d.collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": payload}, nil)

	err := result.Err()
	if err != nil {
		return err
	}

	return nil
}

func (d *mongoDB) UpdateOne(ctx context.Context, fields map[string]interface{}, payload interface{}) error {
	return d.collection.FindOneAndUpdate(ctx, fields, bson.M{"$set": payload}).Err()
}

func (d *mongoDB) UpdateMany(ctx context.Context, fields, payload map[string]interface{}) error {
	if _, err := d.collection.UpdateMany(ctx, fields, bson.M{
		"$set": payload,
	}); err != nil {
		return err
	}

	return nil
}

func (d *mongoDB) IncCount(ctx context.Context, fields map[string]interface{}, payload interface{}) error {
	result := d.collection.FindOneAndUpdate(ctx, fields, bson.M{"$inc": payload})
	return result.Err()
}

func (d *mongoDB) DeleteByID(ctx context.Context, id interface{}) error {
	var u map[string]interface{}

	payload := bson.M{
		"deletedAt": time.Now(),
	}

	if err := d.collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{
		"$set": payload,
	}).Decode(&u); err != nil {
		return err
	}

	return nil
}

func (d *mongoDB) CompleteDeleteByID(ctx context.Context, id interface{}) error {
	var u map[string]interface{}
	if e := d.collection.FindOneAndDelete(ctx, bson.M{
		"_id": id,
	}).Decode(&u); e != nil {
		return e
	}

	return nil
}

func (d *mongoDB) CompleteDeleteOne(ctx context.Context, fields map[string]interface{}) error {
	_, err := d.collection.DeleteOne(ctx, fields)
	if err != nil {
		return err
	}

	return nil
}

func (d *mongoDB) Count(ctx context.Context, fields map[string]interface{}) (int64, error) {
	fields["deletedAt"] = nil

	return d.collection.CountDocuments(ctx, fields)
}
