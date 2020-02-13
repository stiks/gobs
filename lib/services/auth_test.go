package services_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stiks/gobs/lib/providers/mock"
	"github.com/stiks/gobs/lib/services"
)

func _authSrv() services.AuthService {
	os.Setenv("AUTH_SECRET_KEY", "123")
	os.Setenv("AUTH_ACCESS_TOKEN_LIFETIME", "123")
	os.Setenv("AUTH_REFRESH_TOKEN_LIFETIME", "123")

	return services.NewAuthService(mock.NewAuthRepository())
}

func TestService_Auth_NewAuthRepository(t *testing.T) {
	assert.Implements(t, (*services.AuthService)(nil), _authSrv())
}
