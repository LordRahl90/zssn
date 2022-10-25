package reports

import (
	"context"

	"zssn/domains/core"
	"zssn/domains/entities"
	"zssn/domains/reports/repo"
)

var _ IReportService = (*ReportService)(nil)

// ReportService returns the implementation of report service
type ReportService struct {
	Repository repo.IReportRepository
}

// New returns a new service implementation
func New(repo repo.IReportRepository) IReportService {
	return &ReportService{
		Repository: repo,
	}
}

// InfectedSurvivors implements IReportService
func (rs *ReportService) InfectedSurvivors(ctx context.Context) (*entities.Infected, error) {
	return rs.Repository.Infected(ctx)
}

// LostPoints implements IReportService
func (rs *ReportService) LostPoints(ctx context.Context) (uint32, error) {
	var totalPoints uint32
	res, err := rs.Repository.Points(ctx)
	if err != nil {
		return totalPoints, err
	}

	for k, v := range res {
		pt := core.ItemPoints[k]
		totalPoints += pt * v.Balance
	}

	return totalPoints, nil
}

// NonInfectedSurvivors implements IReportService
func (rs *ReportService) NonInfectedSurvivors(ctx context.Context) (*entities.Survivor, error) {
	return rs.Repository.Survivors(ctx)
}

// ResourceSharing implements IReportService
func (rs *ReportService) ResourceSharing(ctx context.Context) (map[string]*entities.ResourceSharing, error) {
	surviors, err := rs.Repository.Survivors(ctx)
	if err != nil {
		return nil, err
	}
	resources, err := rs.Repository.Resources(ctx)
	if err != nil {
		return nil, err
	}
	result := make(map[string]*entities.ResourceSharing)

	for k, v := range resources {
		var pcr uint32
		if surviors.Clean > 0 {
			pcr = v.Balance / surviors.Clean
		}
		result[k.String()] = &entities.ResourceSharing{
			Item:        k.String(),
			Balance:     v.Balance,
			PerSurvivor: pcr,
		}
	}

	return result, nil
}
