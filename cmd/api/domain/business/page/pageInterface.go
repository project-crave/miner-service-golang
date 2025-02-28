package business

import (
	filterBusiness "crave/miner/cmd/api/domain/business/filter"
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
	MakeFrontUrl(name string) string
	MakeBackUrl(name string) string
	GetHtml(url string) (*string, error)
	GetDocument(html *string) (*goquery.Document, error)
	ExtractFrontTargets(doc *goquery.Document, filter filterBusiness.IBusiness, name string) ([]model.ParsedTarget, error)
	ExtractBackTargets(doc *goquery.Document, filter filterBusiness.IBusiness, name string) ([]model.ParsedTarget, error)
	ParseNextTargets(step craveModel.Step, filter filterBusiness.IBusiness, name string) ([]model.ParsedTarget, error)
	ApplyFilter(name string, filterBiz filterBusiness.FilterChain) (craveModel.Filter, error)
}

func (strat *PageStrategy) GetPageBusiness(page craveModel.Page) IBusiness {
	return strat.pageMap[page]
}
