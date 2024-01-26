package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyendpoints "github.com/aws/smithy-go/endpoints"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var s3_client *s3.Client
var bucket = "images"

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
	s3_client = ConnectS3()

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

	// cur, err := collection.Find(c.Context(), bson.M{})
	// if err != nil {
	// 	log.Println("Error finding documents: ", err)
	// 	return c.Status(fiber.StatusInternalServerError).SendString("Error getting all rentals")
	// }
	// defer cur.Close(c.Context())

	// if err := cur.All(c.Context(), &rental_list); err != nil {
	// 	log.Println("Error making rental list: ", err)
	// 	return c.Status(fiber.StatusInternalServerError).SendString("Error getting all rentals")
	// }

	// if len(rental_list) == 0 {
	// 	return c.JSON(Listings{
	// 		Listings: rental_list,
	// 	})
	// }

	// spage := c.Params("*")
	// page, err := strconv.Atoi(spage)
	// fmt.Println(page)
	// var page_list []Rental
	// page_size := 20
	// start_idx := page * page_size
	// last_idx := start_idx + page_size

	// if last_idx >= 0 && last_idx < len(rental_list) { // if listings more than page
	// 	page_list = rental_list[start_idx:last_idx]
	// } else {
	// 	page_list = rental_list[start_idx:]
	// }

	return c.JSON(Listings{
		Listings: rental_list,
	})
}

func post_new_rental(c *fiber.Ctx) error {
	// sample_post := Rental{
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

	_, err = s3_client.PutObject(c.Context(), &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &new_post.ImgId,
		Body:   img_file,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.JSON(new_post)
}

type resolverV2 struct{}

func (*resolverV2) ResolveEndpoint(ctx context.Context, params s3.EndpointParameters) (
	smithyendpoints.Endpoint, error,
) {
	u, err := url.Parse(fmt.Sprintf("http://%s:8333/", os.Getenv("SEAWEEDFS_S3")))
	if err != nil {
		return smithyendpoints.Endpoint{}, err
	}
	return smithyendpoints.Endpoint{
		URI: *u,
	}, nil
}

func ConnectS3() *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	uri := fmt.Sprintf("http://%s:8333/", os.Getenv("SEAWEEDFS_S3"))
	s3_client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = &uri
		o.UsePathStyle = true
	})

	_, err = s3_client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: &bucket,
	})

	out, err := s3_client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	for i := range out.Buckets {
		log.Println(out.Buckets[i].Name)
	}

	log.Println(err)
	return s3_client
}
