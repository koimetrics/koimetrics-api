package api

import (
	"context"
	"fmt"
	"log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client mongo.Client
var analytics *mongo.Collection
var apikeys *mongo.Collection
var domains *mongo.Collection

func DBConnection() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	analytics = client.Database("test").Collection("analytics")
	apikeys = client.Database("test").Collection("apikeys")
	fmt.Println("collection: ", analytics)
}

type AnalyticResult struct {
	//Server stats
	Key         	string
	Session_id		string
	Session_start	string
	Session_end		string		

	// Client stats
	Host         	string
	Path         	string
	Date         	string
	Referrer	 	string
	ReferrerPath 	string
	Time         	string
	Performance  	float64
	Latitude     	float64
	Longitude    	float64
	IsPhone      	bool
	Country      	string
	City         	string
	Region       	string
}

type ApiKey struct {
	Key           	string
	Websites	  	[]string
	AskLocationTo 	[]string
	EndDate 	  	string
}
