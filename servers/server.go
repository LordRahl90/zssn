package servers

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Server contains the server properties that can be propagated across different services.
type Server struct {
	DB     *gorm.DB
	Router *fiber.App
}

// New creates a new instance of the server
func New(db *gorm.DB) *Server {
	router := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
	})
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Zombie Survival Social Network API")
	})

	router.Get("/health", func(c *fiber.Ctx) error {
		env := os.Getenv("ENVIRONMENT")
		return c.SendString("Server Environment " + env + " all green by " + time.Now().Format(time.RFC3339))
	})
	return &Server{
		DB:     db,
		Router: router,
	}
}
