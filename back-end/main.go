package main

import (
	"fmt"
	"log"
	"strconv"

	database "example.com/StuDuwo/db"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

// type Rental_post struct {
// 	House_name  string `json:"house_name"`
// 	Location    string `json:"location"`
// 	Start_date  string `json:"start_date"`
// 	Rent        string `json:"rent"`
// 	Description string `json:"description"`
// 	Image_URL   string `json:"image_url"`
// }

type Rental_post struct {
	ID_          string `json:"_id,omitempty" form:"_id,omitempty"`
	Listing_name string `json:"listing_name" form:"listing_name"`
	Owner_email  string `json:"owner_email" form:"owner_email"`
	Address1     string `json:"address1" form:"address1"`
	Address2     string `json:"address2" form:"address2"`
	Pincode      string `json:"pincode" form:"pincode"`
	Apt_img      string `json:"apt_img" form:"apt_img"`
	Price        string `json:"price" form:"price"`
	Rooms        string `json:"rooms" form :"rooms"`
}

func total_listings(c *fiber.Ctx) error { // Get total number of listings

	collection := database.Get_Collection("rentals")
	count, _ := collection.CountDocuments(c.Context(), nil)

	//fmt.Println(count)
	response := struct {
		Total_listings int `json:"total_listings"`
	}{Total_listings: int(count)}

	return c.JSON(response)
}

func get_all_listings(c *fiber.Ctx) error { // Modify this to enable pagination
	collection := database.Get_Collection("rentals")

	var rental_list []Rental_post

	cur, err := collection.Find(c.Context(), bson.M{})
	if err != nil {
		log.Println("Error finding documents: ", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error getting all rentals")
	}
	defer cur.Close(c.Context())

	if err := cur.All(c.Context(), &rental_list); err != nil {
		log.Println("Error making rental list: ", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error getting all rentals")
	}
	//fmt.Println(rental_list)
	// Send only required page
	spage := c.Params("*")
	page, err := strconv.Atoi(spage)
	fmt.Println(page)
	var page_list []Rental_post
	page_size := 20
	start_idx := (page - 1) * page_size
	last_idx := start_idx + page_size

	if last_idx >= 0 && last_idx < len(rental_list) { // if listings more than page
		page_list = rental_list[start_idx:last_idx]
	} else {
		page_list = rental_list[start_idx:]
	}

	return c.JSON(page_list)
}

func post_new_rental(c *fiber.Ctx) error {
	// sample_post := Rental_post{
	//	ID_:          "1"
	// 	Listing_name: "Uilenstede",
	// 	Owner_email:  "owner@example.com",
	// 	Address1:     "123 Main St",
	// 	Address2:     "Apt 456",
	// 	Pincode:      "12345",
	// 	Apt_img:      "apt_image.jpg",
	// 	Price:        "$1000",
	// 	Rooms:        "3",
	// }
	var new_post Rental_post
	if err := c.BodyParser(new_post); err != nil { //multipart form data
		return err
	}
	collection := database.Get_Collection("rentals")
	// new_post.ID_ = generateUniqueID()
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

	app.Post("/new_listing", post_new_rental)

	app.Get("/listings/*", get_all_listings)

	app.Get("/total_listings", total_listings)

	app.Listen(":5000")

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
