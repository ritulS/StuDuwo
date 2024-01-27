package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"runtime/debug"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var s3_client *s3.Client
var bucket = "images"

func main() {
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

	app.Post("/:img_id", create_image)
	app.Get("/:img_id", get_image)
	app.Listen(":5000")
}

func create_image(c *fiber.Ctx) error {
	img_id := c.Params("img_id")

	_, err := s3_client.PutObject(c.Context(), &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &img_id,
		Body:   bytes.NewReader(c.Body()),
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.SendStatus(200)
}

func get_image(c *fiber.Ctx) error {
	img_id := c.Params("img_id")
	out, err := s3_client.GetObject(c.Context(), &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &img_id,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.SendStream(out.Body, int(*out.ContentLength))
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

	if err != nil {
		log.Println(err)
	}
	return s3_client
}
