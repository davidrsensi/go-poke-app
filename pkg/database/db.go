package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	client *mongo.Client
}

type Config struct {
	Host string
	Name string
}

var mongoDB Client

// connect to the Db
func Init(cfg Config) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI(cfg.Name + "://" + cfg.Host + ":27017")
	mongoDB.client, _ = mongo.Connect(ctx, clientOptions)

}

// return collection
func Collection() *mongo.Collection {
	return mongoDB.client.Database("poke_app").Collection("pokes")

}
