package main

import (
	"fmt"
	"log"

	database "example.com/StuDuwo/db"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

type Rental_post struct {
	House_name  string `json:"house_name"`
	Location    string `json:"location"`
	Start_date  string `json:"start_date"`
	Rent        string `json:"rent"`
	Description string `json:"description"`
	Image_URL   string `json:"image_url"`
}

func get_all_rentals(c *fiber.Ctx) error {
	collection := database.Get_Collection("rentals")

	var rental_list []Rental_post

	cur, err := collection.Find(c.Context(), bson.M{})
	if err != nil {
		log.Println("Error finding documents: ", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error getting all rentals")
	}
	defer cur.Close(c.Context())
	fmt.Println("Cursor sucess")
	if err := cur.All(c.Context(), &rental_list); err != nil {
		log.Println("Error making rental list: ", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error getting all rentals")
	}
	fmt.Println(rental_list)
	return c.JSON(rental_list)
}

func post_new_rental(c *fiber.Ctx) error {
	// sample_post := bson.M{"house_name": "Uilenstede",
	// 	"location":    "Amstelveen",
	// 	"start_date":  "Feb 1",
	// 	"rent":        "500",
	// 	"description": "Nice house",
	// 	"image_url":   "dummy_url"}
	var new_post Rental_post
	if err := c.BodyParser(new_post); err != nil {
		return err
	}
	collection := database.Get_Collection("rentals")

	npost, err := collection.InsertOne(c.Context(), new_post)
	if err != nil {
		log.Println("Error inserting new post: ", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error inserting new post")
	}
	return c.JSON(npost)
}

func main() {

	err := init_App()
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	app.Post("/", post_new_rental)

	app.Get("/", get_all_rentals)

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
	fmt.Println("DB setup done")
	return nil

}

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}
