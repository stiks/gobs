package dummy_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stiks/gobs/lib/providers/dummy"
	"github.com/stiks/gobs/lib/repositories"
)

func TestDummy_Cache_NewCacheRepository(t *testing.T) {
	r := dummy.NewCacheRepository()

	assert.Implements(t, (*repositories.CacheRepository)(nil), r)
}

func TestDummy_Cache_FindByKey(t *testing.T) {
	r := dummy.NewCacheRepository()

	err := r.FindByKey(nil, "random", nil)
	if assert.Error(t, err) {
		assert.EqualError(t, err, "cache not found", "error message %s", "formatted")
	}
}

func TestDummy_Cache_Create(t *testing.T) {
	r := dummy.NewCacheRepository()

	assert.NoError(t, r.Create(nil, "random", nil))
}

func TestDummy_Cache_Update(t *testing.T) {
	r := dummy.NewCacheRepository()

	assert.NoError(t, r.Update(nil, "random", nil))
}

func TestDummy_Cache_Delete(t *testing.T) {
	r := dummy.NewCacheRepository()

	assert.NoError(t, r.Delete(nil, "random"))
}

func TestDummy_Cache_Flush(t *testing.T) {
	r := dummy.NewCacheRepository()

	assert.NoError(t, r.Flush(nil))
}
