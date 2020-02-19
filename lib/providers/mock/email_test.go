package mock_test

import (
	"testing"

	"github.com/matcornic/hermes/v2"
	"github.com/stretchr/testify/assert"

	"github.com/stiks/gobs/lib/providers/mock"
	"github.com/stiks/gobs/lib/repositories"
)

func TestMock_Email_NewEmailRepository(t *testing.T) {
	r := mock.NewEmailRepository()

	assert.Implements(t, (*repositories.EmailRepository)(nil), r)
}

func TestMock_Email_SendEmail(t *testing.T) {
	r := mock.NewEmailRepository()

	t.Run("Existing key", func(t *testing.T) {
		assert.NoError(t, r.SendEmail(nil, "john@snow.com", "Hello world", hermes.Email{}))
	})
}
