package servers

import (
	"net/http"
	"strings"
	"zssn/domains/entities"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) reportRoutes() {
	rsr := s.Router.Group("/reports")
	rsr.Get("/survivors", nonInfectedSurvivor)
	rsr.Get("/infected", infectedSurvivor)
	rsr.Get("/lost-points", lostPoints)
	rsr.Get("/resources", averageResourceShare)
}

func infectedSurvivor(ctx *fiber.Ctx) error {
	res, err := reportService.InfectedSurvivors(ctx.Context())
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(res)
}

func nonInfectedSurvivor(ctx *fiber.Ctx) error {
	res, err := reportService.NonInfectedSurvivors(ctx.Context())
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(res)
}

func averageResourceShare(ctx *fiber.Ctx) error {
	res, err := reportService.ResourceSharing(ctx.Context())
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	var resp []*entities.ResourceSharing
	for _, v := range res {
		v.Item = strings.ToLower(v.Item)
		resp = append(resp, v)
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}

func lostPoints(ctx *fiber.Ctx) error {
	res, err := reportService.LostPoints(ctx.Context())
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
	})
}
