package database

import "context"

type Database interface {
	// Save to DB.
	Save(ctx context.Context, payload interface{}, out interface{}) error
	// Save many contents to DB.
	SaveMany(ctx context.Context, payload []interface{}) error
	// Find content in DB by identifier.
	FindByID(ctx context.Context, id interface{}, projection map[string]interface{}, result interface{}) error
	// Find content with specified field.
	FindOne(ctx context.Context, fields, projection map[string]interface{}, result interface{}) error
	// Find many contents with specified fields.
	FindMany(ctx context.Context, fields, projection, sort map[string]interface{}, limit, skip int64, results interface{}) error
	// Update contents by ID.
	UpdateByID(ctx context.Context, id interface{}, payload interface{}) error
	// Update content with specified field.
	UpdateOne(ctx context.Context, fields map[string]interface{}, payload interface{}) error
	// Update content if found, else update.
	UpsertOne(ctx context.Context, fields map[string]interface{}, payload interface{}) error
	// Update many contents with payload with specified fields.
	UpdateMany(ctx context.Context, fields, payload map[string]interface{}) error
	// Atomically increment fields in DB.
	IncCount(ctx context.Context, fields map[string]interface{}, payload interface{}) error
	// Delete content from DB with ID.
	DeleteByID(ctx context.Context, id interface{}) error
	// Completely delete content from DB with ID.
	CompleteDeleteByID(ctx context.Context, id interface{}) error
	// Completely delete content with specified field.
	CompleteDeleteOne(ctx context.Context, fields map[string]interface{}) error
	// Find number of contents of fields in DB.
	Count(ctx context.Context, fields map[string]interface{}) (int64, error)
}
