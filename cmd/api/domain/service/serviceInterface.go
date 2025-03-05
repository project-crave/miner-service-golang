package service

import (
	craveModel "crave/shared/model"
)

type IService interface {
	Parse(step craveModel.Step, page craveModel.Page, name string) ([]string, error)
	Filter(name string, page craveModel.Page, filter craveModel.Filter) (int64, error)
	Refine(name string, page craveModel.Page, step craveModel.Step, filter craveModel.Filter) ([]string, error)
}
