package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stiks/gobs/lib/providers/mock"
	"github.com/stiks/gobs/lib/services"
)

func _cacheSrv() services.CacheService {
	return services.NewCacheService(mock.NewCacheRepository())
}

func TestService_Cache_NewCacheService(t *testing.T) {
	assert.Implements(t, (*services.CacheService)(nil), _cacheSrv())
}

func TestService_Cache_FindByKey(t *testing.T) {
	srv := _cacheSrv()

	t.Run("Non-existing key", func(t *testing.T) {
		err := srv.GetByKey(nil, "random", nil)
		if assert.Error(t, err) {
			assert.EqualError(t, err, "cache not found", "error message %s", "formatted")
		}
	})

	t.Run("Existing key", func(t *testing.T) {
		assert.NoError(t, srv.GetByKey(nil, "random", true))
	})
}

func TestService_Cache_Create(t *testing.T) {
	srv := _cacheSrv()

	assert.NoError(t, srv.Create(nil, "random", nil))
}

func TestService_Cache_Update(t *testing.T) {
	srv := _cacheSrv()

	assert.NoError(t, srv.Update(nil, "random", nil))
}

func TestService_Cache_Delete(t *testing.T) {
	srv := _cacheSrv()

	assert.NoError(t, srv.Delete(nil, "random"))
}

func TestService_Cache_Flush(t *testing.T) {
	srv := _cacheSrv()

	assert.NoError(t, srv.Flush(nil))
}
