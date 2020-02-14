package mock_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stiks/gobs/lib/providers/mock"
	"github.com/stiks/gobs/lib/repositories"
)

func TestMock_Cache_NewCacheRepository(t *testing.T) {
	r := mock.NewCacheRepository()

	assert.Implements(t, (*repositories.CacheRepository)(nil), r)
}

func TestMock_Cache_FindByKey(t *testing.T) {
	r := mock.NewCacheRepository()

	t.Run("Non-existing key", func(t *testing.T) {
		err := r.FindByKey(nil, "random", nil)
		if assert.Error(t, err) {
			assert.EqualError(t, err, "cache not found", "error message %s", "formatted")
		}
	})

	t.Run("Existing key", func(t *testing.T) {
		assert.NoError(t, r.FindByKey(nil, "random", true))
	})
}

func TestMock_Cache_Create(t *testing.T) {
	r := mock.NewCacheRepository()

	assert.NoError(t, r.Create(nil, "random", nil))
}

func TestMock_Cache_Update(t *testing.T) {
	r := mock.NewCacheRepository()

	assert.NoError(t, r.Update(nil, "random", nil))
}

func TestMock_Cache_Delete(t *testing.T) {
	r := mock.NewCacheRepository()

	assert.NoError(t, r.Delete(nil, "random"))
}

func TestMock_Cache_Flush(t *testing.T) {
	r := mock.NewCacheRepository()

	assert.NoError(t, r.Flush(nil))
}
