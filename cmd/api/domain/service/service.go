package service

import (
	filterBusiness "crave/miner/cmd/api/domain/business/filter"
	pageBusiness "crave/miner/cmd/api/domain/business/page"
	"crave/miner/cmd/api/infrastructure/externalApi"
	"crave/miner/cmd/api/infrastructure/repository"
	"crave/miner/cmd/model"
	craveModel "crave/shared/model"
	"sync"
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
	if filteredBy != 0 {
		s.repo.Remove(name)
	}
	return int64(filteredBy), err
}

func (s *Service) getFilterChain(filter craveModel.Filter) *filterBusiness.FilterChain {
	return s.filterStrat.GetFilterChain(filter)
}

func (s *Service) Refine(name string, page craveModel.Page, step craveModel.Step, filter craveModel.Filter) ([]string, error) {
	pageBiz := s.getPageBusiness(page)
	nameFilterBiz := s.getNameFilterBusiness()
	targets, err := s.getNextTargets(pageBiz, nameFilterBiz, step, name)
	if err != nil {
		return nil, err
	}

	filterChain := s.getFilterChain(filter)

	filteredBy, _ := pageBiz.ApplyFilter(name, *filterChain)
	s.repo.SaveOrigin(name, (int64(filter) &^ int64(filteredBy)))

	indicesToRemove := make(map[int]bool)

	var wg sync.WaitGroup
	for i, target := range targets {
		wg.Add(1)
		go func() {
			defer wg.Done()
			filteredBy, err := pageBiz.ApplyFilter(target.Name, *filterChain)
			if err != nil {
				return
			}
			if filteredBy != 0 {
				indicesToRemove[i] = true
			}
			s.repo.SaveDestination(name, &target, page, (int64(filter) &^ int64(filteredBy)))
		}()

	}
	wg.Wait()

	var refinedTargetNames []model.ParsedTarget
	for i, target := range targets {
		if !indicesToRemove[i] {
			refinedTargetNames = append(refinedTargetNames, target)
		}
	}

	return s.getNames(refinedTargetNames), nil
}
