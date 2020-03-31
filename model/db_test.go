package model

import (
	"context"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestEnvironmentVariable(t *testing.T) {
	if os.Getenv("DB_Action_Host") == "" {
		t.Error("Environment variable not set")
	}
}

func TestDatabase(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DB_Action_Host")))
	if err != nil {
		t.Errorf(err.Error())
	}
}
