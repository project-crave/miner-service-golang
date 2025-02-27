package service

import "crave/shared/model"

type IService interface {
	Parse(step model.Step, page model.Page, name string) ([]string, error)
}
