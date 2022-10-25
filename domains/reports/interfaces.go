package reports

import (
	"context"

	"zssn/domains/entities"
)

// IReportService defines the expectations between the report and service
type IReportService interface {
	InfectedSurvivors(ctx context.Context) (*entities.Infected, error)
	NonInfectedSurvivors(ctx context.Context) (*entities.Survivor, error)
	ResourceSharing(ctx context.Context) (map[string]*entities.ResourceSharing, error)
	LostPoints(ctx context.Context) (uint32, error)
}
