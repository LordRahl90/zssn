package reports

import (
	"context"
	"errors"

	"zssn/domains/core"
	"zssn/domains/entities"
	"zssn/domains/reports/repo"
)

var (
	_ repo.IReportRepository = (*MockReportRepository)(nil)

	errMockNotInitialized = errors.New("mock not initialized")
)

type MockReportRepository struct {
	InfectedFunc  func(ctx context.Context) (*entities.Infected, error)
	PointsFunc    func(ctx context.Context) (map[core.Item]*entities.Resource, error)
	ResourcesFunc func(ctx context.Context) (map[core.Item]*entities.Resource, error)
	SurvivorsFunc func(ctx context.Context) (*entities.Survivor, error)
	TotalFunc     func(ctx context.Context) (uint32, error)
}

// Infected implements repo.IReportRepository
func (m *MockReportRepository) Infected(ctx context.Context) (*entities.Infected, error) {
	if m.InfectedFunc == nil {
		return nil, errMockNotInitialized
	}
	return m.InfectedFunc(ctx)
}

// Points implements repo.IReportRepository
func (m *MockReportRepository) Points(ctx context.Context) (map[core.Item]*entities.Resource, error) {
	if m.PointsFunc == nil {
		return nil, errMockNotInitialized
	}
	return m.PointsFunc(ctx)
}

// Resources implements repo.IReportRepository
func (m *MockReportRepository) Resources(ctx context.Context) (map[core.Item]*entities.Resource, error) {
	if m.ResourcesFunc == nil {
		return nil, errMockNotInitialized
	}
	return m.ResourcesFunc(ctx)
}

// Survivors implements repo.IReportRepository
func (m *MockReportRepository) Survivors(ctx context.Context) (*entities.Survivor, error) {
	if m.SurvivorsFunc == nil {
		return nil, errMockNotInitialized
	}
	return m.SurvivorsFunc(ctx)
}

// Total implements repo.IReportRepository
func (m *MockReportRepository) Total(ctx context.Context) (uint32, error) {
	if m.TotalFunc == nil {
		return 0, errMockNotInitialized
	}

	return m.TotalFunc(ctx)
}
