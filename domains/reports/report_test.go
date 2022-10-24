package reports

import (
	"context"
	"os"
	"testing"

	"zssn/domains/core"
	"zssn/domains/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		os.Exit(code)
	}()

	code = m.Run()
}

func TestInfectedSurvivors(t *testing.T) {
	repo := &MockReportRepository{
		InfectedFunc: func(ctx context.Context) (*entities.Infected, error) {
			return &entities.Infected{
				Total:      10,
				Infected:   5,
				Percentage: 50.0,
			}, nil
		},
	}
	svc := New(repo)
	ctx := context.Background()
	res, err := svc.InfectedSurvivors(ctx)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, uint32(5), res.Infected)
	assert.Equal(t, uint32(10), res.Total)
	assert.Equal(t, float64(50.0), res.Percentage)
}

func TestInfectedWithErrorFromRepo(t *testing.T) {
	repo := &MockReportRepository{
		InfectedFunc: func(ctx context.Context) (*entities.Infected, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := New(repo)
	ctx := context.Background()
	res, err := svc.InfectedSurvivors(ctx)
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	require.Nil(t, res)
}

func TestInfectedWithErrorFromEmptyRepo(t *testing.T) {
	repo := &MockReportRepository{}
	svc := New(repo)
	ctx := context.Background()
	res, err := svc.InfectedSurvivors(ctx)
	require.EqualError(t, err, errMockNotInitialized.Error())
	require.Nil(t, res)
}

func TestNonInfectedSurvivors(t *testing.T) {
	repo := &MockReportRepository{
		SurvivorsFunc: func(ctx context.Context) (*entities.Survivor, error) {
			return &entities.Survivor{
				Total:      10,
				Clean:      5,
				Percentage: 50.0,
			}, nil
		},
	}
	svc := New(repo)
	ctx := context.Background()
	res, err := svc.NonInfectedSurvivors(ctx)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, uint32(5), res.Clean)
	assert.Equal(t, uint32(10), res.Total)
	assert.Equal(t, float64(50.0), res.Percentage)
}

func TestNonInfectedSurvivorsWithErrorFromRepo(t *testing.T) {
	repo := &MockReportRepository{
		SurvivorsFunc: func(ctx context.Context) (*entities.Survivor, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := New(repo)
	ctx := context.Background()
	res, err := svc.NonInfectedSurvivors(ctx)
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	require.Nil(t, res)
}

func TestNonInfectedSurvivorsWithEmptyRepo(t *testing.T) {
	repo := &MockReportRepository{}
	svc := New(repo)
	ctx := context.Background()
	res, err := svc.NonInfectedSurvivors(ctx)
	require.EqualError(t, err, errMockNotInitialized.Error())
	require.Nil(t, res)
}

func TestLostPoints(t *testing.T) {
	repo := &MockReportRepository{
		PointsFunc: func(ctx context.Context) (map[core.Item]*entities.Resource, error) {
			result := map[core.Item]*entities.Resource{
				core.ItemWater: {
					Item:    core.ItemWater,
					Balance: 50,
				},
				core.ItemMedication: {
					Item:    core.ItemMedication,
					Balance: 50,
				},
			}
			return result, nil
		},
	}
	svc := New(repo)
	ctx := context.Background()
	res, err := svc.LostPoints(ctx)
	require.NoError(t, err)
	assert.Equal(t, uint32(300), res)
}

func TestLostPointWithRepoError(t *testing.T) {
	repo := &MockReportRepository{
		PointsFunc: func(ctx context.Context) (map[core.Item]*entities.Resource, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := New(repo)
	ctx := context.Background()
	res, err := svc.LostPoints(ctx)
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	assert.Equal(t, uint32(0), res)
}

func TestLostPointWithUndefinedMock(t *testing.T) {
	repo := &MockReportRepository{}
	svc := New(repo)
	ctx := context.Background()
	res, err := svc.LostPoints(ctx)
	require.EqualError(t, err, errMockNotInitialized.Error())
	assert.Equal(t, uint32(0), res)
}

func TestResourceSharing(t *testing.T) {
	repo := &MockReportRepository{
		SurvivorsFunc: func(ctx context.Context) (*entities.Survivor, error) {
			return &entities.Survivor{
				Total:      10,
				Clean:      5,
				Percentage: 50.0,
			}, nil
		},
		ResourcesFunc: func(ctx context.Context) (map[core.Item]*entities.Resource, error) {
			result := map[core.Item]*entities.Resource{
				core.ItemWater: {
					Item:    core.ItemWater,
					Balance: 50,
				},
				core.ItemMedication: {
					Item:    core.ItemMedication,
					Balance: 63,
				},
			}
			return result, nil
		},
	}
	svc := New(repo)
	ctx := context.Background()
	res, err := svc.ResourceSharing(ctx)
	require.NoError(t, err)
	require.NotNil(t, res)
	water, ok := res[core.ItemWater.String()]
	require.True(t, ok)
	medic, ok := res[core.ItemMedication.String()]
	require.True(t, ok)

	assert.Equal(t, uint32(10), water.PerSurvivor)
	assert.Equal(t, uint32(50), water.Balance)

	assert.Equal(t, uint32(12), medic.PerSurvivor) // nearest whole number
	assert.Equal(t, uint32(63), medic.Balance)
}

func TestResourceSharingSurvivorRepoError(t *testing.T) {
	repo := &MockReportRepository{
		SurvivorsFunc: func(ctx context.Context) (*entities.Survivor, error) {
			return nil, gorm.ErrRecordNotFound
		},
		ResourcesFunc: func(ctx context.Context) (map[core.Item]*entities.Resource, error) {
			result := map[core.Item]*entities.Resource{
				core.ItemWater: {
					Item:    core.ItemWater,
					Balance: 50,
				},
				core.ItemMedication: {
					Item:    core.ItemMedication,
					Balance: 63,
				},
			}
			return result, nil
		},
	}
	svc := New(repo)
	ctx := context.Background()
	res, err := svc.ResourceSharing(ctx)
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	require.Nil(t, res)
}

func TestResourceSharingResourcesRepoError(t *testing.T) {
	repo := &MockReportRepository{
		SurvivorsFunc: func(ctx context.Context) (*entities.Survivor, error) {
			return &entities.Survivor{
				Total:      10,
				Clean:      5,
				Percentage: 50.0,
			}, nil
		},
		ResourcesFunc: func(ctx context.Context) (map[core.Item]*entities.Resource, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := New(repo)
	ctx := context.Background()
	res, err := svc.ResourceSharing(ctx)
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	require.Nil(t, res)
}
