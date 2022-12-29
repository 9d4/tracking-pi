package handler

import (
	logg "github.com/9d4/tracking-pi/log"
	"github.com/gofiber/fiber/v2"
	"log"
)

const LogsPhotoDir = "data/photos/"

func HandleLogsStore(c *fiber.Ctx) error {
	var lo logg.Log

	err := c.BodyParser(&lo)
	if err != nil {
		log.Println(err)
		return fiber.ErrBadRequest
	}

	logg.GetStore().Create(&lo)

	c.Status(201)
	return c.JSON(lo)
}
