package mock_test

import (
	"net/url"
	"testing"

	"github.com/stiks/gobs/lib/providers/mock"
	"github.com/stiks/gobs/lib/repositories"
	"github.com/stretchr/testify/assert"
)

func TestMock_Queue_NewQueueRepository(t *testing.T) {
	r := mock.NewQueueRepository()

	assert.Implements(t, (*repositories.QueueRepository)(nil), r)
}

func TestMock_Queue_Add(t *testing.T) {
	r := mock.NewQueueRepository()

	assert.NoError(t, r.Add(nil, "test", []byte("data")))
}

func TestMock_Queue_AddObject(t *testing.T) {
	r := mock.NewQueueRepository()

	assert.NoError(t, r.AddObject(nil, "test", "data"))
}

func TestMock_Queue_AddToURL(t *testing.T) {
	r := mock.NewQueueRepository()

	var data url.Values

	assert.NoError(t, r.AddToURL(nil, "test", data))
}
