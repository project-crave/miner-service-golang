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

func (c *Controller) Parse(step craveModel.Step, page craveModel.Page, name string) ([]string, error) {
	return c.svc.Parse(step, page, name)
}

func (c *Controller) Filter(name string, page craveModel.Page, filter craveModel.Filter) (int64, error) {
	return c.svc.Filter(name, page, filter)
}
