package service

import (
	craveModel "crave/shared/model"
)

type IService interface {
	Parse(step craveModel.Step, page craveModel.Page, name string) ([]string, error)
	Filter(name string, page craveModel.Page, filter craveModel.Filter) (int64, error)
}
