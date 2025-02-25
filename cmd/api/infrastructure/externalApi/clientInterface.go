package externalApi

import (
	craveModel "crave/shared/model"
)

type IHubClient interface {
	ParseResult(name string, targets []string, step craveModel.Step) error
}
