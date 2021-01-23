package database

import "context"

// Database interface for multiple DB instance
type Database interface {
	Save(ctx context.Context, payload, options interface{}) error
	SaveMany(ctx context.Context, payload []interface{}, option interface{}) error
	FindAndUpdateOne(ctx context.Context, filter interface{}, payload, option interface{}) error
	UpdateMany(ctx context.Context, filter interface{}, payload, option interface{}) error
	DeleteOne(ctx context.Context, filter, option interface{}) error
	DeleteMany(ctx context.Context, filter, option interface{}) error
	FindOne(ctx context.Context, filter, projection interface{}, result interface{}) error
	FindMany(ctx context.Context, filter, projection, results interface{}) error
}
