package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nickemma/hotel-reservation/types"
)

/*
 * @route   GET api/v1/users
 * @desc    Get all users
 * @access  Public
 */

func HandleGetUsers(c *fiber.Ctx) error {
	user := types.User{
		FirstName: "awesome",
		LastName:  "Doe",
	}
	return c.JSON(user)
}

/*
 * @route   GET api/v1/users/:id
 * @desc    Get user by id
 * @access  Public
 */

func HandleGetUserById(c *fiber.Ctx) error {
	return c.JSON("james bond")
}
