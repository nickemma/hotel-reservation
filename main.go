package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nickemma/hotel-reservation/api"
	"github.com/nickemma/hotel-reservation/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoUri = "mongodb://localhost:27017"
const dbname = "hotel-reservation"
const userCollection = "users"

func main() {
	// connection to mongodb
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))

	// check if there is an error while connecting to mongodb
	if err != nil {
		log.Fatal(err)
	}

	// check if the connection is successful
	ctx := context.Background()

	// collection instance to perform operations
	collection := client.Database(dbname).Collection(userCollection)

	// insert a user
	user := types.User{
		FirstName: "Peter",
		LastName:  "Parker",
	}
	// create a new user in the database
	_, err = collection.InsertOne(ctx, user)

	// check if there is an error while inserting a user
	if err != nil {
		log.Fatal(err)
	}

	// find one user
	var userFound types.User
	if err := collection.FindOne(ctx, bson.M{}).Decode(&userFound); err != nil {
		log.Fatal(err)
	}
	// print the result
	fmt.Println(userFound)

	// listen address flag to run the server
	listenAddr := flag.String("listenAddr", ":5000", "Api server listener")
	// parse the flags to get the listen address
	flag.Parse()

	// create a new fiber app
	app := fiber.New()
	// create a new api instance for version 1
	apiV1 := app.Group("api/v1/")

	// routes for the api
	apiV1.Get("/users", api.HandleGetUsers)
	apiV1.Get("/users/:id", api.HandleGetUserById)

	// Listen address
	app.Listen(*listenAddr)
}
