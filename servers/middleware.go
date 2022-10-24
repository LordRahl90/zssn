package servers

import (
	"net/http"
	"strings"

	"zssn/domains/core"

	"github.com/gofiber/fiber/v2"
)

func authMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		auth := ctx.Request().Header.Peek("Authorization")
		tokenString := strings.Split(string(auth), " ")
		if len(tokenString) != 2 {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "invalid token format",
			})
		}
		td, err := core.Decode(tokenString[1])
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}
		ctx.Locals("user_id", td.UserID)
		return ctx.Next()
	}
}
