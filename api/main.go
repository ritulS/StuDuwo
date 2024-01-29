package main

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"os"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type Rental struct {
	ID       string  `json:"id,omitempty" form:"id,omitempty" gorm:"primaryKey;type:varchar(64)"`
	Name     string  `json:"listing_name" form:"listing_name"`
	Email    string  `json:"owner_email" form:"owner_email"`
	Address1 string  `json:"address1" form:"address1"`
	Address2 string  `json:"address2" form:"address2"`
	PinCode  string  `json:"pincode" form:"pincode"`
	ImgId    string  `json:"img_id" form:"img_id" gorm:"type:varchar(64)"`
	Price    float32 `json:"price" form:"price"`
	Rooms    uint8   `json:"rooms" form:"rooms"`
}

func main() {
	db = Connect()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			os.Stderr.WriteString(fmt.Sprintf("panic: %v\n%s\n", e, string(debug.Stack())))
		},
	}))

	app.Use(logger.New(logger.Config{
		Format:     "${pid} ${ip} ${locals:requestid} ${status} ${latency} - ${method} ${path}\n",
		TimeFormat: time.RFC3339,
		Output:     os.Stdout,
	}))

	app.Post("/new_listing", post_new_rental)
	app.Get("/listings/:page", get_all_listings)
	app.Get("/total_listings", total_listings)
	app.Get("/image/:img_id", get_image)
	app.Listen(":5000")
}

func Connect() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=5432 sslmode=disable", os.Getenv("URI"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Rental{})
	return db
}

type Listings struct {
	Listings []Rental `json:"listings"`
}

func total_listings(c *fiber.Ctx) error {
	var count int64
	res := db.Model(&Rental{}).Count(&count)
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(res.Error)
	}

	response := struct {
		Total_listings int `json:"total_listings"`
	}{Total_listings: int(count)}

	return c.JSON(response)
}

func get_all_listings(c *fiber.Ctx) error { // Modify this to enable pagination
	rental_list := make([]Rental, 0)
	page, _ := strconv.Atoi(c.Params("page"))
	if page <= 0 {
		page = 1
	}

	pageSize := 20
	offset := (page - 1) * pageSize
	res := db.Model(&Rental{}).Offset(offset).Limit(pageSize).Find(&rental_list)
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(res.Error)
	}

	return c.JSON(Listings{
		Listings: rental_list,
	})
}

func post_new_rental(c *fiber.Ctx) error {
	new_post := new(Rental)
	img, err := c.FormFile("img_id")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	img_file, err := img.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	if err := c.BodyParser(new_post); err != nil { //multipart form data
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	new_post.ID, _ = GenerateRandomString(64)
	new_post.ImgId, _ = GenerateRandomString(64)
	res := db.Create(&new_post)
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(res.Error)
	}

	u, _ := url.Parse("http://seaweed-proxy-service/")
	u = u.JoinPath(new_post.ImgId)
	agent := fiber.Post(u.String())
	agent.BodyStream(img_file, int(img.Size))

	statusCode, body, err_agent := agent.Bytes()
	if len(err_agent) > 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(err_agent)
	}

	if statusCode != 200 {
		return c.Status(fiber.StatusInternalServerError).JSON(body)
	}

	return c.JSON(new_post)
}

func get_image(c *fiber.Ctx) error {
	img_id := c.Params("img_id")

	u, _ := url.Parse("http://seaweed-proxy-service/")
	u = u.JoinPath(img_id)
	agent := fiber.Get(u.String())

	statusCode, body, err_agent := agent.Bytes()
	if len(err_agent) > 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(err_agent)
	}

	if statusCode != 200 {
		return c.Status(fiber.StatusInternalServerError).JSON(body)
	}

	return c.SendStream(bytes.NewReader((body)))
}
