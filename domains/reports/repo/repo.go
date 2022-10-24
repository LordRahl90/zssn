package repo

import (
	"context"
	"zssn/domains/core"
	"zssn/domains/entities"

	"gorm.io/gorm"
)

var (
	_ IReportRepository = (*ReportRepository)(nil)
)

// ReportRepository implements the IResportRepository interface
type ReportRepository struct {
	DB *gorm.DB
}

// New returns a new instance of IReportRepository
func New(db *gorm.DB) IReportRepository {
	return &ReportRepository{
		DB: db,
	}
}

// Total gets the total number of users in the system.
// This will help us given the db does cache some query results
func (rr *ReportRepository) Total(ctx context.Context) (uint32, error) {
	var total int64

	err := rr.DB.Table("users").Count(&total).Error
	if err != nil {
		return 0, err
	}
	return uint32(total), nil

}

// Infected implements IIReportRepository
func (rr *ReportRepository) Infected(ctx context.Context) (*entities.Infected, error) {
	var (
		infected   int64
		percentage float64
	)

	total, err := rr.Total(ctx)
	if err != nil {
		return nil, err
	}

	err = rr.DB.Table("users").Where("infected = ?", true).Count(&infected).Error
	if err != nil {
		return nil, err
	}

	if total != 0 && infected != 0 {
		percentage = float64(infected) / float64(total) * 100.0
	}

	return &entities.Infected{
		Total:      uint32(total),
		Infected:   uint32(infected),
		Percentage: percentage,
	}, nil
}

// Points returns the accumulated points for each given item
func (rr *ReportRepository) Points(ctx context.Context) (map[core.Item]*entities.Resource, error) {
	var (
		dbResult []*entities.Resource
		result   = make(map[core.Item]*entities.Resource)
	)
	err := rr.DB.Table("inventories").Where("is_accessible = ?", false).Find(&dbResult).Error
	if err != nil {
		return nil, err
	}

	for _, v := range dbResult {
		data, ok := result[v.Item]
		if !ok {
			result[v.Item] = v
			continue
		}
		data.Balance += v.Balance
		result[v.Item] = data
	}

	return result, nil
}

// Resources calculates the total amount of resources available for each
func (rr *ReportRepository) Resources(ctx context.Context) (map[core.Item]*entities.Resource, error) {
	var (
		dbResult []*entities.Resource
		result   = make(map[core.Item]*entities.Resource)
	)
	err := rr.DB.Debug().Table("inventories").Where("is_accessible = ?", true).Find(&dbResult).Error
	if err != nil {
		return nil, err
	}

	for _, v := range dbResult {
		data, ok := result[v.Item]
		if !ok {
			result[v.Item] = v
			continue
		}
		data.Balance += v.Balance
		result[v.Item] = data
	}

	return result, nil
}

// Survivors returns the rate of survivors
func (rr *ReportRepository) Survivors(ctx context.Context) (*entities.Survivor, error) {
	var (
		clean      int64
		percentage float64
	)

	total, err := rr.Total(ctx)
	if err != nil {
		return nil, err
	}

	err = rr.DB.Table("users").Where("infected = false").Count(&clean).Error
	if err != nil {
		return nil, err
	}

	if total > 0 && clean > 0 {
		percentage = float64(clean) / float64(total) * 100.0
	}

	return &entities.Survivor{
		Total:      uint32(total),
		Clean:      uint32(clean),
		Percentage: percentage,
	}, nil
}
