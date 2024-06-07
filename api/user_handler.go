package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nickemma/hotel-reservation/db"
	"github.com/nickemma/hotel-reservation/types"
)

type UserHandler struct {
	userStore db.UserStore
}

// constructor function
func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

/*
 * @route   GET api/v1/users/:id
 * @desc    Get user by id
 * @access  Private
 */

func (h *UserHandler) HandleGetUserById(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)
	user, err := h.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

/*
 * @route   GET api/v1/users
 * @desc    Get all users
 * @access  Public
 */

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}

/*
 * @route   POST api/v1/create
 * @desc    Add a user to the database
 * @access  Public
 */

func (h *UserHandler) HandleCreateUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.Validate(); err != nil {
		return err
	}
	user, err := types.CreateNewUserFromParams(params)
	if err != nil {
		return err
	}
	insertedUser, err := h.userStore.CreateUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}
