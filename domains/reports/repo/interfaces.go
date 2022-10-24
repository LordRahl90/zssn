package repo

import (
	"context"
	"zssn/domains/core"
	"zssn/domains/entities"
)

// IIReportRepository interface to define the report contracts
type IReportRepository interface {
	Total(ctx context.Context) (uint32, error)
	Survivors(ctx context.Context) (*entities.Survivor, error)
	Infected(ctx context.Context) (*entities.Infected, error)
	Resources(ctx context.Context) (map[core.Item]*entities.Resource, error)
	Points(ctx context.Context) (map[core.Item]*entities.Resource, error)
}
