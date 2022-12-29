package handler

import (
	"github.com/9d4/tracking-pi/volunteer"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

const ModelDir = "data/models/"

func HandleVolunteerStore(c *fiber.Ctx) error {
	vol := volunteer.Volunteer{}
	err := c.BodyParser(&vol)
	if err != nil {
		log.Println(err)
		return fiber.ErrBadRequest
	}

	result, err := volunteer.GetStore().Create(&vol)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	resultID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	vol.ID = resultID

	c.Status(fiber.StatusCreated)
	return c.JSON(vol)
}

func HandleVolunteerStorePhoto(c *fiber.Ctx) error {
	fh, err := c.FormFile("photo")
	if err != nil {
		log.Println(err)
		return fiber.ErrBadRequest
	}

	path := ModelDir + c.Params("id")

	if err = os.MkdirAll(path, fs.ModeDir|fs.ModePerm); err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	err = c.SaveFile(fh, filepath.Join(path, fh.Filename))
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	c.Status(fiber.StatusCreated)
	return nil
}

func HandleVolunteerIndex(c *fiber.Ctx) error {
	vols, err := volunteer.GetStore().GetAll()
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(vols)
}
