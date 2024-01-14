package main

import (
	"context"

	database "example.com/StuDuwo/back-end/db"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

type Rental_post struct {
	House_name  string `json:"house_name"`
	Location    string `json:"location"`
	Start_date  string `json:"start_date"`
	Rent        int    `json:"rent"`
	Description string `json:"description"`
	Image_URL   string `json:"image_url"`
}

func get_all_rentals(c *fiber.Ctx) error {
	return nil
}

func post_new_rental(c *fiber.Ctx) error {
	sample_post := bson.M{"house_name": "Uilenstede",
		"location":    "Amstelveen",
		"start_date":  "Feb 1",
		"rent":        "500",
		"description": "Nice house",
		"image_url":   "dummy_url"}

	// if err := c.BodyParser(new_post); err != nil {
	// 	return err
	// }
	collection := database.Get_Collection("rentals")

	_, err := collection.InsertOne(context.TODO(), sample_post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error inserting new post")
	}
	return nil
}

func main() {

	err := init_App()
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello world")
	})

	app.Listen(":3000")

}

func init_App() error {
	// Setup env
	err := loadEnv()
	if err != nil {
		return err
	}

	// Initialize DB
	err = database.Init_db()
	if err != nil {
		return err
	}
	return nil

}

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}
