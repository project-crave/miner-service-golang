package controller

import (
	craveModel "crave/shared/model"
)

type IController interface {
	Parse(step craveModel.Step, page craveModel.Page, name string, filter craveModel.Filter) error
}
