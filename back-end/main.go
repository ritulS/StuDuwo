package main

import (
	database "example.com/StuDuwo/back-end/db"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	//fmt.Println("Hello world")

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello world")
	})

	app.Listen(":3000")

}

func init_App() error {

	err := loadEnv()
	if err != nil {
		return err
	}

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
