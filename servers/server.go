package servers

import (
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
	router := fiber.New()
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Zombie Survival Social Network API")
	})
	return &Server{
		DB:     db,
		Router: router,
	}
}
