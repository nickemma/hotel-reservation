package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nickemma/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

// constructor function
func NewHotelHandler(store db.Store) *HotelHandler {
	return &HotelHandler{
		store: &store,
	}
}

/*
 * @route   GET api/v1/hotels
 * @desc    Get all hotels
 * @access  Public
 */

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	hotels, err := h.store.Hotel.GetHotels(c.Context(), nil)
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}

/*
 * @route   GET api/v1/hotel/:id/room
 * @desc    Get rooms by id
 * @access  Private
 */

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"hotelId": oid}
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return err
	}

	return c.JSON(rooms)
}

/*
 * @route   GET api/v1/hotels/:id
 * @desc    Get hotel by id
 * @access  Private
 */

func (h *HotelHandler) HandleGetHotelById(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), oid)
	if err != nil {
		return err
	}

	return c.JSON(hotel)
}
