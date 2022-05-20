package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/brunoos/cnterra-controller/config"
)

var DB *mongo.Database

func Initialize() {

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/",
		config.DbUser, config.DbPassword, config.DbAddress, config.DbPort)

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalln("[ERRO] Error connecting to mongodb", err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		log.Fatalln("[ERRO] Error connecting to mongodb", err)
	}

	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		log.Fatalln("[ERRO] Error connecting to mongodb", err)
	}

	DB = client.Database(config.DbName)
}
