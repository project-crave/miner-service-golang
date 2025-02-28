package controller

import (
	craveModel "crave/shared/model"
)

type IController interface {
	Parse(step craveModel.Step, page craveModel.Page, name string) ([]string, error)
	Filter(name string, page craveModel.Page, filter craveModel.Filter) (int64, error)
}
