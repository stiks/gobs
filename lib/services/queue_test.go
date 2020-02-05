package services_test

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stiks/gobs/lib/providers/mock"
	"github.com/stiks/gobs/lib/services"
	"github.com/stretchr/testify/assert"
)

func _queueSrv() services.QueueService {
	return services.NewQueueService(mock.NewQueueRepository())
}

func TestService_Queue_Add(t *testing.T) {
	srv := _queueSrv()

	data := "Test data"

	reqByte, err := json.Marshal(data)
	if err != nil {
		assert.NoError(t, err)
	}

	t.Run("Update queue", func(t *testing.T) {
		assert.NoError(t, srv.Add(nil, "775a5b37-1742-4e54-9439-0357e768b011", reqByte))
	})

	t.Run("Update non-existing queue", func(t *testing.T) {
		assert.NoError(t, srv.Add(nil, "5fcc94e5-c6aa-4320-8469-f5021af54b88", nil))
	})
}

func TestService_Queue_AddObject(t *testing.T) {
	srv := _queueSrv()

	newQueue := "Test data"

	t.Run("Create with object", func(t *testing.T) {
		assert.NoError(t, srv.AddObject(nil, "5fcc94e5-c6aa-4320-8469-f5021af54b88", &newQueue))
	})

	t.Run("Create with nil", func(t *testing.T) {
		assert.NoError(t, srv.AddObject(nil, "5fcc94e5-c6aa-4320-8469-f5021af54b88", nil))
	})
}

func TestService_Queue_AddToURL(t *testing.T) {
	srv := _queueSrv()

	data := url.Values{}

	t.Run("Update queue", func(t *testing.T) {
		assert.NoError(t, srv.AddToURL(nil, "775a5b37-1742-4e54-9439-0357e768b011", data))
	})

	t.Run("Update non-existing queue", func(t *testing.T) {
		assert.NoError(t, srv.AddToURL(nil, "5fcc94e5-c6aa-4320-8469-f5021af54b88", nil))
	})
}
