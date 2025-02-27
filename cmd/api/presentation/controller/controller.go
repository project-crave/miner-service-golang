package controller

import (
	"crave/miner/cmd/api/domain/service"
	craveModel "crave/shared/model"
)

type Controller struct {
	svc service.IService
}

func NewController(svc service.IService) *Controller {
	return &Controller{svc: svc}
}

// Parse implements IController.
func (c *Controller) Parse(step craveModel.Step, page craveModel.Page, name string) ([]string, error) {
	return c.svc.Parse(step, page, name)
}
