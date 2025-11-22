package main

import (
	"context"
	"log"
	"task_manager/data"
	"task_manager/router"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	data.InitMongo(client)

	r := router.SetupRouter()
	r.Run(":8080")
}
