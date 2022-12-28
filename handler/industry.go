package handler

import (
	"github.com/9d4/tracking-pi/industry"
	"github.com/gofiber/fiber/v2"
	"log"
)

func HandleStore(c *fiber.Ctx) error {
	ind := industry.Industry{}
	err := c.BodyParser(&ind)
	if err != nil {
		return fiber.ErrBadRequest
	}

	_, err = industry.GetStore().Create(&ind)
	if err != nil {
		log.Println(err)
		return err
	}

	c.Status(fiber.StatusCreated)
	return c.JSON(ind)
}

func HandleIndex(c *fiber.Ctx) error {
	industries, err := industry.GetStore().GetAll()
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(industries)
}
