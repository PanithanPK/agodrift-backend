package handlers

import (
	"strconv"

	"agodrift/internal/model"
	"agodrift/internal/service"

	"github.com/gofiber/fiber/v2"
)

var tourService = service.NewTourService()

func Health(c *fiber.Ctx) error {
	return c.SendString("OK")
}

func ListRoomsHandler(c *fiber.Ctx) error {
	var t model.Tour
	if err := c.BodyParser(&t); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid body")
	}
	created := tourService.Create(t)
	return c.Status(fiber.StatusCreated).JSON(created)
}

func AddRoomHandler(c *fiber.Ctx) error {
	return c.JSON(tourService.List())
}

func RoomByIDHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid id")
	}
	t, found := tourService.Get(id)
	if !found {
		return c.Status(fiber.StatusNotFound).SendString("not found")
	}
	return c.JSON(t)
}
