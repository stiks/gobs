package controllers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stiks/gobs/lib/controllers"
	"github.com/stiks/gobs/lib/providers/mock"
	"github.com/stiks/gobs/lib/services"
)

var (
	_cacheSrv = services.NewCacheService(mock.NewCacheRepository())
	_queueSrv = services.NewQueueService(mock.NewQueueRepository())
	_userSrv  = services.NewUserService(mock.NewUserRepository(), _queueSrv, _cacheSrv)
)

func TestControllers_User_NewUserController(t *testing.T) {
	assert.NotNil(t, controllers.NewUserController(_userSrv))
}
