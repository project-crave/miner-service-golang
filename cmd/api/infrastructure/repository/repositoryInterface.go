package repository

import (
	"crave/miner/cmd/model"
	craveModel "crave/shared/model"
)

type IRepository interface {
	Save(name string, page craveModel.Page, targets []model.ParsedTarget) error
	SaveOrigin(name string, tag int64) error
	SaveDestination(org string, dest *model.ParsedTarget, page craveModel.Page, tag int64) error
	Remove(name string) error
}
