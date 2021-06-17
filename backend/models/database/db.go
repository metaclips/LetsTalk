package database

import "context"

type Database interface {
	Save(ctx context.Context, payload interface{}, out interface{}) error
	SaveMany(ctx context.Context, payload []interface{}) error
	FindByID(ctx context.Context, id interface{}, projection map[string]interface{}, result interface{}) error
	FindOne(ctx context.Context, fields, projection map[string]interface{}, result interface{}) error
	FindMany(ctx context.Context, fields, projection, sort map[string]interface{}, limit, skip int64, results interface{}) error
	UpdateByID(ctx context.Context, id interface{}, payload interface{}) error
	UpdateOne(ctx context.Context, fields map[string]interface{}, payload interface{}) error
	UpsertOne(ctx context.Context, fields map[string]interface{}, payload interface{}) error
	UpdateMany(ctx context.Context, fields, payload map[string]interface{}) error
	IncCount(ctx context.Context, fields map[string]interface{}, payload interface{}) error
	DeleteByID(ctx context.Context, id interface{}) error
	CompleteDeleteByID(ctx context.Context, id interface{}) error
	CompleteDeleteOne(ctx context.Context, fields map[string]interface{}) error
	DeleteOne(ctx context.Context, fields map[string]interface{}) error
	DeleteMany(ctx context.Context, fields map[string]interface{}) error
	Aggregate(ctx context.Context, pipelines interface{}, result interface{}, allowDiskUse bool) error
	Count(ctx context.Context, fields map[string]interface{}) (int64, error)
}
