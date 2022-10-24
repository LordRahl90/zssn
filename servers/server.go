package servers

import (
	"zssn/domains/inventory"
	iinv "zssn/domains/inventory/store"
	"zssn/domains/reports"
	"zssn/domains/reports/repo"
	"zssn/domains/trade"
	itr "zssn/domains/trade/store"
	"zssn/domains/users"
	iusr "zssn/domains/users/store"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"gorm.io/gorm"
)

var (
	inventoryService inventory.IInventoryService
	tradeService     trade.ITradeService
	reportService    reports.IReportService
)

// Server contains the server properties that can be propagated across different services.
type Server struct {
	DB     *gorm.DB
	Router *fiber.App
}

// New creates a new instance of the server
func New(db *gorm.DB) (*Server, error) {
	router := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
	})
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Zombie Survival Social Network API")
	})

	router.Use(requestid.New())
	router.Use(cors.New())
	router.Use(logger.New())
	svr := &Server{
		DB:     db,
		Router: router,
	}
	if err := svr.setupServices(); err != nil {
		return nil, err
	}
	svr.userRoutes()
	svr.tradeRoutes()
	svr.reportRoutes()

	return svr, nil
}

func (s *Server) setupServices() error {
	st, err := iusr.New(s.DB)
	if err != nil {
		return err
	}
	usrSvc, err := users.New(st)
	if err != nil {
		return err
	}
	userService = usrSvc

	invStore, err := iinv.New(s.DB)
	if err != nil {
		return err
	}
	inventoryService = inventory.New(invStore)

	trStore, err := itr.New(s.DB)
	if err != nil {
		return err
	}
	tradeService = trade.New(trStore, userService, inventoryService)

	rpRepo := repo.New(s.DB)
	reportService = reports.New(rpRepo)

	return nil
}
