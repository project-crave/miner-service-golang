package repository

import (
	"crave/miner/cmd/model"
	craveModel "crave/shared/model"
)

type IRepository interface {
	Save(name string, page craveModel.Page, targets []model.ParsedTarget) error
}
