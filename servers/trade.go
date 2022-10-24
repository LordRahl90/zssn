package servers

import (
	"encoding/json"
	"net/http"

	"zssn/requests"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) tradeRoutes() {
	tsr := s.Router.Group("/trades", authMiddleware())

	tsr.Post("", newTrade)
}

func newTrade(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)
	if userID == "" {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "invalid user ID",
		})
	}

	var tr *requests.TradeRequest
	if err := json.Unmarshal(ctx.Body(), &tr); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	seller := tr.Owner.ToServiceEntities()
	seller.UserID = userID // you cannot execute trades on behalf of another person
	buyer := tr.SecondParty.ToServiceEntities()

	if seller.UserID == buyer.UserID {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "you cannot trade with yourself",
		})
	}

	err := tradeService.Execute(ctx.Context(), buyer, seller)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	balance, err := inventoryService.FindUserInventory(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"balance": balance,
	})
}
