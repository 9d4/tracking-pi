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

	//save photo now if available
	if err := os.MkdirAll(ModelDir, fs.ModeDir|fs.ModePerm); err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	file, err := os.Create(filepath.Join(ModelDir, resultID.Hex()))
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	defer file.Close()

	_, err = file.WriteString(vol.Photo)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	c.Status(fiber.StatusCreated)
	return c.JSON(vol)
}

func HandleVolunteerStorePhoto(c *fiber.Ctx) error {
	type req struct {
		Photo string `json:"photo"`
	}
	var r req
	err := c.BodyParser(&r)
	if err != nil {
		return err
	}

	if c.Params("id") == "" || r.Photo == "" {
		return fiber.ErrBadRequest
	}

	if err := os.MkdirAll(ModelDir, fs.ModeDir|fs.ModePerm); err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	file, err := os.Create(filepath.Join(ModelDir, c.Params("id")))
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	defer file.Close()

	_, err = file.WriteString(r.Photo)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	c.Status(fiber.StatusCreated)
	return nil
}

func HandleVolunteerDelete(c *fiber.Ctx) error {
	type reqDelete struct {
		IDs []string `json:"ids"`
	}
	var req reqDelete
	err := c.BodyParser(&req)
	if err != nil {
		return fiber.ErrBadRequest
	}
	log.Println("DELETE:", req)

	result, err := volunteer.GetStore().Delete(req.IDs)
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{
		"n": result.DeletedCount,
	})
}

func HandleVolunteerIndex(c *fiber.Ctx) error {
	vols, err := volunteer.GetStore().GetAll()
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(vols)
}
