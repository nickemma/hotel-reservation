package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("Hello, World!")
	app := fiber.New()
	app.Listen(":5000")
}
