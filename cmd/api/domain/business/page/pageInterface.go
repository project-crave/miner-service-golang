package business

import (
	"crave/miner/cmd/model"
	craveModel "crave/shared/model"

	"github.com/PuerkitoBio/goquery"
)

type PageStrategy struct {
	pageMap map[craveModel.Page]IBusiness
}

func NewStrategy(pageMap map[craveModel.Page]IBusiness) *PageStrategy {
	return &PageStrategy{pageMap: pageMap}
}

type IBusiness interface {
	MakeUrl(step craveModel.Step, name string) string
	GetDocument(url string) (*goquery.Document, error)
	ExtractTargets(doc *goquery.Document) ([]model.ParsedTarget, error)
	ParseNextTargets(step craveModel.Step, name string) ([]model.ParsedTarget, error)
}

func (strat *PageStrategy) GetPageBusiness(page craveModel.Page) IBusiness {
	return strat.pageMap[page]
}
