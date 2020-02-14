package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
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

func TestControllers_User_List(t *testing.T) {
	ctl := controllers.NewUserController(_userSrv)

	t.Run("All users", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set("USER_ID", "775a5b37-1742-4e54-9439-0357e768b011")

		err := ctl.List(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
}
