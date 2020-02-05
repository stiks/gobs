package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"google.golang.org/appengine"

	"github.com/stiks/gobs/lib/controllers"
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

	// Core endpoints
	controllers.NewHealthController().Routes(e.Group("api"))

	appengine.Main()
}
