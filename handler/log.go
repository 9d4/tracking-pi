package handler

import (
	logg "github.com/9d4/tracking-pi/log"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	result, err := logg.GetStore().Create(&lo)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	if objectID, ok := result.InsertedID.(primitive.ObjectID); ok {
		go logg.ProcessLogResult(objectID)
	}

	c.Status(201)
	return c.JSON(lo)
}
