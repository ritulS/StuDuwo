package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {

	//fmt.Println("Hello world")

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello world")
	})

	app.Listen(":3000")

}
