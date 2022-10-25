package servers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"zssn/domains/core"
	"zssn/domains/entities"
	"zssn/domains/users"
	"zssn/requests"
	"zssn/responses"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var userService users.IUserService

func (s *Server) userRoutes() {
	usr := s.Router.Group("/users")
	usr.Post("", newUser)
	usr.Post("/new-token", newToken)

	usr.Use("/me", authMiddleware())
	usr.Get("/me", userDetails)
	usr.Use("/flag", authMiddleware())
	usr.Post("/flag", flagInfectedUser)
	usr.Use("/location", authMiddleware())
	usr.Patch("/location", updateLocation)

}

func updateLocation(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)
	if userID == "" {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "invalid user ID",
		})
	}

	var req *requests.UpdateLocation
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	user, err := userService.Find(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	if user.Infected {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "user is infected and thus you cannot perform this operation",
		})
	}

	err = userService.UpdateLocation(ctx.Context(), userID, req.Latitude, req.Longitude)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "location updated successfully",
	})
}

func newUser(ctx *fiber.Ctx) error {
	var u *requests.Survivor
	if err := json.Unmarshal(ctx.Body(), &u); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	if err := u.Validate(); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	user := &entities.User{
		Email:     u.Email,
		Name:      u.Name,
		Age:       u.Age,
		Gender:    u.Gender,
		Latitude:  u.Latitude,
		Longitude: u.Longitude,
	}
	if err := userService.Create(ctx.Context(), user); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	var invItems []*entities.Inventory
	for _, v := range u.Inventory {
		invItems = append(invItems, &entities.Inventory{
			UserID:   user.ID,
			Item:     v.Item,
			Quantity: v.Quantity,
		})
	}

	if err := inventoryService.Create(ctx.Context(), invItems); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	td := core.TokenData{
		UserID: user.ID,
		Email:  user.Email,
	}
	token, err := td.Generate()
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return ctx.Status(http.StatusCreated).JSON(responses.FromUserEntity(user, token))
}

func userDetails(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)
	if userID == "" {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "invalid user ID",
		})
	}
	// get the user
	user, err := userService.Find(ctx.Context(), userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "invalid user",
			})
		}
	}

	balance, err := inventoryService.FindUserInventory(ctx.Context(), user.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "invalid user",
			})
		}
	}

	td := core.TokenData{
		UserID: user.ID,
		Email:  user.Email,
	}
	tk, err := td.Generate()
	if err != nil {
		return err
	}
	resp := responses.FromUserEntity(user, tk)
	for _, v := range balance {
		resp.Inventory = append(resp.Inventory, &responses.Inventory{
			Item:     strings.ToLower(v.Item.String()),
			Quantity: v.Quantity,
			Balance:  v.Balance,
		})
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}

func newToken(ctx *fiber.Ctx) error {
	var f *requests.NewToken
	if err := json.Unmarshal(ctx.Body(), &f); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	user, err := userService.FindByEmail(ctx.Context(), f.Email)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	td := &core.TokenData{
		UserID: user.ID,
		Email:  user.Email,
	}
	token, err := td.Generate()
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(responses.FromUserEntity(user, token))
}

func flagInfectedUser(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(string)
	if userID == "" {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "invalid user ID",
		})
	}
	var f *requests.FlagUser
	if err := json.Unmarshal(ctx.Body(), &f); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	err := userService.FlagUser(ctx.Context(), userID, f.InfectedUserID)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	details, err := userService.Find(ctx.Context(), f.InfectedUserID)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	// if the user is infected, we want to make all inventory items inaccessible
	if details.Infected {
		if err := inventoryService.BlockUserInventory(ctx.Context(), details.ID); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "user flagged successfully",
	})
}
