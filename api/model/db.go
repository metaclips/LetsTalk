package model

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/metaclips/LetsTalk/api/values"
)

var db *mongo.Database
var ctx, cancel = context.WithCancel(context.Background())

func InitDB() {
	dbHost := values.Config.DbHost
	mongoDB, err := mongo.Connect(ctx, options.Client().ApplyURI(dbHost))
	if err != nil {
		cancel()
		log.Fatalln("could not connect to database", err)
	}

	db = mongoDB.Database(values.Config.DbName)

	// Ping mongo database continuously if down.
	go func() {
		for {
			if err := mongoDB.Ping(ctx, readpref.Primary()); err != nil {
				cancel()
				log.Fatalln("Database down", err)
			}

			time.Sleep(time.Second * 5)
		}
	}()

	getCollection := func(collection string, content interface{}) {
		result, err := db.Collection(collection).Find(ctx, bson.D{})
		if err != nil {
			log.Fatalln("error while getting collection", err)
		}

		err = result.All(ctx, content)
		if err != nil {
			log.Fatalln("error getting collection results", err)
		}
	}

	var users []User
	getCollection(values.UsersCollectionName, &users)

	values.MapEmailToName.Mutex.Lock()
	for _, user := range users {
		values.MapEmailToName.Mapper[user.Email] = user.Name
	}
	values.MapEmailToName.Mutex.Unlock()
}
