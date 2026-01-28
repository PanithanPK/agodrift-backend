package handlers

import (
	"time"

	"agodrift/internal/repository"
	"agodrift/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var bookingService = service.NewBookingService()

type CreateBookingRequest struct {
	HotelID  int    `json:"hotel_id"`
	CheckIn  string `json:"check_in"`
	CheckOut string `json:"check_out"`
	Adults   int    `json:"adults"`
	Children int    `json:"children"`
	Rooms    int    `json:"rooms"`
}

func CreateBooking(c *fiber.Ctx) error {
	user := c.Locals("user")
	if user == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	var uid int
	if tok, ok := user.(*jwt.Token); ok {
		if claims, ok := tok.Claims.(jwt.MapClaims); ok {
			uidFloat, ok := claims["uid"].(float64)
			if ok {
				uid = int(uidFloat)
			}
		}
	}
	if uid == 0 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var req CreateBookingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid body")
	}
	if req.HotelID == 0 {
		return c.Status(fiber.StatusBadRequest).SendString("hotel_id required")
	}
	if req.Rooms <= 0 {
		req.Rooms = 1
	}
	if req.Adults <= 0 {
		req.Adults = 1
	}

	checkIn, err := time.Parse("2006-01-02", req.CheckIn)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid check_in")
	}
	checkOut, err := time.Parse("2006-01-02", req.CheckOut)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid check_out")
	}
	if !checkOut.After(checkIn) {
		return c.Status(fiber.StatusBadRequest).SendString("check_out must be after check_in")
	}

	b, err := bookingService.Create(uid, req.HotelID, checkIn, checkOut, req.Adults, req.Children, req.Rooms)
	if err != nil {
		if err == repository.ErrNotEnoughRooms {
			return c.Status(fiber.StatusConflict).SendString("not enough rooms available")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("failed to create booking")
	}
	return c.Status(fiber.StatusCreated).JSON(b)
}

func ListMyBookings(c *fiber.Ctx) error {
	user := c.Locals("user")
	if user == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	var uid int
	if tok, ok := user.(*jwt.Token); ok {
		if claims, ok := tok.Claims.(jwt.MapClaims); ok {
			uidFloat, ok := claims["uid"].(float64)
			if ok {
				uid = int(uidFloat)
			}
		}
	}
	if uid == 0 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	list, err := bookingService.ListByUserID(uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("failed to list bookings")
	}
	return c.JSON(list)
}
