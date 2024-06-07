package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nickemma/hotel-reservation/api"
	"github.com/nickemma/hotel-reservation/db"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoUri = "mongodb://localhost:27017"

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	// listen address flag to run the server
	listenAddr := flag.String("listenAddr", ":5000", "Api server listener")
	// parse the flags to get the listen address
	flag.Parse()

	// connection to mongodb
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))

	// check if there is an error while connecting to mongodb
	if err != nil {
		log.Fatal(err)
	}
	// handlers initializations
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	// create a new fiber app
	app := fiber.New(config)
	// create a new api instance for version 1
	apiV1 := app.Group("api/v1/")

	// routes for the api
	apiV1.Post("/create", userHandler.HandleCreateUser)
	apiV1.Get("/users", userHandler.HandleGetUsers)
	apiV1.Get("/users/:id", userHandler.HandleGetUserById)

	// Listen address
	app.Listen(*listenAddr)
}
