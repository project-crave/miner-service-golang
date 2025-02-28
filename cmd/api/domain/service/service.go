package service

import (
	filterBusiness "crave/miner/cmd/api/domain/business/filter"
	pageBusiness "crave/miner/cmd/api/domain/business/page"
	"crave/miner/cmd/api/infrastructure/externalApi"
	"crave/miner/cmd/api/infrastructure/repository"
	"crave/miner/cmd/model"

	craveModel "crave/shared/model"
)

type Service struct {
	pageStrat   *pageBusiness.PageStrategy
	filterStrat *filterBusiness.FilterStrategy
	repo        repository.IRepository
	hubClient   externalApi.IHubClient
}

func NewService(pageStrat *pageBusiness.PageStrategy, filterStrat *filterBusiness.FilterStrategy, repo repository.IRepository, hubClient externalApi.IHubClient) *Service {
	return &Service{pageStrat: pageStrat, filterStrat: filterStrat, repo: repo, hubClient: hubClient}
}

func (s *Service) Parse(step craveModel.Step, page craveModel.Page, name string) ([]string, error) {
	pageBiz := s.getPageBusiness(page)
	nameFilterBiz := s.getNameFilterBusiness()
	targets, err := s.getNextTargets(pageBiz, nameFilterBiz, step, name)
	if err != nil {
		return nil, err
	}
	s.repo.Save(name, page, targets)
	return s.getNames(targets), nil
}

func (s *Service) getNameFilterBusiness() filterBusiness.IBusiness {
	return s.filterStrat.GetNameFilterBusiness()
}

func (s *Service) getNames(targets []model.ParsedTarget) []string {
	targetNames := make([]string, len(targets))
	for i, target := range targets {
		targetNames[i] = target.Name
	}
	return targetNames
}

func (s *Service) getPageBusiness(page craveModel.Page) pageBusiness.IBusiness {
	return s.pageStrat.GetPageBusiness(page)
}

func (s *Service) getNextTargets(page pageBusiness.IBusiness, filter filterBusiness.IBusiness, step craveModel.Step, name string) ([]model.ParsedTarget, error) {
	return page.ParseNextTargets(step, filter, name)
}

func (s *Service) Filter(name string, page craveModel.Page, filter craveModel.Filter) (int64, error) {
	pageBiz := s.getPageBusiness(page)
	filterChain := s.getFilterChain(filter)

	filteredBy, err := pageBiz.ApplyFilter(name, *filterChain)
	return int64(filteredBy), err
}

func (s *Service) getFilterChain(filter craveModel.Filter) *filterBusiness.FilterChain {
	return s.filterStrat.GetFilterChain(filter)
}
