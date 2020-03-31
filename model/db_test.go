package model

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestEnvironmentVariable(t *testing.T) {
	fmt.Println(os.Getenv("DBHost"), 1)
	if os.Getenv("DBHost") == "" {
		t.Error("Environment variable not set")
	}
}

func TestDatabase(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DBHost")))
	if err != nil {
		log.Fatal(err)
	}
}
