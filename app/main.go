package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"google.golang.org/appengine"

	"github.com/stiks/gobs/lib/controllers"
	"github.com/stiks/gobs/lib/providers/dummy"
	"github.com/stiks/gobs/lib/providers/mock"
	"github.com/stiks/gobs/lib/services"
)

func main() {
	e := echo.New()

	// Hide banner
	e.HideBanner = true

	e.Use(middleware.Logger())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server started, v%s", echo.Version)

	// Start server
	http.Handle("/", e)

	// Some stuff
	var (
		cacheSrv = services.NewCacheService(dummy.NewCacheRepository())
		queueSrv = services.NewQueueService(mock.NewQueueRepository())
		authSrv  = services.NewAuthService(mock.NewAuthRepository())
		userSrv  = services.NewUserService(mock.NewUserRepository(), queueSrv, cacheSrv)
	)

	// Core endpoints
	controllers.NewHealthController().Routes(e.Group("api"))
	controllers.NewWorkerController(queueSrv).Routes(e.Group("api"))

	// Base controllers
	controllers.NewAuthController(userSrv, authSrv).Routes(e.Group("api"))
	controllers.NewUserController(userSrv).Routes(e.Group("api"))

	appengine.Main()
}
