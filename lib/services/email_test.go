package services_test

import (
	"testing"

	"github.com/matcornic/hermes/v2"
	"github.com/stretchr/testify/assert"

	"github.com/stiks/gobs/lib/providers/mock"
	"github.com/stiks/gobs/lib/services"
)

func _emailSrv() services.EmailService {
	return services.NewEmailService(mock.NewEmailRepository())
}

func TestService_Email_NewEmailService(t *testing.T) {
	assert.Implements(t, (*services.EmailService)(nil), _emailSrv())
}

func TestService_Email_SendEmail(t *testing.T) {
	srv := _emailSrv()

	t.Run("Existing key", func(t *testing.T) {
		assert.NoError(t, srv.SendEmail(nil, "john@snow.com", "Hello world", hermes.Email{}))
	})
}
