package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/stiks/gobs/lib/services"
)

type workerController struct {
	queue services.QueueService
}

// WorkerControllerInterface ...
type WorkerControllerInterface interface {
	Routes(g *echo.Group)
}

// NewWorkerController returns a controller
func NewWorkerController(queueSrv services.QueueService) WorkerControllerInterface {
	return &workerController{
		queue: queueSrv,
	}
}

// Routes registers routes
func (ctl *workerController) Routes(g *echo.Group) {
}
