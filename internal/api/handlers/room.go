package handlers

import (
	"strconv"

	"agodrift/internal/model"
	"agodrift/internal/service"

	"github.com/gofiber/fiber/v2"
)

var roomService = service.NewRoomService()

func Health(c *fiber.Ctx) error {
	return c.SendString("OK")
}

func ListRoomsHandler(c *fiber.Ctx) error {
	return c.JSON(roomService.List())
}

func AddRoomHandler(c *fiber.Ctx) error {
	var r model.Room
	if err := c.BodyParser(&r); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid body")
	}
	created := roomService.Create(r)
	return c.Status(fiber.StatusCreated).JSON(created)
}

func RoomByIDHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid id")
	}
	r, found := roomService.Get(id)
	if !found {
		return c.Status(fiber.StatusNotFound).SendString("not found")
	}
	return c.JSON(r)
}
