package model

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var db *mongo.Client
var ctx context.Context
var cancel func()

func InitDB() {
	var err error

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	dbHost := os.Getenv("DB_Host")
	db, err = mongo.Connect(ctx, options.Client().ApplyURI(dbHost))

	if err != nil {
		log.Fatalln(err, "Host file is", dbHost)
	}

	// Ping mongo database if up
	go func() {
		for {
			if err := db.Ping(ctx, readpref.Primary()); err != nil {
				log.Fatalln(err)
			}
			time.Sleep(time.Second * 5)
		}
	}()
}
