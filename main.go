package main

import (
	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/nickemma/hotel-reservation/api"
)

func main() {

	listenAddr := flag.String("listenAddr", ":5000", "Api server listener")
	flag.Parse()

	app := fiber.New()
	apiV1 := app.Group("api/v1/")

	// routes
	apiV1.Get("/user", api.HandleGetUsers)
	apiV1.Get("/user/:id", api.HandleGetUserById)

	// Listen address
	app.Listen(*listenAddr)
}
