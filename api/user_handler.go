package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nickemma/hotel-reservation/types"
)

func HandleGetUsers(c *fiber.Ctx) error {
	user := types.User{
		FirstName: "awesome",
		LastName:  "Doe",
	}
	return c.JSON(user)
}

func HandleGetUserById(c *fiber.Ctx) error {
	return c.JSON("james bond")
}
