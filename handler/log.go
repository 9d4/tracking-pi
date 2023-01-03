package handler

import (
	"encoding/hex"
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
		//go logg.ProcessLogResult(objectID)
		logg.GetQueue().Add(objectID, true)
	}

	c.Status(201)
	return c.JSON(lo)
}

func HandleLogsTrigger(c *fiber.Ctx) error {
	id := c.Params("id")
	if len(id) != 24 {
		return fiber.ErrBadRequest
	}

	idBytes, err := hex.DecodeString(id)
	if err != nil {
		return fiber.ErrBadRequest
	}

	objectID := primitive.ObjectID{}
	for i := 0; i < len(objectID); i++ {
		objectID[i] = idBytes[i]
	}

	go logg.ProcessLogResult(objectID)
	c.WriteString("Processing")
	return nil
}
