package mock_test

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/stiks/gobs/lib/providers/mock"
	"github.com/stiks/gobs/lib/repositories"
)

func TestMock_Stats_NewStatsRepository(t *testing.T) {
	r := mock.NewStatsRepository()

	assert.Implements(t, (*repositories.StatsRepository)(nil), r)
}

func TestMock_Stats_Process(t *testing.T) {
	r := mock.NewStatsRepository()

	assert.NoError(t, r.Process(echo.New().AcquireContext()))
}

func TestMock_Stats_GetStats(t *testing.T) {
	r := mock.NewStatsRepository()

	_, err := r.GetStats(nil)
	assert.NoError(t, err)
}
