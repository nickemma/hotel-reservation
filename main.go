package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nickemma/hotel-reservation/api"
	"github.com/nickemma/hotel-reservation/db"
	"github.com/nickemma/hotel-reservation/middleware"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.MongoUri))

	// check if there is an error while connecting to mongodb
	if err != nil {
		log.Fatal(err)
	}
	var (
		// handlers initializations
		hotelStore = db.NewMongoHotelStore(client)
		roomStore  = db.NewMongoRoomStore(client, hotelStore)
		userStore  = db.NewMongoUserStore(client)
		store      = &db.Store{
			Room:  roomStore,
			Hotel: hotelStore,
			User:  userStore,
		}
		userHandler  = api.NewUserHandler(userStore)
		authHandler  = api.NewAuthHandler(userStore)
		hotelHandler = api.NewHotelHandler(*store)

		// create a new fiber app
		app = fiber.New(config)
		// create a new api instance for version 1
		apiV1 = app.Group("api/v1/", middleware.JWTAuthentication)
		auth  = app.Group("api/")
	)

	// route for authentication
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// routes for the api User Handlers
	apiV1.Post("/create", userHandler.HandleCreateUser)
	apiV1.Get("/users", userHandler.HandleGetUsers)
	apiV1.Get("/users/:id", userHandler.HandleGetUserById)
	apiV1.Delete("/users/:id", userHandler.HandleDeleteUser)
	apiV1.Put("/users/:id", userHandler.HandleUpdateUser)

	// routes for the api User Handlers
	apiV1.Get("/hotels", hotelHandler.HandleGetHotel)
	apiV1.Get("/hotels/:id/rooms", hotelHandler.HandleGetRooms)
	apiV1.Get("/hotels/:id", hotelHandler.HandleGetHotelById)

	// Listen address
	app.Listen(*listenAddr)
}
